package main

import "github.com/syndtr/goleveldb/leveldb"
import "time"
import "fmt"
import  "encoding/binary"
//import "math/rand" 

func main() {
    
    count  := 1000*1000000

    db, _ := leveldb.OpenFile("/home/ubuntu/environment/nats-test/db", nil)
    defer db.Close()

    start := time.Now()
    
    bs := make([]byte, 4)
     
    for i := 0; i < count; i++ {
        binary.LittleEndian.PutUint32(bs, uint32(i))
        //rand.Read(bs)
        db.Put([]byte(fmt.Sprintf("key-%010d", i)) , bs, nil)
	}
	elapsed := time.Since(start)
    fmt.Println("Set wps", float64(count)/elapsed.Seconds())

   for i := 0; i < count; i++ {
        binary.LittleEndian.PutUint32(bs, uint32(i))
        //rand.Read(bs)
        data, _ := db.Get( []byte(fmt.Sprintf("key-%010d", i)), nil)
        val := binary.LittleEndian.Uint32(data)
        if val != uint32(i) {
            fmt.Println("Validation error with key ", i)
            break  
        }
        
	}
    
  
}