package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elazarl/goproxy"
	"github.com/jasper-zsh/ones-hijacker-proxy/errors"
	"github.com/jasper-zsh/ones-hijacker-proxy/models"
	"github.com/jasper-zsh/ones-hijacker-proxy/types"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
	"time"
)

var _ goproxy.ReqHandler = (*ONESRequestHandler)(nil)

type AuthUpdatedCallback = func(binding *models.Binding)
type AuthExpiredCallback = func(binding *models.Binding)

type ONESRequestHandler struct {
	Instance    *models.Instance
	Account     *models.Account
	Binding     *models.Binding
	AuthUpdated AuthUpdatedCallback
	AuthExpired AuthExpiredCallback
	loginLock   sync.Mutex
}

func NewONESRequestHandler() *ONESRequestHandler {
	r := &ONESRequestHandler{}
	return r
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

			nReq, err := http.NewRequest(req.Method, fmt.Sprintf("%s/%s", h.Instance.BaseURL, matches[3]), req.Body)
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
				ctx.Warnf("Auth expired for %s", h.Account.Email)
				if h.AuthExpired != nil {
					h.AuthExpired(h.Binding)
				}
				h.Binding = nil
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
	if h.Binding == nil {
		err := h.Login(ctx)
		if err != nil {
			ctx.Warnf("Failed to login. %s", err.Error())
			return err
		}
		ctx.Warnf("Logged as %s", h.Account.Email)
	}
	req.Header.Set("Ones-Auth-Token", h.Binding.Token)
	req.AddCookie(&http.Cookie{
		Name:  "uid",
		Value: h.Binding.UserUUID,
	})
	req.AddCookie(&http.Cookie{
		Name:  "lt",
		Value: h.Binding.Token,
	})
	req.Header.Set("Ones-User-Id", h.Binding.UserUUID)
	return nil
}

func (h *ONESRequestHandler) Login(ctx *goproxy.ProxyCtx) error {
	h.loginLock.Lock()
	defer h.loginLock.Unlock()
	if h.Binding != nil {
		return nil
	}
	if ctx != nil {
		ctx.Warnf("Logging %s to  %s", h.Account.Email, h.Instance.BaseURL)
	}
	body, err := json.Marshal(types.LoginRequest{
		Email:    h.Account.Email,
		Password: h.Account.Password,
	})
	if err != nil {
		return err
	}
	loginUrl := fmt.Sprintf("%s/auth/login", h.Instance.BaseURL)
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
	h.Binding = &models.Binding{
		AccountID:  h.Account.ID,
		InstanceID: h.Instance.ID,
		UserUUID:   loginRes.User.UUID,
		Token:      loginRes.User.Token,
	}
	if h.AuthUpdated != nil {
		h.AuthUpdated(h.Binding)
	}

	return nil
}
