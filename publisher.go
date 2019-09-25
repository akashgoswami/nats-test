package main

import (
	"fmt"
	"time"
	"math/rand"
     "encoding/binary"	
	"github.com/nats-io/nats.go"
)


func main() {
    
    total := 10*1000000
    chunk := total / 100 
  
    nc, err := nats.Connect("nats://172.31.46.240:4222")
  
    if err != nil {
         fmt.Println("Unable to connect to the server ", err)
         return
    }
    
	start := time.Now()
    
    bs := make([]byte, 24)
  
	for i := 0; i < total; i++ {
	    binary.LittleEndian.PutUint32(bs, uint32(i))
	    
		nc.Publish( fmt.Sprintf("state.%d", i), bs)

		if i % chunk == 0{
		    fmt.Printf("+")
		    nc.Flush()
		}
	}
  	elapsed := time.Since(start)
    fmt.Println("Published at rate m/s", float64(total)/elapsed.Seconds())
    nc.Flush()
    time.Sleep(10 * time.Second)
 
	for i := 0; i < chunk; i++ {
	    
        query := rand.Intn(total)
        //fmt.Println("Querying", query)
        msg, err := nc.Request("state.query", []byte(fmt.Sprintf("state.%d", query)), 3*time.Second)
    	if err != nil {
    		fmt.Println("%v for request", err)
    		return;
    	}
    	//fmt.Println(msg.Data)
        val := binary.LittleEndian.Uint32(msg.Data)
        if val != uint32(query) {
            fmt.Println("Validation error with key ", val, "expected", query)
            //break  
        }
		if i % 1000 == 0{
		    fmt.Printf(".")
		}
		

	}


	nc.Flush()
	nc.Drain()
	
        
}
         
            