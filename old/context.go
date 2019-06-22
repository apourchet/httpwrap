package httpwrap

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Context contains the raw http request and response interfaces, as
// well as all of the previous results.
type Context struct {
	Res http.ResponseWriter
	Req *http.Request

	Results []Result

	bodies [][]byte
}

func newContext(w http.ResponseWriter, req *http.Request) (Context, error) {
	body, err := ioutil.ReadAll(req.Body)
	return Context{
		Res:    w,
		Req:    req,
		bodies: [][]byte{body},
	}, err
}

func (ctx Context) unpack(in interface{}) error {
	for _, body := range ctx.bodies {
		if len(body) == 0 {
			continue
		}
		if err := json.Unmarshal(body, in); err != nil {
			return err
		}
	}
	return nil
}
