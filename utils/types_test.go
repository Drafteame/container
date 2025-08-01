package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFirstReturnType(t *testing.T) {
	t.Run("get type from function with single return", func(t *testing.T) {
		fun := func() string { return "" }

		tfun := GetFirstReturnType(fun)

		assert.NotNil(t, tfun)
		assert.Equal(t, reflect.TypeOf(""), tfun)
	})

	t.Run("get type from function with no return", func(t *testing.T) {
		fun := func() {}

		tfun := GetFirstReturnType(fun)

		assert.Nil(t, tfun)
	})

	t.Run("get type from function with multi-return values", func(t *testing.T) {
		fun := func() (string, error) { return "", nil }

		tfun := GetFirstReturnType(fun)

		assert.NotNil(t, tfun)
		assert.Equal(t, reflect.TypeOf(""), tfun)
	})

	t.Run("get type from non function", func(t *testing.T) {
		tfun := GetFirstReturnType("")

		assert.Nil(t, tfun)
	})
}
