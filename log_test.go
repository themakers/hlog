package hlog

import (
	"io/ioutil"
	"testing"

	"github.com/themakers/hlog/demo/events"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//> Run with `go test -v -benchmem -bench=. .`

func newLogger() *zap.Logger {
	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(ioutil.Discard),
			zapcore.DebugLevel,
		),
	).WithOptions(
		zap.AddCallerSkip(2),
	)
}

func BenchmarkEmit(b *testing.B) {
	logger := New(newLogger())

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		logger.Emit(events.WarnSampleEvent{
			Field1: 10000000,
			Field2: "value",
		})
	}
}

func BenchmarkPure(b *testing.B) {
	logger := newLogger()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		logger.Warn(
			"github.com/themakers/log/demo/events!WarnSampleEvent",
			zap.Int("Field1", 10000000),
			zap.String("Field2", "value"),
		)
	}
}
