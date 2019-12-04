package main

import (
	"math/rand"

	"github.com/themakers/hlog"
	"github.com/themakers/hlog/demo/events"
	"go.uber.org/zap"
)

func main() {
	zlog, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	//> You must add 2 to callerSkip to see correct caller info
	zlog = zlog.WithOptions(zap.AddCallerSkip(2))

	hlog := hlog.New(zlog)

	hlog.Emit(events.WarnSampleEvent{
		Field1: rand.Int(),
		Field2: "value",
	})
}
