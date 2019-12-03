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
		emit(log, m)
	}
}

type Event interface {
	Emit(*zap.Logger)
}

func Name(m interface{}) string {
	_, _, _, name := prepare1(m)
	return name
}

func prepare1(m interface{}) (rt reflect.Type, rv reflect.Value, typeName, name string) {
	rv = reflect.ValueOf(m)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	rt = rv.Type()

	typeName = rt.Name()

	name = fmt.Sprintf("%s!%s", rt.PkgPath(), typeName)

	return rt, rv, typeName, name
}

func prepare2(m interface{}) (typeName, name string, fields []zap.Field) {
	rt, rv, typeName, name := prepare1(m)

	fields = make([]zap.Field, rt.NumField())
	for i := 0; i < rv.NumField(); i++ {
		sf := rt.Field(i)
		fv := rv.Field(i)

		var v interface{}
		if fv.CanInterface() {
			v = fv.Interface()
		} else {
			v = nil
		}

		fields[i] = zap.Any(sf.Name, v)
	}

	return typeName, name, fields
}

type level string

const (
	levelDebug level = "Debug"
	levelInfo  level = "Info"
	levelWarn  level = "Warn"
	levelError level = "Error"
	levelPanic level = "Panic"
)

func emit(logger *zap.Logger, m interface{}) {
	typeName, name, fields := prepare2(m)

	emit := logger.Panic

	switch {
	case strings.HasPrefix(typeName, string(levelDebug)):
		emit = logger.Debug
	case strings.HasPrefix(typeName, string(levelInfo)):
		emit = logger.Info
	case strings.HasPrefix(typeName, string(levelWarn)):
		emit = logger.Warn
	case strings.HasPrefix(typeName, string(levelError)):
		emit = logger.Error
	case strings.HasPrefix(typeName, string(levelPanic)):
		emit = logger.Panic
	default:
		panic(fmt.Sprintf("could not determine event level: %s", typeName))
	}

	emit(name, fields...)
}
