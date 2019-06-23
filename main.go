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
	inputs := make([]reflect.Value, len(fn.inTypes))
	for i, inType := range fn.inTypes {
		if val, found := ctx.get(inType); found {
			inputs[i] = val
			continue
		}
		input, err := ctx.construct(inType)
		if err != nil {
			ctx.provide(reflect.TypeOf(err), reflect.ValueOf(err))
			return nil
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

	if len(outs) == 1 && fn.outTypes[0] == _errorType {
		return nil
	}
	return outs[0].Interface()
}
