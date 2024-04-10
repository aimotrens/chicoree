package noroute

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	staticFiles fs.FS
	httpFS      http.FileSystem
}

func NewController(fs fs.FS) *Controller {
	return &Controller{staticFiles: fs, httpFS: http.FS(fs)}
}

func (c *Controller) RegisterRoutes(g *gin.Engine) {
	g.NoRoute(func(ctx *gin.Context) {
		s, err := fs.Stat(c.staticFiles, path.Join("static"+ctx.Request.URL.Path))

		if ctx.Request.URL.Path == "/" || errors.Is(err, os.ErrNotExist) {
			ctx.FileFromFS("static/index.htm", c.httpFS)
			return
		}

		if s.IsDir() {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.FileFromFS("static"+ctx.Request.URL.Path, c.httpFS)
	})
}
