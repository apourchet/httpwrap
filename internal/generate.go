package internal

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// GenVal generates an interface{} from the string values given.
func GenVal(t reflect.Type, value string, values ...string) (reflect.Value, error) {
	if len(values) > 0 || t.Kind() == reflect.Slice {
		return genVals(t, append([]string{value}, values...))
	}

	val := reflect.New(t)
	if err := json.Unmarshal([]byte(value), val.Interface()); err == nil {
		return val.Elem(), nil
	}
	if err := json.Unmarshal([]byte(strconv.Quote(value)), val.Interface()); err != nil {
		return reflect.Zero(t), fmt.Errorf("failed to generate value for type %v and string content %v", t, value)
	}
	return val.Elem(), nil
}

func genVals(t reflect.Type, values []string) (reflect.Value, error) {
	val := reflect.New(t)

	joined := "[" + strings.Join(values, ",") + "]"
	if err := json.Unmarshal([]byte(joined), val.Interface()); err == nil {
		return val.Elem(), nil
	}

	for i, strval := range values {
		values[i] = strconv.Quote(strval)
	}
	joined = "[" + strings.Join(values, ",") + "]"
	if err := json.Unmarshal([]byte(joined), val.Interface()); err != nil {
		msg := "failed to generate value for type %v and string content %v: %v"
		return reflect.Zero(t), fmt.Errorf(msg, t, values, err)
	}
	return val.Elem(), nil
}
