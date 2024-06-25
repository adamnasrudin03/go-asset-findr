package dto

import "github.com/adamnasrudin03/go-template/pkg/helpers"

type PostRes struct {
	ID      uint64   `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (m *PostRes) CheckResp() {
	m.Title = helpers.ToTitle(m.Title)

	if m.Tags == nil || len(m.Tags) == 0 {
		m.Tags = []string{}
	}

}
