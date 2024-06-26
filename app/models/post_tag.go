package models

type PostTag struct {
	ID     uint64 `json:"id" gorm:"primaryKey"`
	PostID uint64 `json:"post_id" gorm:"not null"`
	TagID  uint64 `json:"tag_id" gorm:"not null"`
	Tag    *Tag   `json:"tag" gorm:"foreignKey:TagID"`
	Post   *Post  `json:"post" gorm:"foreignKey:PostID"`
}

func (PostTag) TableName() string {
	return "post_tag"
}
