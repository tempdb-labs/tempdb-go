package main

import (
	"log"
	"net/http"

	"github.com/ThembinkosiThemba/zen"
	tempdb "github.com/tempdb-labs/tempdb-go/lib"
)

func main() {
	app := zen.New()

	config := tempdb.Config{
		Addr: "0.0.0.0:8081",
		URL:  "tempdb://admin:123456789@workspace:cb4552273c5c/ecommerce",
	}

	client, err := tempdb.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	if err != nil {
		log.Fatalf("Failed to get client: %v", err)
	}
	defer client.Close()

	app.Apply(zen.Logger())
	DocumentGroup(app, client)

	app.Serve(":8082")

}

func DocumentGroup(r *zen.Engine, client *tempdb.TempDBClient) {
	r.GET("/c/1", func(ctx *zen.Context) {
		res, err := client.PubSubChannels()
		if err != nil {
			ctx.Error(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.Success(http.StatusOK, res, "OK")
	})

	r.GET("/c/2", func(ctx *zen.Context) {
		res, err := client.XList()
		if err != nil {
			ctx.Error(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.Success(http.StatusOK, res, "OK")
	})

	r.GET("/c/3", func(ctx *zen.Context) {
		res, err := client.QPeek("myqueue")
		if err != nil {
			ctx.Error(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.Success(http.StatusOK, res, "OK")
	})

	r.GET("/c/4", func(ctx *zen.Context) {
		res, err := client.QList()
		if err != nil {
			ctx.Error(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.Success(http.StatusOK, res, "OK")
	})
}
