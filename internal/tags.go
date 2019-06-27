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
	if len(values) > 0 {
		return genVals(t, append([]string{value}, values...))
	}

	val := reflect.New(t)
	if err := json.Unmarshal([]byte(value), val.Interface()); err == nil {
		return val, nil
	}
	if err := json.Unmarshal([]byte(strconv.Quote(value)), val.Interface()); err != nil {
		return reflect.Zero(t), fmt.Errorf("failed to generate value for type %v and string content %v", t, value)
	}
	return val, nil
}

func genVals(t reflect.Type, values []string) (reflect.Value, error) {
	if t.Kind() != reflect.Slice {
		return reflect.Zero(t), fmt.Errorf("can only generate multiple string values into slice field")
	}
	val := reflect.New(t)

	joined := "[" + strings.Join(values, ",") + "]"
	if err := json.Unmarshal([]byte(joined), val.Interface()); err == nil {
		return val, nil
	}

	for i, strval := range values {
		values[i] = strconv.Quote(strval)
	}
	joined = "[" + strings.Join(values, ",") + "]"
	if err := json.Unmarshal([]byte(joined), val.Interface()); err != nil {
		return reflect.Zero(t), fmt.Errorf("failed to generate value for type %v and string content %v", t, values)
	}
	return val, nil
}
