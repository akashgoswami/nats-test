package main

import (
	"time"
	//"fmt"
	"github.com/aerogo/aero"
	"github.com/nats-io/nats.go"	
)

func main() {
	
	
    nc, _ := nats.Connect("nats://172.31.46.240:4222")
    
	app := aero.New()
	configure(app,nc).Run()
}

func configure(app *aero.Application, nc *nats.Conn) *aero.Application {
	
	app.Get("/api/:key", func(ctx aero.Context) error {
		
		key := ctx.Get("key")

		msg, err := nc.Request("query", []byte(key), 5*time.Second)
    	if err != nil {
    		return ctx.Error(404)
    	} else {
    		return ctx.Bytes(msg.Data);
    	}
	})

	app.Post("/api/:key", func(ctx aero.Context) error {
		
		key := ctx.Get("key")

		key = "store."+key
		body, _ := ctx.Request().Body().Bytes()

		//fmt.Println("Save key", key, "with content", string(body))
		
		nc.Publish(key, body)
		nc.Flush()
   		return ctx.String("OK")
   		
   	})

	app.Delete("/api/:key", func(ctx aero.Context) error {
		
		key := ctx.Get("key")
		//fmt.Println("Save key", key, "with content", string(body))
		_, err := nc.Request("delete", []byte(key), 5*time.Second)		
    	if err != nil {
    		return ctx.Error(404)
    	} else {
    		return ctx.String("Ok")
    	}
   	})
   	
	return app
}
