package main

import nats "github.com/nats-io/nats.go"
import "github.com/syndtr/goleveldb/leveldb"
import "time"
import "fmt"
//import "encoding/binary"

import (
	"log"
	"os"
	"os/signal"
)


func main () {
    // Connect to a server
    nc, _ := nats.Connect("nats://172.31.46.240:4222")
    
    db, _ := leveldb.OpenFile("/tmp/db", nil)
    defer db.Close()
    
    count  := 10*1000000

    received := 0
    start := time.Now()
    
    sub, _ :=  nc.QueueSubscribe("test", "result.workers", func(m *nats.Msg) {
        
         received++;

         if received == 1 {
            start = time.Now()
         }
        
         db.Put([]byte(fmt.Sprintf("key-%010d", received)) , []byte(m.Data), nil)
    
         if received % 100000 == 0{
             log.Printf("+")
         }
         if received == count {
        	elapsed := time.Since(start)
            fmt.Println("Set wps", float64(count)/elapsed.Seconds())
        
            // Unsubscribe
            //sub.Unsubscribe()
            // Drain
            //sub.Drain()
            // Close connection
            nc.Close()
         }
    })
    sub.SetPendingLimits(-1, -1)


	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Printf("Draining...")
	nc.Drain()
	log.Fatalf("Exiting")
	

}



