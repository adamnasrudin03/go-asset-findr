package models

type Post struct {
	ID      uint64 `json:"id" gorm:"primaryKey"`
	Title   string `json:"title" gorm:"not null"`
	Content string `json:"content" gorm:"not null;type:text"`
}

func (Post) TableName() string {
	return "post"
}
