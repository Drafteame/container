package container

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Drafteame/container/dependency"
)

const depName = "userTest"

func TestGet(t *testing.T) {
	t.Run("get instance of specific type by name", func(t *testing.T) {
		c := New()

		dep := dependency.New(newUser, name, age)

		err := Register(depName, dep, WithContainer(c))
		require.NoError(t, err)

		ui, err := Get[*user](depName, WithContainer(c))

		assert.NoError(t, err)
		assert.NotEmpty(t, ui)
		assert.Equal(t, ui.age, age)
		assert.Equal(t, ui.name, name)
	})

	t.Run("get instance from container error", func(t *testing.T) {
		c := New()

		dep := dependency.New(newUserError, name, age)

		err := Register(depName, dep, WithContainer(c))
		require.NoError(t, err)

		ui, err := Get[*user](depName, WithContainer(c))
		expErr := errors.New("inject: error building dependency instance: inject: error constructing `func(string, int) (*container.user, error)`: some error")

		assert.Error(t, err)
		assert.Empty(t, ui)
		assert.Equal(t, expErr.Error(), err.Error())
	})

	t.Run("cast type error", func(t *testing.T) {
		c := New()

		dep := dependency.New(newUser, name, age)

		err := Register(depName, dep, WithContainer(c))
		require.NoError(t, err)

		ui, err := Get[string](depName, WithContainer(c))
		expErr := errors.New("inject: error casting instance of `userTest` dependency to `string`")

		assert.Error(t, err)
		assert.Empty(t, ui)
		assert.Equal(t, expErr, err)
	})
}

func TestMustGet(t *testing.T) {
	t.Run("get instance of specific type by name", func(t *testing.T) {
		c := New()

		defer func() {
			if r := recover(); r != nil {
				t.Error(r)
			}
		}()

		dep := dependency.New(newUser, name, age)

		err := Register(depName, dep, WithContainer(c))
		require.NoError(t, err)

		ui := MustGet[*user](depName, WithContainer(c))

		assert.NotEmpty(t, ui)
		assert.Equal(t, ui.age, age)
		assert.Equal(t, ui.name, name)
	})

	t.Run("get instance from container error", func(t *testing.T) {
		c := New()
		defer func() {
			r := recover()
			expErr := errors.New("inject: error building dependency instance: inject: error constructing `func(string, int) (*container.user, error)`: some error")

			assert.Equal(t, expErr, fmt.Errorf("%v", r))
		}()

		dep := dependency.New(newUserError, name, age)

		err := Register(depName, dep, WithContainer(c))
		require.NoError(t, err)

		val := MustGet[*user](depName, WithContainer(c))
		assert.NotEmpty(t, val)
		assert.Equal(t, val.age, age)
		assert.Equal(t, val.name, name)
	})

	t.Run("cast type error", func(t *testing.T) {
		c := New()

		defer func() {
			r := recover()
			expErr := errors.New("inject: error casting instance of `userTest` dependency to `string`")

			assert.Equal(t, expErr, fmt.Errorf("%v", r))
		}()

		dep := dependency.New(newUser, name, age)

		err := Register(depName, dep, WithContainer(c))
		require.NoError(t, err)

		_ = MustGet[string](depName, WithContainer(c))
	})
}
