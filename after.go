package httpwrap

import "reflect"

type afterFn struct {
	val      reflect.Value
	inTypes  []reflect.Type
	outTypes []reflect.Type
}

func newAfter(fn any) (afterFn, error) {
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
	inputs, err := ctx.generate(fn.inTypes)
	if err != nil {
		return
	}

	outs := fn.val.Call(inputs)
	for i := 0; i < len(outs); i++ {
		ctx.provide(outs[i].Interface())
	}
}
