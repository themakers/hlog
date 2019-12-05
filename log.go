package hlog

import (
	"fmt"
	"reflect"
	"strings"

	"go.uber.org/zap"
)

type FluidEvent map[string]interface{}

type Emitter interface {
	Emit(*zap.Logger)
}

type Logger interface {
	Emit(e interface{}) interface{}
	With(e interface{}) Logger
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

func (l *logger) Emit(e interface{}) interface{} {
	if em, ok := e.(Emitter); ok {
		em.Emit(l.log)
	} else {
		l.emit(e)
	}
	return e
}

func (l *logger) emit(e interface{}) {
	type level string

	const (
		levelDebug level = "Debug"
		levelInfo  level = "Info"
		levelWarn  level = "Warn"
		levelError level = "Error"
	)

	typeName, name, fields := prepare2(e)

	emit := l.log.Debug

	switch {
	case strings.HasPrefix(typeName, string(levelDebug)):
		emit = l.log.Debug
	case strings.HasPrefix(typeName, string(levelInfo)):
		emit = l.log.Info
	case strings.HasPrefix(typeName, string(levelWarn)):
		emit = l.log.Warn
	case strings.HasPrefix(typeName, string(levelError)):
		emit = l.log.Error
	default:
		panic(fmt.Sprintf("could not determine event level: %s", typeName))
	}

	emit(name, fields...)
}

func (l *logger) With(e interface{}) Logger {
	_, _, fields := prepare2(e)
	return New(l.log.With(fields...))
}

func (l *logger) Module(name string) Logger {
	return New(l.log.Named(name))
}

func EventName(m interface{}) string {
	_, _, _, name := prepare1(m)
	return name
}

func prepare1(e interface{}) (rt reflect.Type, rv reflect.Value, typeName, name string) {
	rv = reflect.ValueOf(e)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	rt = rv.Type()

	typeName = rt.Name()

	name = fmt.Sprintf("%s!%s", rt.PkgPath(), typeName)

	return rt, rv, typeName, name
}

func prepareFieldsStruct(rt reflect.Type, rv reflect.Value) (fields []zap.Field) {
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

	return fields
}

func prepareFieldsMap(rt reflect.Type, rv reflect.Value) (fields []zap.Field) {
	fields = make([]zap.Field, rv.Len())
	mi := rv.MapRange()
	for i := 0; mi.Next(); i++ {
		fields[i] = zap.Any(mi.Key().String(), mi.Value().Interface())
	}

	return fields
}

func prepare2(e interface{}) (typeName, name string, fields []zap.Field) {
	rt, rv, typeName, name := prepare1(e)

	switch rt.Kind() {
	case reflect.Struct:
		fields = prepareFieldsStruct(rt, rv)
	case reflect.Map:
		fields = prepareFieldsMap(rt, rv)
	default:
		panic(fmt.Sprintf("bad kind of event: %s is %v", typeName, rt.Kind()))
	}

	return typeName, name, fields
}
