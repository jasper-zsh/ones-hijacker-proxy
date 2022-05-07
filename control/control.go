package control

import (
	"github.com/jasper-zsh/ones-hijacker-proxy/handlers"
	"github.com/jasper-zsh/ones-hijacker-proxy/models"
	"github.com/jasper-zsh/ones-hijacker-proxy/types"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

type Control struct {
	app  *iris.Application
	db   *gorm.DB
	ones *handlers.ONESRequestHandler
}

func NewControl(db *gorm.DB, ones *handlers.ONESRequestHandler) *Control {
	c := &Control{
		app:  iris.New(),
		db:   db,
		ones: ones,
	}
	c.app.Get("/status", c.status)
	c.app.Get("/instances", c.listInstances)
	c.app.Post("/instances", c.createInstance)
	c.app.Delete("/instances/{id:uint}", c.deleteInstance)
	c.app.Post("/instances/{id:uint}/select", c.selectInstance)
	c.app.Get("/accounts", c.listAccounts)
	c.app.Post("/accounts", c.createAccount)
	c.app.Delete("/accounts/{id:uint}", c.deleteAccount)
	c.app.Post("/accounts/{id:uint}/select", c.selectAccount)
	return c
}

func (c *Control) SelectDefaultInstance() error {
	var instance models.Instance
	q := c.db.First(&instance)
	if q.Error != nil {
		return q.Error
	}
	c.ones.SetInstance(&instance)
	return nil
}

func (c *Control) SelectDefaultAccount() error {
	var account models.Account
	q := c.db.First(&account)
	if q.Error != nil {
		return q.Error
	}
	c.ones.SetAccount(&account)
	return nil
}

func (c *Control) Run() {
	c.app.Listen(":9090")
}

func (c *Control) listInstances(ctx iris.Context) {
	var instances []models.Instance
	q := c.db.Find(&instances)
	if q.Error != nil {
		ctx.StopWithError(500, q.Error)
		return
	}
	ctx.JSON(instances)
}

func (c *Control) createInstance(ctx iris.Context) {
	var instance models.Instance
	ctx.ReadJSON(&instance)

	c.db.Save(&instance)
	ctx.JSON(instance)
}

func (c *Control) deleteInstance(ctx iris.Context) {
	id, err := ctx.Params().GetUint("id")
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
	q := c.db.Delete(&models.Instance{}, id)
	if q.Error != nil {
		ctx.StopWithError(500, q.Error)
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
	var instance models.Instance
	q := c.db.First(&instance, id)
	if q.Error != nil {
		ctx.StopWithError(500, q.Error)
		return
	}
	c.ones.SetInstance(&instance)
	c.status(ctx)
}

func (c *Control) listAccounts(ctx iris.Context) {
	var accounts []models.Account
	q := c.db.Find(&accounts)
	if q.Error != nil {
		ctx.StopWithError(500, q.Error)
		return
	}

	ctx.JSON(accounts)
}

func (c *Control) createAccount(ctx iris.Context) {
	var account models.Account
	ctx.ReadJSON(&account)

	c.db.Save(&account)
	ctx.JSON(account)
}

func (c *Control) deleteAccount(ctx iris.Context) {
	id, err := ctx.Params().GetUint("id")
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}

	q := c.db.Delete(&models.Account{}, id)
	if q.Error != nil {
		ctx.StopWithError(500, q.Error)
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

	var account models.Account
	q := c.db.First(&account, id)
	if q.Error != nil {
		ctx.StopWithError(500, q.Error)
		return
	}

	c.ones.SetAccount(&account)
	c.ones.SetAuthUpdatedCallback(func(info *types.User) {
		account.Token = info.Token
		account.UserUUID = info.UUID
		c.db.Save(&account)
	})
	c.status(ctx)
}

func (c *Control) status(ctx iris.Context) {
	status := &types.StatusResponse{
		Account:  c.ones.Account(),
		Instance: c.ones.Instance(),
	}
	ctx.JSON(status)
}
