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
	inputs := make([]reflect.Value, len(fn.inTypes))
	for i, inType := range fn.inTypes {
		if isEmptyInterface(inType) {
			inputs[i] = ctx.response
			continue
		} else if val, found := ctx.get(inType); found {
			inputs[i] = val
			continue
		}

		input, err := ctx.construct(inType)
		if err != nil {
			ctx.provide(_errorType, reflect.ValueOf(err))
			return err
		}
		inputs[i] = input
	}

	outs := fn.val.Call(inputs)
	if len(outs) == 0 {
		return nil
	}

	for i := 0; i < len(outs); i++ {
		ctx.provide(fn.outTypes[i], outs[i])
	}

	lastVal := outs[len(outs)-1]
	if lastVal.IsNil() {
		return nil
	}
	return lastVal.Interface().(error)
}
