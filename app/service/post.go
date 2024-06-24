package service

import (
	"github.com/adamnasrudin03/go-asset-findr/app/configs"
	"github.com/adamnasrudin03/go-asset-findr/app/repository"
	"github.com/sirupsen/logrus"
)

type PostService interface {
}

type PostSrv struct {
	Repo   repository.PostRepository
	Cfg    *configs.Configs
	Logger *logrus.Logger
}

func NewPostService(
	tmRepo repository.PostRepository,
	cfg *configs.Configs,
	logger *logrus.Logger,
) PostService {
	return PostSrv{
		Repo:   tmRepo,
		Cfg:    cfg,
		Logger: logger,
	}
}
