package dto

import (
	"github.com/adamnasrudin03/go-template/pkg/helpers"
)

type PostGetReq struct {
	ID uint64 `json:"id"`
}

func (m *PostGetReq) Validate() error {
	if m.ID == 0 {
		return helpers.ErrIsRequired("id", "id")
	}

	return nil
}
