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
	"strings"
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
    deleted := 0
    
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
                // write any unstored transactions to db
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
    

    sub, _ :=  nc.QueueSubscribe("store.>", "store.workers", func(m *nats.Msg) {
        
         received++;

         if received == 1 {
            start = time.Now()
         }
        key := strings.TrimPrefix(m.Subject, "store.")
        
        //fmt.Println("Received store", m.Data)
        batch.Put([]byte(key) , []byte(m.Data))
    
        if received % 2048 == 0 {
            mutex.Lock()
            err = db.Write(batch, nil)
            batch.Reset();
            mutex.Unlock()
        }
        
        if received % chunk == 0{
        	elapsed := time.Since(start)
            fmt.Println("Received last", chunk, "at rate m/s", float64(chunk)/elapsed.Seconds())
            
            start = time.Now()
         }
    })
    sub.SetPendingLimits(100000, -1)

    nc.Subscribe("query", func(m *nats.Msg) {
        
        //fmt.Println("Received query", string(m.Data))
        data, err := db.Get([]byte(m.Data), nil)
        if err != nil {
            //fmt.Println("Could not find", string(m.Data))
            missed++
        } else {
            //fmt.Println("Found", string(data))
            m.Respond(data)
            nc.Flush();
            hit++
        }
    })

    nc.Subscribe("delete", func(m *nats.Msg) {
        
        //fmt.Println("Received query", string(m.Data))
        exists, _ := db.Has([]byte(m.Data), nil)

        if exists {
            db.Delete([]byte(m.Data), nil)
            //fmt.Println("Found", string(data))
            m.Respond([]byte("OK"))
            nc.Flush();
            deleted++
        } else {
            missed++;
        }
    })
    
    
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	ticker.Stop();
	log.Printf("Draining...")
	nc.Drain()
	log.Fatalf("Exiting")


}



