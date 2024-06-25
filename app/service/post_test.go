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
	"github.com/adamnasrudin03/go-template/pkg/helpers"
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
		Title:   "Title Test",
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
			Title:   "Test Title",
			Content: "test content",
			Tags:    []string{"tag1", "tag2"},
		},
		{
			ID:      102,
			Title:   "Test Title 102",
			Content: "test content 102",
			Tags:    []string{},
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

func (srv *PostServiceTestSuite) TestPostSrv_Create() {
	resp := &dto.PostRes{
		ID:      1,
		Title:   "title 1",
		Content: "content 2",
		Tags:    []string{"tags1", "tags2", "tags1"},
	}
	resp.CheckResp()

	tests := []struct {
		name     string
		req      dto.PostCreateReq
		mockFunc func(input dto.PostCreateReq)
		want     *dto.PostRes
		wantErr  bool
	}{
		{
			name: "invalid input",
			req: dto.PostCreateReq{
				Title:   "",
				Content: "",
				Tags:    []string{},
			},
			mockFunc: func(input dto.PostCreateReq) {
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed create db",
			req: dto.PostCreateReq{
				Title:   "title 1",
				Content: "content 2",
				Tags:    []string{"tags1", "tags2", "tags1"},
			},
			mockFunc: func(input dto.PostCreateReq) {
				input.Validate()
				srv.repo.On("Create", mock.Anything, input).Return(nil, errors.New("db error")).Once()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			req: dto.PostCreateReq{
				Title:   "title 1",
				Content: "content 2",
				Tags:    []string{"tags1", "tags2", "tags1"},
			},
			mockFunc: func(input dto.PostCreateReq) {
				input.Validate()
				srv.repo.On("Create", mock.Anything, input).Return(resp, nil).Once()
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

			got, err := srv.service.Create(srv.ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostSrv.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostSrv.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (srv *PostServiceTestSuite) TestPostSrv_DeleteByID() {

	tests := []struct {
		name     string
		postID   uint64
		mockFunc func(input uint64)
		wantErr  bool
	}{
		{
			name:   "err db",
			postID: 101,
			mockFunc: func(input uint64) {
				srv.repo.On("DeleteByID", mock.Anything, input).Return(errors.New("db error")).Once()
			},
			wantErr: true,
		},
		{
			name:   "Success",
			postID: 101,
			mockFunc: func(input uint64) {
				srv.repo.On("DeleteByID", mock.Anything, input).Return(nil).Once()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		srv.T().Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(tt.postID)
			}

			if err := srv.service.DeleteByID(srv.ctx, tt.postID); (err != nil) != tt.wantErr {
				t.Errorf("PostSrv.DeleteByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (srv *PostServiceTestSuite) TestPostSrv_UpdateByID() {

	tests := []struct {
		name     string
		req      dto.PostUpdateReq
		mockFunc func(input dto.PostUpdateReq)
		wantErr  bool
	}{
		{
			name: "invalid request",
			req: dto.PostUpdateReq{
				ID:      101,
				Title:   "",
				Content: "",
				Tags:    []string{},
			},
			mockFunc: func(input dto.PostUpdateReq) {
			},
			wantErr: true,
		},
		{
			name: "err db",
			req: dto.PostUpdateReq{
				ID:      101,
				Title:   "title test 101",
				Content: "content test 101",
				Tags:    []string{"tags1", "tags2"},
			},
			mockFunc: func(input dto.PostUpdateReq) {
				input.Validate()
				srv.repo.On("UpdateByID", mock.Anything, input).Return(helpers.ErrNotFound()).Once()
			},
			wantErr: true,
		},
		{
			name: "Success",
			req: dto.PostUpdateReq{
				ID:      101,
				Title:   "title test 101",
				Content: "content test 101",
				Tags:    []string{"tags1", "tags2"},
			},
			mockFunc: func(input dto.PostUpdateReq) {
				input.Validate()
				srv.repo.On("UpdateByID", mock.Anything, input).Return(nil).Once()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		srv.T().Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(tt.req)
			}

			if err := srv.service.UpdateByID(srv.ctx, tt.req); (err != nil) != tt.wantErr {
				t.Errorf("PostSrv.UpdateByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
