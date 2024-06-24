package app

import (
	"github.com/adamnasrudin03/go-asset-findr/app/configs"
	"github.com/adamnasrudin03/go-asset-findr/app/controller"
	"github.com/adamnasrudin03/go-asset-findr/app/repository"
	"github.com/adamnasrudin03/go-asset-findr/app/service"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func WiringRepository(db *gorm.DB, cfg *configs.Configs, logger *logrus.Logger) *repository.Repositories {
	return &repository.Repositories{
		Post: repository.NewPostRepository(db, cfg, logger),
	}
}

func WiringService(repo *repository.Repositories, cfg *configs.Configs, logger *logrus.Logger) *service.Services {
	return &service.Services{
		Post: service.NewPostService(repo.Post, cfg, logger),
	}
}

func WiringController(srv *service.Services, cfg *configs.Configs, logger *logrus.Logger) *controller.Controllers {
	return &controller.Controllers{
		Post: controller.NewPostDelivery(srv.Post, logger),
	}
}
