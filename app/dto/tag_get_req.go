package dto

import (
	"github.com/adamnasrudin03/go-template/pkg/helpers"
)

type TagGetReq struct {
	ID           uint64 `json:"id"`
	Label        string `json:"label"`
	ColumnCustom string `json:"column_custom"`
}

func (m *TagGetReq) Validate() error {
	m.ColumnCustom = helpers.ToLower(m.ColumnCustom)
	m.Label = helpers.ToLower(m.Label)

	isNoParams := m.ID == 0 && m.Label == ""
	if isNoParams {
		return helpers.NewError(helpers.ErrValidation, helpers.NewResponseMultiLang(
			helpers.MultiLanguages{
				ID: "Harap masukkan minimal satu parameter yang diperlukan",
				EN: "Please provide at least one required parameter",
			},
		))
	}

	return nil
}
