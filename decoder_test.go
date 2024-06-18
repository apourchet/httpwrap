package httpwrap

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

type holder struct {
	Body1    int     `json:"body1"`
	Query1   string  `http:"query=query1"`
	Query2   []int   `http:"query=query2"`
	Header1  *string `http:"header=header1"`
	Cookie1  bool    `http:"cookie=cookie1"`
	Segment1 string  `http:"segment=segment1"`
}

func TestDecoder(t *testing.T) {
	t.Run("decode full", func(t *testing.T) {
		body := strings.NewReader(`{"body1":42}`)
		url := fmt.Sprintf("http://localhost/path?query1=query1val&query2=1&query2=2")
		req := httptest.NewRequest("POST", url, body)
		req.Header.Set("header1", "header1val")
		req.AddCookie(&http.Cookie{Name: "cookie1", Value: "true"})
		req.SetPathValue("segment1", "segment1val")

		decoder := NewDecoder()
		into := holder{}
		err := decoder.Decode(req, &into)
		require.NoError(t, err)

		_headerVal := "header1val"
		expected := holder{
			Body1:    42,
			Header1:  &_headerVal,
			Segment1: "segment1val",
			Cookie1:  true,
			Query1:   "query1val",
			Query2:   []int{1, 2},
		}
		require.Equal(t, expected, into)
	})
}

func TestDecoderGorillaMux(t *testing.T) {
	t.Run("decode full", func(t *testing.T) {
		body := strings.NewReader(`{"body1":42}`)
		url := fmt.Sprintf("http://localhost/path?query1=query1val&query2=1&query2=2")
		req := httptest.NewRequest("POST", url, body)
		req.Header.Set("header1", "header1val")
		req.AddCookie(&http.Cookie{Name: "cookie1", Value: "true"})
		req = mux.SetURLVars(req, map[string]string{"segment1": "segment1val"})

		decoder := NewDecoder()
		into := holder{}
		err := decoder.Decode(req, &into)
		require.NoError(t, err)

		_headerVal := "header1val"
		expected := holder{
			Body1:    42,
			Header1:  &_headerVal,
			Segment1: "segment1val",
			Cookie1:  true,
			Query1:   "query1val",
			Query2:   []int{1, 2},
		}
		require.Equal(t, expected, into)
	})
}
