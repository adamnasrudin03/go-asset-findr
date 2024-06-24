package router

import (
	"net/http"

	"github.com/adamnasrudin03/go-asset-findr/app/controller"
	"github.com/adamnasrudin03/go-template/pkg/helpers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type routes struct {
	router *gin.Engine
}

func NewRoutes(h controller.Controllers) routes {
	var err error
	r := routes{
		router: gin.Default(),
	}

	r.router.Use(gin.Logger())
	r.router.Use(gin.Recovery())
	r.router.Use(cors.Default())

	r.router.GET("/", func(c *gin.Context) {
		helpers.RenderJSON(c.Writer, http.StatusOK, "welcome this server")
	})

	v1 := r.router.Group("/api")
	r.postRouter(v1, h.Post)

	r.router.NoRoute(func(c *gin.Context) {
		err = helpers.ErrRouteNotFound()
		helpers.RenderJSON(c.Writer, http.StatusNotFound, err)
	})
	return r
}

func (r routes) Run(addr string) error {
	return r.router.Run(addr)
}
