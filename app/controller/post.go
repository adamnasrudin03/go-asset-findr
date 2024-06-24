package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/adamnasrudin03/go-asset-findr/app/dto"
	"github.com/adamnasrudin03/go-asset-findr/app/service"
	"github.com/adamnasrudin03/go-template/pkg/helpers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PostController interface {
	GetList(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
}

type PostHandler struct {
	Service service.PostService
	Logger  *logrus.Logger
}

func NewPostDelivery(
	srv service.PostService,
	logger *logrus.Logger,
) PostController {
	return &PostHandler{
		Service: srv,
		Logger:  logger,
	}
}
func (c *PostHandler) GetList(ctx *gin.Context) {
	var (
		opName = "UserDelivery-GetList"
		resp   = []dto.PostRes{}
		err    error
	)

	resp, err = c.Service.GetList(ctx)
	if err != nil {
		c.Logger.Errorf("%v error: %v ", opName, err)
		helpers.RenderJSON(ctx.Writer, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *PostHandler) GetDetail(ctx *gin.Context) {
	var (
		opName  = "UserDelivery-GetDetail"
		idParam = strings.TrimSpace(ctx.Param("id"))
		err     error
	)

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.Logger.Errorf("%v error parse param: %v ", opName, err)
		helpers.RenderJSON(ctx.Writer, http.StatusBadRequest, helpers.ErrInvalid("ID Post", "Post ID"))
		return
	}

	res, err := c.Service.GetDetail(ctx, dto.PostGetReq{
		ID: id,
	})

	if err != nil {
		c.Logger.Errorf("%v error: %v ", opName, err)
		helpers.RenderJSON(ctx.Writer, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
