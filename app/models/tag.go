package models

type Tag struct {
	ID    uint64 `json:"id" gorm:"primaryKey"`
	Label string `json:"label" gorm:"not null;unique"`
}

func (Tag) TableName() string {
	return "tag"
}
