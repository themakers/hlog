package hlog

import (
	"fmt"
	"reflect"
	"strings"

	"go.uber.org/zap"
)

type Logger interface {
	Emit(m interface{})
	With(m interface{}) Logger
	Module(name string) Logger
}

var _ Logger = new(logger)

type logger struct {
	log *zap.Logger
}

func New(log *zap.Logger) Logger {
	return &logger{
		log: log,
	}
}

func (l *logger) Emit(m interface{}) {
	if e, ok := m.(Emitter); ok {
		e.Emit(l.log)
	} else {
		l.emit(m)
	}
}

func (l *logger) emit(m interface{}) {
	type level string

	const (
		levelDebug level = "Debug"
		levelInfo  level = "Info"
		levelWarn  level = "Warn"
		levelError level = "Error"
		levelPanic level = "Panic"
	)

	typeName, name, fields := prepare2(m)

	emit := l.log.Panic

	switch {
	case strings.HasPrefix(typeName, string(levelDebug)):
		emit = l.log.Debug
	case strings.HasPrefix(typeName, string(levelInfo)):
		emit = l.log.Info
	case strings.HasPrefix(typeName, string(levelWarn)):
		emit = l.log.Warn
	case strings.HasPrefix(typeName, string(levelError)):
		emit = l.log.Error
	case strings.HasPrefix(typeName, string(levelPanic)):
		emit = l.log.Panic
	default:
		panic(fmt.Sprintf("could not determine event level: %s", typeName))
	}

	emit(name, fields...)
}

func (l *logger) With(m interface{}) Logger {
	_, _, fields := prepare2(m)
	return New(l.log.With(fields...))
}

func (l *logger) Module(name string) Logger {
	return New(l.log.Named(name))
}

type Emitter interface {
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
