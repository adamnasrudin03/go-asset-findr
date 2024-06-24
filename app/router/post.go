package router

import (
	"net/http"

	"github.com/adamnasrudin03/go-asset-findr/app/controller"
	"github.com/adamnasrudin03/go-template/pkg/helpers"
	"github.com/gin-gonic/gin"
)

func (r routes) postRouter(rg *gin.RouterGroup, handler controller.PostController) {
	tm := rg.Group("/posts")
	{
		tm.GET("/", func(c *gin.Context) {
			helpers.RenderJSON(c.Writer, http.StatusOK, "Build with love by adamnasrudin03")
		})
	}

}
