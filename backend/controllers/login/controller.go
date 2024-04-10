package login

import (
	"chicoree/controllers"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	controllers.ControllerBase
}

func NewController(ecc controllers.EntClientConstructor) *Controller {
	return &Controller{ControllerBase: controllers.NewBaseController(ecc)}
}

func (c *Controller) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/authenticate", c.authenticate)
}
