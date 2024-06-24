package controller

import (
	"github.com/adamnasrudin03/go-asset-findr/app/service"
	"github.com/sirupsen/logrus"
)

type PostController interface {
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
