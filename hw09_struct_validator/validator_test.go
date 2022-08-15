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
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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
		{in: App{Version: "1234"}, expectedErr: ValidationErrors{{Field: "Version", Err: ErrInvalidLen}}},
		{in: App{Version: "123456"}, expectedErr: ValidationErrors{{Field: "Version", Err: ErrInvalidLen}}},
		{
			in:          Response{Code: 505, Body: `{"result":"Тест"}`},
			expectedErr: ValidationErrors{{Field: "Code", Err: ErrInvalidIn}},
		},
		{in: User{
			ID:     "12345678_",
			Name:   "Юзер",
			Age:    17,
			Email:  "qq@qq123.ru",
			Role:   "user",
			Phones: []string{"12345678"},
			meta:   nil,
		}, expectedErr: ValidationErrors{
			{Field: "ID", Err: ErrInvalidLen},
			{Field: "Age", Err: ErrInvalidMin},
			{Field: "Role", Err: ErrInvalidIn},
			{Field: "Phones", Err: ErrInvalidLen},
		}},
		{in: User{
			ID:     "123456789012345678901234567890123456",
			Name:   "User",
			Age:    19,
			Email:  "qq@qq.ru",
			Role:   "admin",
			Phones: []string{"12345678901"},
			meta:   nil,
		}, expectedErr: ValidationErrors{}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			// Place your code here.
			_ = tt
		})
	}
}
