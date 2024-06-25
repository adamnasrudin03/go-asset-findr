package router

import (
	"github.com/adamnasrudin03/go-asset-findr/app/controller"
	"github.com/gin-gonic/gin"
)

func (r routes) postRouter(rg *gin.RouterGroup, handler controller.PostController) {
	post := rg.Group("/posts")
	{
		post.GET("/:id", handler.GetDetail)
		post.DELETE("/:id", handler.Delete)
		post.PUT("/:id", handler.Update)
		post.GET("", handler.GetList)
		post.POST("", handler.Create)
	}

}
