package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kelseyhightower/envconfig"

	"github.com/mrksmt/test-task-7/internal/client"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := new(sync.WaitGroup)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		<-sigChan
		cancel()
	}()

	cfg := client.Parameters{}

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	c := client.NewSentenceClient(&cfg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := c.Run(ctx)
		if err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()
}
