package service

import (
	"context"

	"github.com/adamnasrudin03/go-asset-findr/app/configs"
	"github.com/adamnasrudin03/go-asset-findr/app/dto"
	"github.com/adamnasrudin03/go-asset-findr/app/repository"
	"github.com/adamnasrudin03/go-template/pkg/helpers"
	"github.com/sirupsen/logrus"
)

type PostService interface {
	GetList(ctx context.Context) ([]dto.PostRes, error)
	GetDetail(ctx context.Context, req dto.PostGetReq) (*dto.PostRes, error)
	Create(ctx context.Context, req dto.PostCreateReq) (*dto.PostRes, error)
	DeleteByID(ctx context.Context, postID uint64) error
}

type PostSrv struct {
	Repo   repository.PostRepository
	Cfg    *configs.Configs
	Logger *logrus.Logger
}

// NewPostService creates a new instance of PostService.
func NewPostService(
	postRepo repository.PostRepository,
	cfg *configs.Configs,
	logger *logrus.Logger,
) PostService {
	return &PostSrv{
		Repo:   postRepo,
		Cfg:    cfg,
		Logger: logger,
	}
}

func (srv *PostSrv) GetList(ctx context.Context) ([]dto.PostRes, error) {
	var (
		opName = "PostService-GetList"
		err    error
	)

	res, err := srv.Repo.GetAll(ctx)
	if err != nil {
		srv.Logger.Errorf("%s failed get data: %v \n", opName, err)
		return []dto.PostRes{}, helpers.ErrDB()
	}
	if len(res) == 0 {
		return []dto.PostRes{}, nil
	}

	for i := range res {
		res[i].CheckResp()
	}

	return res, nil
}

func (srv *PostSrv) GetDetail(ctx context.Context, req dto.PostGetReq) (*dto.PostRes, error) {
	var (
		opName = "PostService-GetDetail"
		err    error
	)

	err = req.Validate()
	if err != nil {
		return nil, err
	}

	res, err := srv.Repo.GetDetail(ctx, req)
	if err != nil {
		srv.Logger.Errorf("%s failed get data: %v \n", opName, err)
		return nil, helpers.ErrDB()
	}

	isExist := res != nil && res.ID > 0
	if !isExist {
		return nil, helpers.ErrNotFound()
	}

	res.CheckResp()
	return res, nil
}

func (srv *PostSrv) Create(ctx context.Context, req dto.PostCreateReq) (*dto.PostRes, error) {
	var (
		opName = "PostService-Create"
		err    error
	)

	err = req.Validate()
	if err != nil {
		return nil, err
	}

	result, err := srv.Repo.Create(ctx, req)
	if err != nil {
		srv.Logger.Errorf("%s failed create data: %v \n", opName, err)
		return nil, helpers.ErrCreatedDB()
	}
	result.CheckResp()
	return result, nil
}

func (srv *PostSrv) DeleteByID(ctx context.Context, postID uint64) error {
	var (
		opName = "PostService-DeleteByID"
		err    error
	)

	err = srv.Repo.DeleteByID(ctx, postID)
	if err != nil {
		srv.Logger.Errorf("%s failed delete data: %v \n", opName, err)
		return err
	}

	return nil
}
