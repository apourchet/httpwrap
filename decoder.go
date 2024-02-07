package httpwrap

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/apourchet/httpwrap/defaults"
)

// Decoder is a struct that allows for the decoding of http requests
// into arbitrary objects.
type Decoder struct {
	// DecodeBody is the function that will be used to decode the
	// request body into a target object.
	DecodeBody DecodeFunc

	// Header is the function used to get the string value of a header.
	Header func(*http.Request, string) (string, error)

	// Segment is the function used to get the string value of a path
	// parameter.
	Segment func(*http.Request, string) (string, error)

	// Queries is the function used to get the string values of a query
	// parameter.
	Queries func(*http.Request, string) ([]string, error)

	// Cookie is the function used to get the value of a cookie from a
	// request.
	Cookie func(*http.Request, string) (string, error)
}

// DecodeFunc is the function signature for decoding a request into an
// object.
type DecodeFunc func(req *http.Request, obj any) error

// NewDecoder returns a new decoder with sensible defaults for the
// DecodeBody, Header and Query functions.
// By default, it uses a JSON decoder on the request body.
func NewDecoder() *Decoder {
	return &Decoder{
		DecodeBody: defaults.DecodeBody,
		Header:     defaults.GetHeader,
		Segment:    defaults.GetSegment,
		Queries:    defaults.GetQueries,
		Cookie:     defaults.GetCookie,
	}
}

// Decode will (by default), given a struct definition:
//
//	type Request struct {
//			AuthString string    `http:"header=Authorization"`
//			Limit int            `http:"query=limit"`
//			Resource string      `http:"segment=resource"`
//			UserCookie float64   `http:"cookie=user_cookie"`
//			Extra map[string]int `json:"extra"`
//	}
//
// The Authorization header will be parsed into the field Token of the
// request struct.
//
// The Limit field will come from the query string.
//
// The Resource field will come from the resource value of the path (e.g: /api/pets/{resource}).
//
// The Extra field will come from deserializing the request body from JSON encoding.
func (d *Decoder) Decode(req *http.Request, obj any) error {
	if err := d.DecodeBody(req, obj); err != nil {
		return err
	}

	v, valid := defaults.DerefValue(obj)
	if !valid || v.Kind() != reflect.Struct {
		return nil
	}

	t, valid := defaults.DerefType(obj)
	if !valid {
		return nil
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		directive, found := field.Tag.Lookup("http")
		if !found || directive == "" {
			continue
		}

		f := v.Field(i)
		if !f.IsValid() {
			return fmt.Errorf("field %s is not valid to decode into from request", field.Name)
		} else if !f.CanSet() {
			continue
		}

		if err := d.decodeDirective(req, f, directive); err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) decodeDirective(req *http.Request, field reflect.Value, directive string) error {
	split := strings.SplitN(directive, "=", 2)
	if len(split) != 2 {
		return fmt.Errorf("malformed http struct tag: %v", directive)
	}

	tagkey, tagval := split[0], split[1]
	return d.decodeValue(req, field, tagkey, tagval)
}

func (d *Decoder) decodeValue(req *http.Request, field reflect.Value, tagkey, tagval string) error {
	strvals := []string{""}
	var err error

	switch tagkey {
	case "header":
		strvals[0], err = d.Header(req, tagval)
	case "segment":
		strvals[0], err = d.Segment(req, tagval)
	case "cookie":
		strvals[0], err = d.Cookie(req, tagval)
	case "query":
		strvals, err = d.Queries(req, tagval)
	default:
		return fmt.Errorf("unrecognized http tag %v", tagkey)
	}

	if len(strvals) == 0 {
		return nil
	}

	if err == defaults.ErrValueNotFound {
		return nil
	} else if err != nil {
		return err
	}

	val, err := defaults.GenVal(field.Type(), strvals[0], strvals[1:]...)
	if err != nil {
		return err
	}

	field.Set(val)
	return nil
}
