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
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
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
		opName = "PostController-GetList"
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
		opName  = "PostController-GetDetail"
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

func (c *PostHandler) Create(ctx *gin.Context) {
	var (
		opName = "PostController-Create"
		input  dto.PostCreateReq
		err    error
	)

	err = ctx.ShouldBindJSON(&input)
	if err != nil {
		c.Logger.Errorf("%v error bind json: %v ", opName, err)
		helpers.RenderJSON(ctx.Writer, http.StatusBadRequest, helpers.ErrGetRequest())
		return
	}

	res, err := c.Service.Create(ctx, input)
	if err != nil {
		c.Logger.Errorf("%v error: %v ", opName, err)
		helpers.RenderJSON(ctx.Writer, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (c *PostHandler) Delete(ctx *gin.Context) {
	var (
		opName  = "PostController-Delete"
		idParam = strings.TrimSpace(ctx.Param("id"))
		err     error
	)

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.Logger.Errorf("%v error parse param: %v ", opName, err)
		helpers.RenderJSON(ctx.Writer, http.StatusBadRequest, helpers.ErrInvalid("ID Post", "Post ID"))
		return
	}

	err = c.Service.DeleteByID(ctx, id)
	if err != nil {
		c.Logger.Errorf("%v error: %v ", opName, err)
		helpers.RenderJSON(ctx.Writer, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.ResponseMessage{Message: "Deleted data post successfully"})
}

func (c *PostHandler) Update(ctx *gin.Context) {
	var (
		opName  = "PostController-Update"
		idParam = strings.TrimSpace(ctx.Param("id"))
		input   dto.PostUpdateReq
		err     error
	)

	err = ctx.ShouldBindJSON(&input)
	if err != nil {
		c.Logger.Errorf("%v error bind json: %v ", opName, err)
		helpers.RenderJSON(ctx.Writer, http.StatusBadRequest, helpers.ErrGetRequest())
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.Logger.Errorf("%v error parse param: %v ", opName, err)
		helpers.RenderJSON(ctx.Writer, http.StatusBadRequest, helpers.ErrInvalid("ID Post", "Post ID"))
		return
	}

	input.ID = id
	err = c.Service.UpdateByID(ctx, input)
	if err != nil {
		c.Logger.Errorf("%v error: %v ", opName, err)
		helpers.RenderJSON(ctx.Writer, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, dto.ResponseMessage{Message: "Updated data post successfully"})
}
