package injector

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Drafteame/container/dependency"
)

func TestContainer_Flush(t *testing.T) {
	ic := New()

	if err := ic.Provide(depName, dependency.New(func() int { return 10 })); err != nil {
		t.Error(err)
		return
	}

	ic.Flush()

	assert.Empty(t, ic.deps)
	assert.Empty(t, ic.solvedDeps)
}

func TestContainer_Remove(t *testing.T) {
	depName2 := "test2"

	c := New()

	if err := c.Provide(depName, dependency.NewSingleton(func() int { return 10 })); err != nil {
		t.Fatal(err)
	}

	if err := c.Provide(depName2, dependency.NewSingleton(func() int { return 10 })); err != nil {
		t.Fatal(err)
	}

	_, _ = c.Get(depName)
	_, _ = c.Get(depName2)

	assert.Len(t, c.deps, 2)
	assert.Len(t, c.solvedDeps, 2)

	c.Remove(depName)

	assert.Len(t, c.deps, 1)
	assert.Len(t, c.solvedDeps, 1)

	c.Remove(depName2)

	assert.Empty(t, c.deps)
	assert.Empty(t, c.solvedDeps)
}

func TestContainer_Override(t *testing.T) {
	c := New()

	if err := c.Provide(depName, dependency.NewSingleton(func() int { return 10 })); err != nil {
		t.Fatal(err)
	}

	v, err := c.Get(depName)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 10, v)

	if err = c.Override(depName, dependency.NewSingleton(func() int { return 20 })); err != nil {
		t.Fatal(err)
	}

	v, err = c.Get(depName)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 20, v)
}
