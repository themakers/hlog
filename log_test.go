package hlog_test

import (
	"io/ioutil"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/themakers/hlog"
	"github.com/themakers/hlog/demo/events"
)

//> Run with `go test -v -benchmem -bench=. .`

func newBlackholeLogger() *zap.Logger {
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

func BenchmarkEmitStruct2(b *testing.B) {
	logger := hlog.New(newBlackholeLogger())

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		logger.Emit(events.Warn_SampleEvent2{
			Field1: 10000000,
			Field2: "value",
		})
	}
}

func BenchmarkEmitStruct8(b *testing.B) {
	logger := hlog.New(newBlackholeLogger())

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		logger.Emit(events.Warn_SampleEvent8{
			Field1: 10000000,
			Field2: "value",
			Field3: 10000000,
			Field4: "value",
			Field5: 10000000,
			Field6: "value",
			Field7: 10000000,
			Field8: "value",
		})
	}
}

func BenchmarkEmitMap2(b *testing.B) {
	logger := hlog.New(newBlackholeLogger())

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		logger.Emit(events.Debug_UnimportantEvent{
			"Field1": 10000000,
			"Field2": "value",
		})
	}
}

func BenchmarkEmitMap8(b *testing.B) {
	logger := hlog.New(newBlackholeLogger())

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		logger.Emit(events.Debug_UnimportantEvent{
			"Field1": 10000000,
			"Field2": "value",
			"Field3": 10000000,
			"Field4": "value",
			"Field5": 10000000,
			"Field6": "value",
			"Field7": 10000000,
			"Field8": "value",
		})
	}
}

func BenchmarkPure2(b *testing.B) {
	logger := newBlackholeLogger()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		logger.Warn(
			"github.com/themakers/hlog/demo/events!Warn_SampleEvent2",
			zap.Int("Field1", 10000000),
			zap.String("Field2", "value"),
		)
	}
}


func BenchmarkPure8(b *testing.B) {
	logger := newBlackholeLogger()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		logger.Warn(
			"github.com/themakers/hlog/demo/events!Warn_SampleEvent2",
			zap.Int("Field1", 10000000),
			zap.String("Field2", "value"),
			zap.Int("Field3", 10000000),
			zap.String("Field4", "value"),
			zap.Int("Field5", 10000000),
			zap.String("Field6", "value"),
			zap.Int("Field7", 10000000),
			zap.String("Field8", "value"),
		)
	}
}
