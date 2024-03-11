package main

import (
	"github.com/goava/di"
	"github.com/goava/slice"

	"github.com/mrksmt/test-task-7/internal/service"
)

func main() {
	slice.Run(
		slice.WithName("sentence-service"),
		slice.WithParameters(
			&service.Parameters{},
		),
		slice.WithComponents(
			slice.Provide(service.NewSentenceService, di.As(new(slice.Dispatcher))),
		),
	)
}
