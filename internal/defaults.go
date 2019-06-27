package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	// ErrValueNotFound is the error returned from the Get* functions
	// when this value was not found in the request.
	ErrValueNotFound = errors.New("value not found")
)

// DecodeBody uses a json decoder to decode the body of the request
// into the target object.
func DecodeBody(req *http.Request, obj interface{}) error {
	err := json.NewDecoder(req.Body).Decode(obj)
	if err == io.EOF {
		return nil
	}
	return err
}

// GetHeader returns the value of the header if it was found.
func GetHeader(req *http.Request, key string) (string, error) {
	val := req.Header.Get(key)
	if val == "" {
		return "", ErrValueNotFound
	}
	return val, nil
}

// GetSegment returns ErrValueNotFound because we have no way of
// knowing how the server mux has altered the request to let this
// information resurface.
func GetSegment(req *http.Request, key string) (string, error) {
	// We need some insight into the parsing of the request path if we
	// want access to this information. Therefore this must be
	// user-supplied.
	return "", ErrValueNotFound
}

// GetQueries returns the list of values that this query parameter
// had in the request.
func GetQueries(req *http.Request, key string) ([]string, error) {
	vals, found := req.URL.Query()[key]
	if !found || len(vals) == 0 {
		return nil, ErrValueNotFound
	}

	for i, val := range vals {
		val, err := url.QueryUnescape(val)
		if err != nil {
			return nil, fmt.Errorf("failed to unescape query value %s: %v", val, err)
		}
		vals[i] = val
	}
	return vals, nil
}

// GetCookie returns the value of the cookie by the name given.
func GetCookie(req *http.Request, key string) (string, error) {
	cookie, err := req.Cookie(key)
	if err == http.ErrNoCookie {
		return "", ErrValueNotFound
	} else if err != nil {
		return "", fmt.Errorf("failed to get cookie information from request: %v", err)
	}

	val := cookie.Value
	if unescaped, err := url.PathUnescape(cookie.Value); err == nil {
		return unescaped, nil
	}
	return val, nil
}
