package main

import (
	"github.com/elazarl/goproxy"
	"github.com/jasper-zsh/ones-hijacker-proxy/cert"
	"github.com/jasper-zsh/ones-hijacker-proxy/control"
	"github.com/jasper-zsh/ones-hijacker-proxy/dao"
	"github.com/jasper-zsh/ones-hijacker-proxy/errors"
	"github.com/jasper-zsh/ones-hijacker-proxy/handlers"
	"github.com/jasper-zsh/ones-hijacker-proxy/services"
	"go.uber.org/dig"
	"net/http"
	_ "net/http/pprof"
	"regexp"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.CertStore = cert.NewCertStorage()
	c := dig.New()
	dao.Provide(c)
	services.Provide(c)
	handlers.Provide(c)
	control.Provide(c)
	err := c.Invoke(func(ctl *control.Control, ones *handlers.ONESRequestHandler) {
		err := ctl.SelectDefaultInstance()
		errors.OrPanic(err)
		err = ctl.SelectDefaultAccount()
		errors.OrPanic(err)
		proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("dev\\.myones\\.net"))).
			HandleConnect(goproxy.AlwaysMitm)
		proxy.OnRequest().Do(ones)

		proxy.Verbose = false
		go func() {
			err := http.ListenAndServe(":6789", proxy)
			errors.OrPanic(err)
		}()
		ctl.Run()
	})
	errors.OrPanic(err)
}
