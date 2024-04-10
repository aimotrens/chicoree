package login

import (
	"chicoree/ent"

	"github.com/gin-gonic/gin"
)

func (c *Controller) authenticate(ctx *gin.Context) {
	c.WithTx(ctx, func(tx *ent.Tx) error {
		// do something
		return nil
	})
}
