package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kelseyhightower/envconfig"

	"github.com/mrksmt/test-task-7/internal/server/controller"
	"github.com/mrksmt/test-task-7/internal/server/service"
	"github.com/mrksmt/test-task-7/internal/server/storage"
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

	cfg := controller.Parameters{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	stor := storage.NewFakeStorage()
	srv := service.NewService(stor)
	c := controller.NewController(&cfg, srv)

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
