package dto

import "testing"

func TestPostUpdateReq_Validate(t *testing.T) {
	tests := []struct {
		name    string
		m       *PostUpdateReq
		wantErr bool
	}{
		{
			name: "id required",
			m: &PostUpdateReq{
				ID:      0,
				Title:   "title test",
				Content: "test content",
				Tags:    []string{},
			},
			wantErr: true,
		},
		{
			name: "title required",
			m: &PostUpdateReq{
				ID:      1,
				Title:   "",
				Content: "test content",
				Tags:    []string{},
			},
			wantErr: true,
		},
		{
			name: "content required",
			m: &PostUpdateReq{
				ID:      1,
				Title:   "test title",
				Content: "",
				Tags:    []string{},
			},
			wantErr: true,
		},
		{
			name: "success",
			m: &PostUpdateReq{
				ID:      1,
				Title:   "title 1",
				Content: "content 2",
				Tags:    []string{"tags1", "tags2", "tags1"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("PostUpdateReq.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
