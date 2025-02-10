package hw09structvalidator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type User struct {
	ID     string `validate:"len:36"`
	Name   string
	Age    int      `validate:"min:18|max:50"`
	Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	Role   string   `validate:"in:admin,stuff"`
	Phones []string `validate:"len:11"`
}

func TestValidateUser(t *testing.T) {
	t.Run("valid user", func(t *testing.T) {
		u := User{
			ID:     "012345678901234567890123456789012345",
			Name:   "John Doe",
			Age:    33,
			Email:  "test@mail.ru",
			Role:   "stuff",
			Phones: []string{"12345678901"},
		}
		err := Validate(u)
		require.NoError(t, err)
	})

	t.Run("user with errors", func(t *testing.T) {
		u := User{
			ID:     "0",
			Name:   "Name",
			Age:    0,
			Email:  "invalid",
			Role:   "none",
			Phones: []string{"0"},
		}
		expected := ValidationErrors{
			{Field: "ID", Err: NewErrStringLengthMismatch(36)},
			{Field: "Age", Err: NewErrValueOutOfRange("greater", 18)},
			{Field: "Email", Err: NewErrStringDoesNotMatchPattern("^\\w+@\\w+\\.\\w+$")},
			{Field: "Role", Err: NewErrValueNotInSet("admin,stuff")},
			{Field: "Phones", Err: NewErrStringLengthMismatch(11)},
		}
		err := Validate(u)
		require.Error(t, err)
		require.True(t, reflect.DeepEqual(expected, err))
	})
}
