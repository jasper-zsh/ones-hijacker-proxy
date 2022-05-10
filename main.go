package main

import (
	"github.com/jasper-zsh/ones-hijacker-proxy/cert"
	"github.com/jasper-zsh/ones-hijacker-proxy/control"
	"github.com/jasper-zsh/ones-hijacker-proxy/dao"
	"github.com/jasper-zsh/ones-hijacker-proxy/errors"
	"github.com/jasper-zsh/ones-hijacker-proxy/handlers"
	"github.com/jasper-zsh/ones-hijacker-proxy/services"
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"
	_ "net/http/pprof"
	"regexp"
)

func main() {
	db, err := dao.InitDB()
	if err != nil {
		panic(err)
	}
	proxy := goproxy.NewProxyHttpServer()
	proxy.CertStore = cert.NewCertStorage()
	ones := handlers.NewONESRequestHandler()
	ctl := control.NewControl(
		db,
		ones,
		services.NewAccountService(db, ones),
		services.NewInstanceService(db, ones),
	)
	err = ctl.SelectDefaultInstance()
	errors.OrPanic(err)
	err = ctl.SelectDefaultAccount()
	errors.OrPanic(err)
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("dev\\.myones\\.net"))).
		HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().Do(ones)

	proxy.Verbose = false
	go http.ListenAndServe(":6789", proxy)
	go http.ListenAndServe(":22222", nil)
	ctl.Run()
}
