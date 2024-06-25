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
	"gorm.io/gorm/clause"
)

type PostRepository interface {
	GetAll(ctx context.Context) (result []dto.PostRes, err error)
	GetDetail(ctx context.Context, req dto.PostGetReq) (*dto.PostRes, error)
	GetDetailTag(ctx context.Context, req dto.TagGetReq) (*models.Tag, error)
	Create(ctx context.Context, req dto.PostCreateReq) (*dto.PostRes, error)
	DeleteByID(ctx context.Context, postID uint64) error
	UpdateByID(ctx context.Context, req dto.PostUpdateReq) error
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
		column = "*"
	)
	if req.ColumnCustom != "" {
		column = req.ColumnCustom
	}

	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	}

	err := query.Select(column).First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.ErrNotFound()
		}

		r.Logger.Errorf("%s failed get data post: %v \n", opName, err)
		return nil, helpers.ErrDB()
	}

	result := &dto.PostRes{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
	}

	if column != "*" {
		return result, nil
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
		return nil, helpers.ErrDB()
	}

	return &result, nil
}

func (r *PostRepo) trxEnd(trx *gorm.DB, err error) {
	if rc := recover(); rc != nil {
		r.Logger.Errorf(`trxEnd Panic Error %v`, r)
		trx.Rollback()
		return
	}
	if err != nil {
		r.Logger.Errorf(`trxEnd Error %v`, err)
		trx.Rollback()
		return
	}
	if err := trx.Commit(); err != nil {
		r.Logger.Errorf(`trxEnd err commit %v`, err)
		trx.Rollback()
		return
	}
}

func (r *PostRepo) Create(ctx context.Context, req dto.PostCreateReq) (*dto.PostRes, error) {
	var (
		opName = "PostRepository-Create"
		err    error
		trx    *gorm.DB
		post   = models.Post{
			Title:   req.Title,
			Content: req.Content,
		}
	)

	trx = r.DB.Begin().WithContext(ctx)
	defer func() {
		r.trxEnd(trx, err)
	}()

	err = trx.Clauses(clause.Returning{}).Create(&post).Error
	if err != nil {
		r.Logger.Errorf("%s failed create data: %v \n", opName, err)
		return nil, err
	}

	for _, val := range req.Tags {
		_, err = r.createPostTag(ctx, trx, post.ID, models.Tag{Label: val})
		if err != nil {
			r.Logger.Errorf("%s failed create post_tag: %v \n", opName, err)
			return nil, err
		}

	}

	result := &dto.PostRes{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Tags:    req.Tags,
	}
	return result, nil
}

func (r *PostRepo) DeleteByID(ctx context.Context, postID uint64) error {
	var (
		opName = "PostRepository-DeleteByID"
		err    error
		trx    *gorm.DB
	)

	_, err = r.GetDetail(ctx, dto.PostGetReq{
		ID:           postID,
		ColumnCustom: "id",
	})
	if err != nil {
		r.Logger.Errorf("%s failed get data post: %v \n", opName, err)
		return err
	}

	trx = r.DB.Begin().WithContext(ctx)
	defer r.trxEnd(trx, err)

	err = trx.Where("post_id = ?", postID).Delete(&models.PostTag{}).Error
	if err != nil {
		r.Logger.Errorf("%s failed delete data post-tag: %v \n", opName, err)
		return helpers.ErrDB()
	}

	err = trx.Where("id = ?", postID).Delete(&models.Post{}).Error
	if err != nil {
		r.Logger.Errorf("%s failed delete data post: %v \n", opName, err)
		return helpers.ErrDB()
	}

	return nil
}

func (r *PostRepo) UpdateByID(ctx context.Context, req dto.PostUpdateReq) error {
	var (
		opName = "PostRepository-UpdateByID"
		err    error
		trx    *gorm.DB
	)

	_, err = r.GetDetail(ctx, dto.PostGetReq{
		ID:           req.ID,
		ColumnCustom: "id",
	})
	if err != nil {
		r.Logger.Errorf("%s failed get data post: %v \n", opName, err)
		return err
	}

	trx = r.DB.Begin().WithContext(ctx)
	defer r.trxEnd(trx, err)

	err = trx.Model(&models.Post{}).Where("id = ?", req.ID).Updates(&models.Post{
		ID:      req.ID,
		Title:   req.Title,
		Content: req.Content,
	}).Error
	if err != nil {
		r.Logger.Errorf("%s failed update data post: %v \n", opName, err)
		return helpers.ErrUpdatedDB()
	}

	err = trx.Where("post_id = ?", req.ID).Delete(&models.PostTag{}).Error
	if err != nil {
		r.Logger.Errorf("%s failed delete data post-tag: %v \n", opName, err)
		return helpers.ErrUpdatedDB()
	}

	for _, val := range req.Tags {
		_, err = r.createPostTag(ctx, trx, req.ID, models.Tag{Label: val})
		if err != nil {
			r.Logger.Errorf("%s failed create post_tag: %v \n", opName, err)
			return err
		}
	}

	return nil
}

func (r *PostRepo) createPostTag(ctx context.Context, trx *gorm.DB, postID uint64, req models.Tag) (*models.Tag, error) {
	var (
		opName = "PostRepository-createPostTag"
		err    error
	)

	tag, err := r.GetDetailTag(ctx, dto.TagGetReq{
		ColumnCustom: "id",
		Label:        req.Label,
	})
	if err != nil {
		r.Logger.Errorf("%s failed get data tags: %v \n", opName, err)
		return nil, err
	}

	isExist := tag != nil && tag.ID > 0
	if !isExist {
		tag = &models.Tag{
			Label: req.Label,
		}
		err = trx.Clauses(clause.Returning{}).Create(tag).Error
		if err != nil {
			r.Logger.Errorf("%s failed create data tags: %v \n", opName, err)
			return nil, helpers.ErrDB()
		}
	}

	err = trx.Create(&models.PostTag{
		PostID: postID,
		TagID:  tag.ID,
	}).Error
	if err != nil {
		r.Logger.Errorf("%s failed create data post-tag: %v \n", opName, err)
		return nil, helpers.ErrDB()
	}

	return tag, nil
}
