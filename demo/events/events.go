package events

import "github.com/themakers/hlog"

type Warn_SampleEvent2 struct {
	Field1 int
	Field2 string
}

type Warn_SampleEvent8 struct {
	Field1 int
	Field2 string
	Field3 int
	Field4 string
	Field5 int
	Field6 string
	Field7 int
	Field8 string
}

type Debug_UnimportantEvent hlog.FluidEvent
