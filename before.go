package httpwrap

import "reflect"

type beforeFn struct {
	val      reflect.Value
	inTypes  []reflect.Type
	outTypes []reflect.Type
}

func newBefore(fn interface{}) (beforeFn, error) {
	val := reflect.ValueOf(fn)
	fnType := val.Type()
	inTypes, outTypes := []reflect.Type{}, []reflect.Type{}
	for i := 0; i < fnType.NumIn(); i++ {
		inTypes = append(inTypes, fnType.In(i))
	}
	for i := 0; i < fnType.NumOut(); i++ {
		outTypes = append(outTypes, fnType.Out(i))
	}

	if err := validateBefore(inTypes, outTypes); err != nil {
		return beforeFn{}, err
	}

	return beforeFn{
		val:      val,
		inTypes:  inTypes,
		outTypes: outTypes,
	}, nil
}

func (fn beforeFn) run(ctx *runctx) error {
	inputs, err := ctx.generate(fn.inTypes)
	if err != nil {
		return err
	}

	outs := fn.val.Call(inputs)
	if len(outs) == 0 {
		return nil
	}

	for i := 0; i < len(outs); i++ {
		ctx.provide(outs[i].Interface())
	}

	if !isError(fn.outTypes[len(outs)-1]) {
		return nil
	}

	lastVal := outs[len(outs)-1]
	if lastVal.IsNil() {
		return nil
	}
	return lastVal.Interface().(error)
}
