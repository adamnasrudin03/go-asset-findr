package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/adamnasrudin03/go-asset-findr/app/configs"
	"github.com/adamnasrudin03/go-asset-findr/app/dto"
	"github.com/adamnasrudin03/go-asset-findr/app/repository/mocks"
	"github.com/adamnasrudin03/go-asset-findr/pkg/driver"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PostServiceTestSuite struct {
	suite.Suite
	repo    *mocks.PostRepository
	ctx     context.Context
	service PostService
}

func (srv *PostServiceTestSuite) SetupTest() {
	var (
		cfg    = configs.GetInstance()
		logger = driver.Logger(cfg)
	)

	srv.repo = &mocks.PostRepository{}
	srv.ctx = context.Background()
	srv.service = NewPostService(srv.repo, cfg, logger)
}

func TestPostService(t *testing.T) {
	suite.Run(t, new(PostServiceTestSuite))
}

func (srv *PostServiceTestSuite) TestPostSrv_GetDetail() {
	resp := &dto.PostRes{
		ID:      101,
		Title:   "test title",
		Content: "test content",
		Tags:    []string{"tag1", "tag2"},
	}

	tests := []struct {
		name     string
		req      dto.PostGetReq
		mockFunc func(input dto.PostGetReq)
		want     *dto.PostRes
		wantErr  bool
	}{
		{
			name: "validation error",
			req: dto.PostGetReq{
				ID: 0,
			},
			mockFunc: func(input dto.PostGetReq) {
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed query db",
			req: dto.PostGetReq{
				ID: 101,
			},
			mockFunc: func(input dto.PostGetReq) {
				srv.repo.On("GetDetail", mock.Anything, input).Return(nil, errors.New("db error")).Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not found",
			req: dto.PostGetReq{
				ID: 101,
			},
			mockFunc: func(input dto.PostGetReq) {
				srv.repo.On("GetDetail", mock.Anything, input).Return(nil, nil).Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			req: dto.PostGetReq{
				ID: 101,
			},
			mockFunc: func(input dto.PostGetReq) {
				srv.repo.On("GetDetail", mock.Anything, input).Return(resp, nil).Once()
			},
			want:    resp,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		srv.T().Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(tt.req)
			}
			got, err := srv.service.GetDetail(srv.ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostSrv.GetDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostSrv.GetDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (srv *PostServiceTestSuite) TestPostSrv_GetList() {
	resp := []dto.PostRes{
		{
			ID:      101,
			Title:   "test title",
			Content: "test content",
			Tags:    []string{"tag1", "tag2"},
		},
	}

	tests := []struct {
		name     string
		mockFunc func()
		want     []dto.PostRes
		wantErr  bool
	}{
		{
			name: "error query db",
			mockFunc: func() {
				srv.repo.On("GetAll", mock.Anything).Return([]dto.PostRes{}, errors.New("db error")).Once()
			},
			want:    []dto.PostRes{},
			wantErr: true,
		},
		{
			name: "empty data",
			mockFunc: func() {
				srv.repo.On("GetAll", mock.Anything).Return([]dto.PostRes{}, nil).Once()
			},
			want:    []dto.PostRes{},
			wantErr: false,
		},
		{
			name: "Success",
			mockFunc: func() {
				srv.repo.On("GetAll", mock.Anything).Return(resp, nil).Once()
			},
			want:    resp,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		srv.T().Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc()
			}
			got, err := srv.service.GetList(srv.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostSrv.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostSrv.GetList() = %v, want %v", got, tt.want)
			}
		})
	}
}
