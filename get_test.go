package container

import (
	"errors"
	"fmt"
	"github.com/Drafteame/container/dependency"
	"github.com/Drafteame/container/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	t.Run("get instance of specific type by name", func(t *testing.T) {
		defer Flush()

		depName := types.Symbol("userTest1")
		dep := dependency.New(newUser, name, age)

		if err := Register(depName, dep); err != nil {
			t.Error(err)
			return
		}

		ui, err := Get[*user](depName)

		assert.NoError(t, err)
		assert.NotEmpty(t, ui)
		assert.Equal(t, ui.age, age)
		assert.Equal(t, ui.name, name)
	})

	t.Run("get instance from container error", func(t *testing.T) {
		defer Flush()

		depName := types.Symbol("userTest2")
		dep := dependency.New(newUserError, name, age)

		if err := Register(depName, dep); err != nil {
			t.Error(err)
			return
		}

		ui, err := Get[*user](string(depName))
		expErr := errors.New("inject: error building dependency instance: inject: error constructing `func(string, int) (*container.user, error)`: some error")

		assert.Error(t, err)
		assert.Empty(t, ui)
		assert.Equal(t, expErr, err)
	})

	t.Run("cast type error", func(t *testing.T) {
		defer Flush()

		depName := types.Symbol("userTest3")
		dep := dependency.New(newUser, name, age)

		if err := Register(depName, dep); err != nil {
			t.Error(err)
			return
		}

		ui, err := Get[string](depName)
		expErr := errors.New("inject: error casting instance of `userTest3` dependency to `string`")

		assert.Error(t, err)
		assert.Empty(t, ui)
		assert.Equal(t, expErr, err)
	})
}

func TestMustGet(t *testing.T) {
	t.Run("get instance of specific type by name", func(t *testing.T) {
		defer Flush()
		defer func() {
			if r := recover(); r != nil {
				t.Error(r)
			}
		}()

		depName := types.Symbol("userTest1")
		dep := dependency.New(newUser, name, age)

		if err := Register(depName, dep); err != nil {
			t.Error(err)
			return
		}

		ui := MustGet[*user](depName)

		assert.NotEmpty(t, ui)
		assert.Equal(t, ui.age, age)
		assert.Equal(t, ui.name, name)
	})

	t.Run("get instance from container error", func(t *testing.T) {
		defer Flush()
		defer func() {
			r := recover()
			expErr := errors.New("inject: error building dependency instance: inject: error constructing `func(string, int) (*container.user, error)`: some error")

			assert.Equal(t, expErr, fmt.Errorf("%v", r))
		}()

		depName := types.Symbol("userTest2")
		dep := dependency.New(newUserError, name, age)

		if err := Register(depName, dep); err != nil {
			t.Error(err)
			return
		}

		_ = MustGet[*user](string(depName))
	})

	t.Run("cast type error", func(t *testing.T) {
		defer Flush()
		defer func() {
			r := recover()
			expErr := errors.New("inject: error casting instance of `userTest3` dependency to `string`")

			assert.Equal(t, expErr, fmt.Errorf("%v", r))
		}()

		depName := types.Symbol("userTest3")
		dep := dependency.New(newUser, name, age)

		if err := Register(depName, dep); err != nil {
			t.Error(err)
			return
		}

		_ = MustGet[string](depName)
	})
}
