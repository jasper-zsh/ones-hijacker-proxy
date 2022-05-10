package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jasper-zsh/ones-hijacker-proxy/errors"
	"github.com/jasper-zsh/ones-hijacker-proxy/models"
	"github.com/jasper-zsh/ones-hijacker-proxy/types"
	"gopkg.in/elazarl/goproxy.v1"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
	"time"
)

var _ goproxy.ReqHandler = (*ONESRequestHandler)(nil)

type AuthUpdatedCallback = func(info *types.User)

type ONESRequestHandler struct {
	instance    *models.Instance
	account     *models.Account
	authInfo    *types.User
	authUpdated AuthUpdatedCallback
	loginLock   sync.Mutex
}

func NewONESRequestHandler() *ONESRequestHandler {
	r := &ONESRequestHandler{}
	return r
}

func (h *ONESRequestHandler) SetInstance(instance *models.Instance) {
	h.instance = instance
}

func (h *ONESRequestHandler) SetAccount(account *models.Account) {
	h.account = account
}

func (h *ONESRequestHandler) SetAuthInfo(authInfo *types.User) {
	h.authInfo = authInfo
}

func (h *ONESRequestHandler) ClearAuthInfo() {
	h.authInfo = nil
}

func (h *ONESRequestHandler) SetAuthUpdatedCallback(cb AuthUpdatedCallback) {
	h.authUpdated = cb
}

func (h *ONESRequestHandler) Account() *models.Account {
	return h.account
}

func (h *ONESRequestHandler) Instance() *models.Instance {
	return h.instance
}

func (h *ONESRequestHandler) AuthInfo() *types.User {
	return h.authInfo
}

func (h *ONESRequestHandler) Handle(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	ctx.UserData = make([]types.Timing, 0)
	timing := func(period string, fun func()) {
		ts := time.Now().UnixMilli()
		fun()
		ctx.UserData = append(ctx.UserData.([]types.Timing), types.Timing{
			Period:   period,
			Duration: time.Now().UnixMilli() - ts,
		})
	}
	rule := regexp.MustCompile("/project/(.*?)/api/(.*?)/(.*)")
	if req.Host == "dev.myones.net" {
		matches := rule.FindStringSubmatch(req.RequestURI)
		if len(matches) > 0 {
			ctx.Warnf("Hijacking ONES %s API request %s", matches[2], req.RequestURI)

			nReq, err := http.NewRequest(req.Method, fmt.Sprintf("%s/%s", h.instance.BaseURL, matches[3]), req.Body)
			if resp := errors.ErrorResponse(req, err); resp != nil {
				return req, resp
			}
			nReq.Header = req.Header

			timing("inject_auth", func() {
				err = h.injectAuth(ctx, nReq)
			})
			if resp := errors.ErrorResponse(req, err); resp != nil {
				return req, resp
			}

			var resp *http.Response
			timing("request", func() {
				resp, err = http.DefaultClient.Do(nReq)
			})
			if resp.StatusCode == 401 {
				ctx.Warnf("Auth expired for %s", h.authInfo.Email)
				h.authInfo = nil
			}
			if resp := errors.ErrorResponse(req, err); resp != nil {
				return req, resp
			}

			for _, timing := range ctx.UserData.([]types.Timing) {
				resp.Header.Add("Server-Timing", fmt.Sprintf("%s;dur=%d", timing.Period, timing.Duration))
			}
			return req, resp
		}
	}
	return req, nil
}

func (h *ONESRequestHandler) injectAuth(ctx *goproxy.ProxyCtx, req *http.Request) error {
	if h.authInfo == nil {
		err := h.Login(ctx)
		if err != nil {
			ctx.Warnf("Failed to login. %s", err.Error())
			return err
		}
		ctx.Warnf("Logged as %s", h.authInfo.Email)
	}
	req.Header.Set("Ones-Auth-Token", h.authInfo.Token)
	req.AddCookie(&http.Cookie{
		Name:  "uid",
		Value: h.authInfo.UUID,
	})
	req.AddCookie(&http.Cookie{
		Name:  "lt",
		Value: h.authInfo.Token,
	})
	req.Header.Set("Ones-User-Id", h.authInfo.UUID)
	return nil
}

func (h *ONESRequestHandler) Login(ctx *goproxy.ProxyCtx) error {
	h.loginLock.Lock()
	defer h.loginLock.Unlock()
	if h.authInfo != nil {
		return nil
	}
	if ctx != nil {
		ctx.Warnf("Logging %s to  %s", h.account.Email, h.instance.BaseURL)
	}
	body, err := json.Marshal(types.LoginRequest{
		Email:    h.account.Email,
		Password: h.account.Password,
	})
	if err != nil {
		return err
	}
	loginUrl := fmt.Sprintf("%s/auth/login", h.instance.BaseURL)
	req, err := http.NewRequest("POST", loginUrl, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if ctx != nil {
		req.Header.Set("Referer", ctx.Req.Header.Get("Referer"))
	} else {
		req.Header.Set("Referer", "https://dev.myones.net/project/master/")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		switch resp.StatusCode {
		case 401:
			return fmt.Errorf("账号或密码错误")
		case 502:
			return fmt.Errorf("服务挂了")
		default:
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			return fmt.Errorf("error response: %s", body[0:30])
		}
	}
	respDecoder := json.NewDecoder(resp.Body)
	loginRes := types.LoginResponse{}
	err = respDecoder.Decode(&loginRes)
	if err != nil {
		return err
	}
	h.authInfo = loginRes.User
	if h.authUpdated != nil {
		h.authUpdated(loginRes.User)
	}

	return nil
}
