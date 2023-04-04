package injector

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Drafteame/container/dependency"
	"github.com/Drafteame/container/types"
)

func TestContainer_Get(t *testing.T) {
	t.Run("get instance", func(t *testing.T) {
		c := New()

		type test struct {
			name string
		}

		sym := types.Symbol("test")

		err := c.Provide(sym, dependency.New(func() test {
			return test{name: "test"}
		}))
		if err != nil {
			t.Fatal(err)
		}

		obj, errGet := c.Get(sym)
		if errGet != nil {
			t.Fatal(errGet)
		}

		tobj, ok := obj.(test)
		if !ok {
			t.Fatal("cant cast obj to test")
		}

		assert.Equal(t, "test", tobj.name)
	})

	t.Run("get singleton instance in regular mode", func(t *testing.T) {
		c := New()

		type test struct {
			number int
		}

		sym := types.Symbol("test")

		err := c.Provide(sym, dependency.NewSingleton(func() test {
			return test{number: rand.Int()}
		}))
		if err != nil {
			t.Fatal(err)
		}

		obj1, errGet := c.Get(sym)
		if errGet != nil {
			t.Fatal(errGet)
		}

		obj2, errGet2 := c.Get(sym)
		if errGet2 != nil {
			t.Fatal(errGet2)
		}

		fmt.Printf("obj1: %p\n", obj1)
		fmt.Printf("obj2: %p\n", obj2)

		assert.Equal(t, obj1.(test).number, obj2.(test).number)
	})

	t.Run("get singleton instance on test mode", func(t *testing.T) {
		c := New()
		c.TestMode()

		type test struct {
			number int
		}

		sym := types.Symbol("test")

		err := c.Provide(sym, dependency.NewSingleton(func() test {
			return test{number: rand.Int()}
		}))
		if err != nil {
			t.Fatal(err)
		}

		obj1, errGet := c.Get(sym)
		if errGet != nil {
			t.Fatal(errGet)
		}

		obj2, errGet2 := c.Get(sym)
		if errGet2 != nil {
			t.Fatal(errGet2)
		}

		assert.NotEqual(t, obj1.(test).number, obj2.(test).number)
	})
}
