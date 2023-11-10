package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"stocksapp/pkg/discovery"
	"stocksapp/pkg/discovery/consul"
	"stocksapp/stock/internal/controller/stock"
	metadatagateway "stocksapp/stock/internal/gateway/metadata/http"
	ratinggateway "stocksapp/stock/internal/gateway/rating/http"
	httphandler "stocksapp/stock/internal/handler/http"
	"time"
)

const serviceName = "stock"

func main() {
	var port int

	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()

	log.Println("Starting the movie service")

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}

			time.Sleep(1 * time.Second)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := stock.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)

	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
