package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"stocksapp/metadata/internal/controller/metadata"
	httphandler "stocksapp/metadata/internal/handler/http"
	"stocksapp/metadata/internal/repository/memory"
	"stocksapp/pkg/discovery"
	"stocksapp/pkg/discovery/consul"
	"time"
)

const serviceName = "metadata"

func main() {
	var port int

	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()

	log.Printf("Starting the movie metadata service on port %d", port)

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
	defer fmt.Println("sending deregister")

	repo := memory.New()
	ctrl := metadata.New(repo)
	handler := httphandler.New(ctrl)

	http.Handle("/metadata", http.HandlerFunc(handler.GetMetadata))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
