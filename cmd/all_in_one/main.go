package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/goava/di"
	"github.com/goava/slice"

	"github.com/mrksmt/test-task-7/internal/client"
	"github.com/mrksmt/test-task-7/internal/service"
)

var address = "localhost:8090"

func xxx() {

	ctx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)

	// run server

	srv := service.NewSentenceService(&service.Parameters{Address: address})
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := srv.Run(ctx)
		if err != nil {
			log.Println("server run error:", err)
			return
		}
	}()

	<-time.After(time.Second)

	// run client

	cl := client.NewSentenceClient(&client.Parameters{Address: address})
	sentence, err := cl.GetSentence(ctx)
	if err != nil {
		log.Println("client run error:", err)
		return
	}

	log.Println("client got sentence:", sentence)

	// wait for server stop

	cancel()
	wg.Wait()
}

func main() {
	slice.Run(
		slice.WithName("sentence-all-in-one"),
		slice.WithParameters(
			&client.Parameters{},
			&service.Parameters{},
		),
		slice.WithComponents(
			slice.Provide(client.NewSentenceClient, di.As(new(slice.Dispatcher))),
			slice.Provide(service.NewSentenceService, di.As(new(slice.Dispatcher))),
		),
	)
}
