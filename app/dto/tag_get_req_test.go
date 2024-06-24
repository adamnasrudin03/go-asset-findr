package dto

import "testing"

func TestTagGetReq_Validate(t *testing.T) {
	tests := []struct {
		name    string
		m       *TagGetReq
		wantErr bool
	}{
		{
			name: "failed params must be filled",
			m: &TagGetReq{
				ID:    0,
				Label: "",
			},
			wantErr: true,
		},
		{
			name: "success",
			m: &TagGetReq{
				ID: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("TagGetReq.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
