package injector

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Drafteame/container/dependency"
)

func TestContainer_Get(t *testing.T) {
	t.Run("get instance", func(t *testing.T) {
		c := New()

		type test struct {
			name string
		}

		err := c.Provide(depName, dependency.New(func() test {
			return test{name: "test"}
		}))
		require.NoError(t, err)

		obj, errGet := c.Get(depName)
		require.NoError(t, errGet)

		tobj, ok := obj.(test)
		require.True(t, ok)

		assert.Equal(t, "test", tobj.name)
	})

	t.Run("get singleton instance", func(t *testing.T) {
		c := New()

		type test struct {
			number int
		}

		err := c.Provide(depName, dependency.NewSingleton(func() test {
			return test{number: rand.Int()}
		}))
		require.NoError(t, err)

		obj1, errGet := c.Get(depName)
		require.NoError(t, errGet)

		obj2, errGet2 := c.Get(depName)
		require.NoError(t, errGet2)

		fmt.Printf("obj1: %p\n", obj1)
		fmt.Printf("obj2: %p\n", obj2)

		assert.Equal(t, obj1.(test).number, obj2.(test).number)
	})
}
