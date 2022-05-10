package errors

import (
	"github.com/elazarl/goproxy"
	"net/http"
)

func OrPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func ErrorResponse(req *http.Request, err error) *http.Response {
	if err != nil {
		return goproxy.NewResponse(req, "text/plain", 500, err.Error())
	}
	return nil
}
