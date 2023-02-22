package injector

import (
	"github.com/Drafteame/container/dependency"
	"github.com/Drafteame/container/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainer_Get(t *testing.T) {
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
}