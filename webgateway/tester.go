package main

import (
  "fmt"
  "time"
  "gopkg.in/h2non/gentleman.v2"
  "gopkg.in/h2non/gentleman.v2/plugins/body"
  "gopkg.in/h2non/gentleman.v2/plugins/timeout"
 
)

func main() {

  for i:=0; i < 10000; i++ {

      // Create a new client
      cli := gentleman.New()
    
      // Define the Base URL
    
      cli.URL("http://172.31.10.75:8080")
      
      // Define the max timeout for the whole HTTP request
      cli.Use(timeout.Request(3 * time.Second))
    
      // Define dial specific timeouts
      cli.Use(timeout.Dial(1*time.Second, 3*time.Second))
  
      // Create a new request based on the current client
      req := cli.Request()
    
      // Method to be used
      req.Method("POST")

      req.Path(fmt.Sprintf("/api/%u",i))

      // Define the JSON payload via body plugin
      data := map[string]string{"counter": string(i)}
      req.Use(body.JSON(data))
    
      // Perform the request
      res, err := req.Send()
      if err != nil {
        fmt.Printf("Request error: %s\n", err)
      }
      if !res.Ok {
        fmt.Printf("Invalid server response: %d\n", res.StatusCode)
      }
      res.Close()
 
      query := cli.Request()
    
      // Method to be used
      query.Method("GET")

      query.Path(fmt.Sprintf("/api/%u",i))
      
      res, err = query.Send()
      if err != nil {
        fmt.Printf("Request error: %s\n", err)
      }
      if !res.Ok {
        fmt.Printf("Invalid server response: %d\n", res.StatusCode, "key", i)
      }
      res.Close()
      fmt.Printf("+")
  }
  //fmt.Printf("Status: %d\n", res.StatusCode)
  //fmt.Printf("Body: %s", res.String())
}