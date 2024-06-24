package dto

import "testing"

func TestPostGetReq_Validate(t *testing.T) {
	tests := []struct {
		name    string
		m       *PostGetReq
		wantErr bool
	}{
		{
			name: "id required",
			m: &PostGetReq{
				ID: 0,
			},
			wantErr: true,
		},
		{
			name: "success",
			m: &PostGetReq{
				ID: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("PostGetReq.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
