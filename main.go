package main

import (
	"github.com/elazarl/goproxy"
	"github.com/jasper-zsh/ones-hijacker-proxy/cert"
	"github.com/jasper-zsh/ones-hijacker-proxy/control"
	"github.com/jasper-zsh/ones-hijacker-proxy/dao"
	"github.com/jasper-zsh/ones-hijacker-proxy/errors"
	"github.com/jasper-zsh/ones-hijacker-proxy/handlers"
	"github.com/jasper-zsh/ones-hijacker-proxy/models"
	"github.com/jasper-zsh/ones-hijacker-proxy/services"
	"go.uber.org/dig"
	"gorm.io/gorm"
	"net/http"
	_ "net/http/pprof"
	"regexp"
)

type ReqHandlerWrapper struct {
	*handlers.ONESRequestHandler
}

func (r ReqHandlerWrapper) Handle(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return r.FilterReq(req, ctx)
}

var _ goproxy.ReqHandler = (*ReqHandlerWrapper)(nil)

type RespHandlerWrapper struct {
	*handlers.ONESRequestHandler
}

func (r RespHandlerWrapper) Handle(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	return r.FilterResp(resp, ctx)
}

var _ goproxy.RespHandler = (*RespHandlerWrapper)(nil)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.CertStore = cert.NewCertStorage()
	c := dig.New()
	dao.Provide(c)
	services.Provide(c)
	handlers.Provide(c)
	control.Provide(c)
	err := c.Invoke(func(db *gorm.DB, ctl *control.Control, ones *handlers.ONESRequestHandler) {
		ones.AuthUpdated = func(binding *models.Binding) {
			db.Save(binding)
		}
		ones.AuthExpired = func(binding *models.Binding) {
			db.Delete(binding)
		}
		err := ctl.SelectDefaultInstance()
		errors.OrPanic(err)
		err = ctl.SelectDefaultAccount()
		errors.OrPanic(err)
		err = ctl.LoadBinding()
		errors.OrPanic(err)
		proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("dev\\.myones\\.net"))).
			HandleConnect(goproxy.AlwaysMitm)
		proxy.OnRequest().Do(ReqHandlerWrapper{ones})
		proxy.OnResponse().Do(RespHandlerWrapper{ones})

		proxy.Verbose = false
		go func() {
			err := http.ListenAndServe(":6789", proxy)
			errors.OrPanic(err)
		}()
		ctl.Run()
	})
	errors.OrPanic(err)
}
