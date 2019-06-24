package httpwrap

import "reflect"

type mainFn struct {
	val      reflect.Value
	inTypes  []reflect.Type
	outTypes []reflect.Type
}

func newMain(fn interface{}) (mainFn, error) {
	val := reflect.ValueOf(fn)
	fnType := val.Type()
	inTypes, outTypes := []reflect.Type{}, []reflect.Type{}
	for i := 0; i < fnType.NumIn(); i++ {
		inTypes = append(inTypes, fnType.In(i))
	}
	for i := 0; i < fnType.NumOut(); i++ {
		outTypes = append(outTypes, fnType.Out(i))
	}

	if err := validateMain(inTypes, outTypes); err != nil {
		return mainFn{}, err
	}

	return mainFn{
		val:      val,
		inTypes:  inTypes,
		outTypes: outTypes,
	}, nil
}

func (fn mainFn) run(ctx *runctx) interface{} {
	inputs, err := ctx.generate(fn.inTypes)
	if err != nil {
		return nil
	}

	outs := fn.val.Call(inputs)
	for i := 0; i < len(outs); i++ {
		ctx.provide(outs[i].Interface())
	}

	if len(outs) == 0 {
		return nil
	} else if len(outs) == 1 && isError(fn.outTypes[0]) {
		return nil
	}
	return outs[0].Interface()
}
