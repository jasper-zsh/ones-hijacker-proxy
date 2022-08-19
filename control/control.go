package control

import (
	"github.com/jasper-zsh/ones-hijacker-proxy/handlers"
	"github.com/jasper-zsh/ones-hijacker-proxy/models"
	"github.com/jasper-zsh/ones-hijacker-proxy/services"
	"github.com/jasper-zsh/ones-hijacker-proxy/types"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/cors"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type ControlDeps struct {
	dig.In

	DB              *gorm.DB
	Handler         *handlers.ONESRequestHandler
	AccountService  *services.AccountService
	InstanceService *services.InstanceService
}

type Control struct {
	deps ControlDeps
	app  *iris.Application
}

func NewControl(deps ControlDeps) *Control {
	return &Control{
		deps: deps,
		app:  iris.New(),
	}
}

func (c *Control) Run() {
	c.app.UseRouter(cors.New().Handler())
	c.app.Get("/status", c.status)
	c.app.Get("/instances", c.listInstances)
	c.app.Post("/instances", c.createInstance)
	c.app.Post("/instances/{id:uint}", c.updateInstance)
	c.app.Delete("/instances/{id:uint}", c.deleteInstance)
	c.app.Post("/instances/{id:uint}/select", c.selectInstance)
	c.app.Get("/accounts", c.listAccounts)
	c.app.Post("/accounts", c.createAccount)
	c.app.Post("/accounts/{id:uint}", c.updateAccount)
	c.app.Delete("/accounts/{id:uint}", c.deleteAccount)
	c.app.Post("/accounts/{id:uint}/select", c.selectAccount)
	c.app.Listen(":9090")
}

func (c *Control) SelectDefaultInstance() error {
	var instance models.Instance
	q := c.deps.DB.First(&instance)
	if q.Error != nil {
		return q.Error
	}
	c.deps.Handler.Instance = &instance
	return nil
}

func (c *Control) SelectDefaultAccount() error {
	var account models.Account
	q := c.deps.DB.First(&account)
	if q.Error != nil {
		return q.Error
	}
	c.deps.Handler.Account = &account
	return nil
}

func (c *Control) listInstances(ctx iris.Context) {
	instances, err := c.deps.InstanceService.ListInstances()
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
	ctx.JSON(instances)
}

func (c *Control) createInstance(ctx iris.Context) {
	var instance models.Instance
	ctx.ReadJSON(&instance)

	err := c.deps.InstanceService.SaveInstance(&instance)
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
	ctx.JSON(instance)
}

func (c *Control) updateInstance(ctx iris.Context) {
	id, err := ctx.Params().GetUint("id")
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
	var instance models.Instance
	err = ctx.ReadJSON(&instance)
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
	instance.ID = id
	err = c.deps.InstanceService.SaveInstance(&instance)
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}

	ctx.JSON(instance)
}

func (c *Control) deleteInstance(ctx iris.Context) {
	id, err := ctx.Params().GetUint("id")
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
	err = c.deps.InstanceService.DeleteInstance(id)
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
	ctx.StopWithStatus(204)
}

func (c *Control) selectInstance(ctx iris.Context) {
	id, err := ctx.Params().GetUint("id")
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}

	err = c.deps.InstanceService.SelectInstance(id)
	resp := &types.StatusResponse{
		Account:  c.deps.Handler.Account,
		Instance: c.deps.Handler.Instance,
	}
	if err != nil {
		resp.ErrorMsg = err.Error()
	}

	ctx.JSON(resp)
}

func (c *Control) listAccounts(ctx iris.Context) {
	accounts, err := c.deps.AccountService.ListAccounts()
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}

	ctx.JSON(accounts)
}

func (c *Control) createAccount(ctx iris.Context) {
	var account models.Account
	ctx.ReadJSON(&account)

	c.deps.AccountService.SaveAccount(&account)
	ctx.JSON(account)
}

func (c *Control) updateAccount(ctx iris.Context) {
	id, err := ctx.Params().GetUint("id")
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
	var account models.Account
	err = ctx.ReadJSON(&account)
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
	account.ID = id
	err = c.deps.AccountService.SaveAccount(&account)
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
	ctx.JSON(account)
}

func (c *Control) deleteAccount(ctx iris.Context) {
	id, err := ctx.Params().GetUint("id")
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}

	err = c.deps.AccountService.DeleteAccount(id)
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
	ctx.StopWithStatus(204)
}

func (c *Control) selectAccount(ctx iris.Context) {
	id, err := ctx.Params().GetUint("id")
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}

	err = c.deps.AccountService.SelectAccount(id)
	resp := &types.StatusResponse{
		Account:  c.deps.Handler.Account,
		Instance: c.deps.Handler.Instance,
	}
	if err != nil {
		resp.ErrorMsg = err.Error()
	}
	ctx.JSON(resp)
}

func (c *Control) status(ctx iris.Context) {
	status := &types.StatusResponse{
		Account:  c.deps.Handler.Account,
		Instance: c.deps.Handler.Instance,
	}
	ctx.JSON(status)
}
