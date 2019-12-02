package log

import (
	"fmt"
	"reflect"
	"strings"

	"go.uber.org/zap"
)

func Emit(log *zap.Logger, m interface{}) {
	if e, ok := m.(Event); ok {
		e.Emit(log)
	} else {
		// FIXME Replace sugared logger with regular
		emit(log.Sugar(), m)
	}
}

type Event interface {
	Emit(*zap.Logger)
}

type level string

const (
	levelDebug level = "Debug"
	levelInfo  level = "Info"
	levelWarn  level = "Warn"
	levelError level = "Error"
	levelPanic level = "Panic"
)

func emit(logger *zap.SugaredLogger, m interface{}) {
	rv := reflect.ValueOf(m)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	rt := rv.Type()

	var fields = make([]interface{}, rt.NumField()*2)
	for i := 0; i < rv.NumField(); i++ {
		sf := rt.Field(i)
		fv := rv.Field(i)

		fields[i*2] = sf.Name
		var v interface{}
		if fv.CanInterface() {
			v = fv.Interface()
		} else {
			v = nil
		}
		fields[i*2+1] = v
	}

	typeName := rt.Name()

	name := fmt.Sprintf("%s!%s", rt.PkgPath(), typeName)

	emit := logger.Panicw

	switch {
	case strings.HasPrefix(typeName, string(levelDebug)):
		emit = logger.Debugw
	case strings.HasPrefix(typeName, string(levelInfo)):
		emit = logger.Infow
	case strings.HasPrefix(typeName, string(levelWarn)):
		emit = logger.Warnw
	case strings.HasPrefix(typeName, string(levelError)):
		emit = logger.Errorw
	case strings.HasPrefix(typeName, string(levelPanic)):
		emit = logger.Panicw
	default:
		panic(fmt.Sprintf("could not determine event level: %s", typeName))
	}

	emit(name, fields...)

	return
}
