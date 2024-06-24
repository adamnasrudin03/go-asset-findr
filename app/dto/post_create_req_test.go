package dto

import "testing"

func TestPostCreateReq_Validate(t *testing.T) {
	tests := []struct {
		name    string
		m       *PostCreateReq
		wantErr bool
	}{
		{
			name: "title required",
			m: &PostCreateReq{
				Title:   "",
				Content: "test content",
				Tags:    []string{},
			},
			wantErr: true,
		}, {
			name: "content required",
			m: &PostCreateReq{
				Title:   "test title",
				Content: "",
				Tags:    []string{},
			},
			wantErr: true,
		},
		{
			name: "success",
			m: &PostCreateReq{
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
				t.Errorf("PostCreateReq.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
