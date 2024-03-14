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
	"github.com/mrksmt/test-task-7/internal/server/controller"
	"github.com/mrksmt/test-task-7/internal/server/service"
	"github.com/mrksmt/test-task-7/internal/server/storage"
)

func main() {

	// common part ----------

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

	// run server ----------

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

	// run client ----------

	ClientCfg := client.Parameters{}

	err = envconfig.Process("", &ClientCfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	cl := client.NewSentenceClient(&ClientCfg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := cl.Run(ctx)
		if err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()
}
