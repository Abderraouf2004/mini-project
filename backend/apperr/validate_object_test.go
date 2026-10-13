package apperr

import "testing"

type testDTO struct {
	Title string `json:"title" validate:"required,notblank,max=200"`
}

func TestValidateObject(t *testing.T) {
	tests := []struct {
		name    string
		data    testDTO
		wantErr bool
	}{
		{
			name:    "valid title",
			data:    testDTO{Title: "Learn Go"},
			wantErr: false,
		},
		{
			name:    "missing title",
			data:    testDTO{Title: ""},
			wantErr: true,
		},
		{
			name:    "whitespace title",
			data:    testDTO{Title: "   "},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateObject(tt.data)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateObject() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}
