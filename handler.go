package main

import (
	"reflect"
)

type handler struct {
	val      reflect.Value
	inTypes  []reflect.Type
	outTypes []reflect.Type
}

func newHandler(fn interface{}) handler {
	val := reflect.ValueOf(fn)
	fnType := val.Type()
	switch {
	case fnType.NumIn() > 2:
		panic("handler function has at most 2 inputs")
	case fnType.NumOut() > 1:
		panic("handler function has at most 1 output")
	case fnType.NumOut() == 1 && fnType.Out(0).String() != "main.Result":
		panic("handler function output type must be Result, got " + fnType.Out(0).String())
	}

	inTypes, outTypes := []reflect.Type{}, []reflect.Type{}
	for i := 0; i < fnType.NumIn(); i++ {
		inType := fnType.In(i)
		if inType.Kind() != reflect.Struct {
			panic("handler function input must be a struct, got " + inType.Kind().String())
		}
		inTypes = append(inTypes, inType)
	}
	for i := 0; i < fnType.NumOut(); i++ {
		outTypes = append(outTypes, fnType.Out(i))
	}

	return handler{
		val:      val,
		inTypes:  inTypes,
		outTypes: outTypes,
	}
}

func (h handler) handle(ctx Context) *Result {
	inputs := make([]reflect.Value, len(h.inTypes))
	for i, inType := range h.inTypes {
		if inType.String() == "main.Context" {
			inputs[i] = reflect.ValueOf(ctx)
			continue
		}
		input := reflect.New(inType)
		if err := ctx.unpack(input.Interface()); err != nil {
			return &Result{Done: true, ErrInternal: err}
		}
		inputs[i] = input.Elem()
	}
	outs := h.val.Call(inputs)
	if len(h.outTypes) == 0 {
		return nil
	}
	result := outs[0].Interface().(Result)
	return &result
}
