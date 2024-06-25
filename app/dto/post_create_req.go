package dto

import (
	"strings"

	"github.com/adamnasrudin03/go-template/pkg/helpers"
)

type PostCreateReq struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (m *PostCreateReq) Validate() error {
	m.Title = strings.TrimSpace(m.Title)
	if m.Title == "" {
		return helpers.ErrIsRequired("judul", "title")
	}

	m.Content = strings.TrimSpace(m.Content)
	if m.Content == "" {
		return helpers.ErrIsRequired("konten", "content")
	}

	// filter duplicate value
	tagsNotDuplicate := map[string]bool{}
	tags := []string{}
	for _, v := range m.Tags {
		v = helpers.ToLower(v)
		if v == "" || tagsNotDuplicate[v] {
			continue
		}

		tagsNotDuplicate[v] = true
		tags = append(tags, v)
	}

	m.Tags = tags

	return nil
}
