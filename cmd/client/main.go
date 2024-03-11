package main

import (
	"github.com/goava/di"
	"github.com/goava/slice"

	"github.com/mrksmt/test-task-7/internal/client"
)

func main() {
	slice.Run(
		slice.WithName("sentence-client"),
		slice.WithParameters(
			&client.Parameters{},
		),
		slice.WithComponents(
			slice.Provide(client.NewSentenceClient, di.As(new(slice.Dispatcher))),
		),
	)
}
