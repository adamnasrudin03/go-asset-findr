package dto

import (
	"github.com/adamnasrudin03/go-template/pkg/helpers"
)

type PostGetReq struct {
	ID           uint64 `json:"id"`
	ColumnCustom string `json:"column_custom"`
}

func (m *PostGetReq) Validate() error {
	m.ColumnCustom = helpers.ToLower(m.ColumnCustom)
	if m.ID == 0 {
		return helpers.ErrIsRequired("id", "id")
	}

	return nil
}
