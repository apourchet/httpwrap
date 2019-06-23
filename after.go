package httpwrap

import "reflect"

type afterFn struct {
	val      reflect.Value
	inTypes  []reflect.Type
	outTypes []reflect.Type
}

func newAfter(fn interface{}) (afterFn, error) {
	val := reflect.ValueOf(fn)
	fnType := val.Type()
	inTypes, outTypes := []reflect.Type{}, []reflect.Type{}
	for i := 0; i < fnType.NumIn(); i++ {
		inTypes = append(inTypes, fnType.In(i))
	}
	for i := 0; i < fnType.NumOut(); i++ {
		outTypes = append(outTypes, fnType.Out(i))
	}

	if err := validateAfter(inTypes, outTypes); err != nil {
		return afterFn{}, err
	}

	return afterFn{
		val:      val,
		inTypes:  inTypes,
		outTypes: outTypes,
	}, nil
}

func (fn afterFn) run(ctx *runctx) {
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
			ctx.provide(reflect.TypeOf(err), reflect.ValueOf(err))
		}
		inputs[i] = input
	}

	outs := fn.val.Call(inputs)
	for i := 0; i < len(outs); i++ {
		ctx.provide(fn.outTypes[i], outs[i])
	}
}
