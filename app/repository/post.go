package repository

import (
	"context"
	"errors"
	"sort"

	"github.com/adamnasrudin03/go-asset-findr/app/configs"
	"github.com/adamnasrudin03/go-asset-findr/app/dto"
	"github.com/adamnasrudin03/go-asset-findr/app/models"
	"github.com/adamnasrudin03/go-template/pkg/helpers"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PostRepository interface {
	GetAll(ctx context.Context) (result []dto.PostRes, err error)
	GetDetail(ctx context.Context, req dto.PostGetReq) (*dto.PostRes, error)
	GetDetailTag(ctx context.Context, req dto.TagGetReq) (*models.Tag, error)
}

type PostRepo struct {
	DB     *gorm.DB
	Cfg    *configs.Configs
	Logger *logrus.Logger
}

func NewPostRepository(
	db *gorm.DB,
	cfg *configs.Configs,
	logger *logrus.Logger,
) PostRepository {
	return &PostRepo{
		DB:     db,
		Cfg:    cfg,
		Logger: logger,
	}
}

func (r *PostRepo) findTags(ctx context.Context, postID uint64) (result []string, err error) {
	var (
		opName = "PostRepository-findTags"
		tags   = []models.Tag{}
	)

	err = r.DB.WithContext(ctx).
		Raw("SELECT tag.id, tag.label FROM post_tag "+
			" INNER JOIN tag ON tag.id = post_tag.tag_id "+
			" WHERE post_tag.post_id = ?", postID).
		Scan(&tags).Error
	if err != nil {
		r.Logger.Errorf("%s failed get data: %v \n", opName, err)
		return []string{}, err
	}

	sort.SliceStable(tags, func(i, j int) bool {
		return tags[i].Label < tags[j].Label
	})

	for _, v := range tags {
		result = append(result, helpers.ToTitle(v.Label))
	}
	return result, nil
}

func (r *PostRepo) GetAll(ctx context.Context) (result []dto.PostRes, err error) {
	var (
		opName = "PostRepository-FindAll"
		query  = r.DB.WithContext(ctx)
		posts  = []models.Post{}
	)
	err = query.Find(&posts).Error
	if err != nil {
		r.Logger.Errorf("%s failed get data posts: %v \n", opName, err)
		return result, err
	}

	for _, v := range posts {
		temp := dto.PostRes{
			ID:      v.ID,
			Title:   v.Title,
			Content: v.Content,
		}
		resTags, err := r.findTags(ctx, v.ID)
		if err != nil {
			r.Logger.Errorf("%s failed get data tags: %v \n", opName, err)
			return result, err
		}

		temp.Tags = resTags
		result = append(result, temp)
	}

	return result, nil
}

func (r *PostRepo) GetDetail(ctx context.Context, req dto.PostGetReq) (*dto.PostRes, error) {
	var (
		opName = "PostRepository-GetDetail"
		query  = r.DB.WithContext(ctx)
		post   = models.Post{}
	)

	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	}

	err := query.First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		r.Logger.Errorf("%s failed get data post: %v \n", opName, err)
		return nil, err
	}

	result := &dto.PostRes{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
	}

	resTags, err := r.findTags(ctx, post.ID)
	if err != nil {
		r.Logger.Errorf("%s failed get data tags: %v \n", opName, err)
		return nil, err
	}
	result.Tags = resTags
	return result, nil
}

func (r *PostRepo) GetDetailTag(ctx context.Context, req dto.TagGetReq) (*models.Tag, error) {
	var (
		opName = "PostRepository-GetDetailTag"
		query  = r.DB.WithContext(ctx)
		result = models.Tag{}
		column = "*"
	)

	err := req.Validate()
	if err != nil {
		r.Logger.Errorf("%s validate params: %v \n", opName, err)
		return nil, err
	}

	if req.ColumnCustom != "" {
		column = req.ColumnCustom
	}

	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	}
	if req.Label != "" {
		query = query.Where("label = ?", req.Label)
	}

	err = query.Select(column).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		r.Logger.Errorf("%s failed get data: %v \n", opName, err)
		return nil, err
	}

	return &result, nil
}
