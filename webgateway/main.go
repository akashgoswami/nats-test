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
	
	app.Get("/query/:key", func(ctx aero.Context) error {
		
		key := ctx.Get("key")
		key = "state."+key
		msg, err := nc.Request("state.query", []byte(key), 5*time.Second)
    	if err != nil {
    		return ctx.String("Not Found")
    	} else {
    		return ctx.Bytes(msg.Data);
    	}
	})

	app.Post("/store/:key", func(ctx aero.Context) error {
		
		key := ctx.Get("key")
		key = "state."+key
		body, _ := ctx.Request().Body().Bytes()

		//fmt.Println("Save key", key, "with content", string(body))
		
		nc.Publish(key, body)
   		return ctx.String("OK")
   		
   	})
	
	return app
}
