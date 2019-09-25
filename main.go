package main

import nats "github.com/nats-io/nats.go"
import "time"
import "fmt"
//import "encoding/binary"
import "github.com/syndtr/goleveldb/leveldb"
import "github.com/syndtr/goleveldb/leveldb/opt"
import "github.com/syndtr/goleveldb/leveldb/filter"


import (
	"log"
	"os"
	"os/signal"
	"sync"
)


func main () {
    // Connect to a server
    nc, _ := nats.Connect("nats://172.31.46.240:4222")
    
    o := &opt.Options{
    	Filter: filter.NewBloomFilter(10),
    }

    db, err := leveldb.OpenFile("/var/db/state", o)
    
    if err != nil {
         fmt.Println("Unable to open Database ", err)
         return
    }
    defer db.Close()
    
    chunk := 1000000
    
    received := 0
    missed := 0
    hit := 0
    
    start := time.Now()
    batch := new(leveldb.Batch)
    var mutex = &sync.Mutex{}

    ticker := time.NewTicker(1000 * time.Millisecond)
    done := make(chan bool)
    go func() {
        for {
            select {
            case <-done:
                return
            case t := <-ticker.C:
                //fmt.Println("Tick at", t)
                t = t
                if batch.Len() > 0 {
                    mutex.Lock()
                    err = db.Write(batch, nil)
                    batch.Reset();
                    //batch = new(leveldb.Batch)
                    mutex.Unlock()
                }
            }
        }
    }()
    

    sub, _ :=  nc.QueueSubscribe("state.>", "state.workers", func(m *nats.Msg) {
        
         received++;

         if received == 1 {
            start = time.Now()
         }

        mutex.Lock()
        batch.Put([]byte(m.Subject) , []byte(m.Data))
        if received % 1000 == 0{
            err = db.Write(batch, nil)
            batch.Reset();
        }
        mutex.Unlock()
        
        if received % chunk == 0{
        	elapsed := time.Since(start)
            fmt.Println("Received last", chunk, "at rate m/s", float64(chunk)/elapsed.Seconds())
            
            start = time.Now()
         }
    })
    sub.SetPendingLimits(-1, -1)

    querysub, _ := nc.Subscribe("state.query", func(m *nats.Msg) {
        //fmt.Println("Received query", m.Data)
        
        data, err := db.Get([]byte(m.Data), nil)
        if err != nil {
            //fmt.Println("Coud not find", string(m.Data))
            missed++
        } else {
            //fmt.Println("Found", string(data))
            m.Respond(data)
            nc.Flush();
            hit++
        }
            


    })
    querysub.SetPendingLimits(-1, -1)


	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	ticker.Stop();
	log.Printf("Draining...")
	nc.Drain()
	log.Fatalf("Exiting")


}



