package main

import (
	"math/rand"

	"github.com/themakers/log"
	"github.com/themakers/log/demo/events"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	//> You must add 2 to callerSkip to see correct caller info
	logger = logger.WithOptions(zap.AddCallerSkip(2))

	log.Emit(logger, events.WarnSampleEvent{
		Field1: rand.Int(),
		Field2: "value",
	})
}
