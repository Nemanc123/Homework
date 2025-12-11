package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		Vers   App             `validate:"nested"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "111111111111111111111111111111111111",
				Name:   "",
				Age:    60,
				Email:  "a@mail.ru",
				Role:   "admin",
				Phones: []string{"11122233355"},
				meta:   nil,
				Vers:   App{Version: "1.0.0"},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Age",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			for l := range tests[i].expectedErr.(ValidationErrors) {
				if err.(ValidationErrors)[l].Field != tt.expectedErr.(ValidationErrors)[l].Field {
					t.Errorf("test doesnt PASS")
				}
			}
			// Place your code here.
			_ = tt
		})
	}
}
