package injector

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Drafteame/container/dependency"
	"github.com/Drafteame/container/types"
)

func TestContainer_EnableTestMode(t *testing.T) {
	ic := New()
	assert.False(t, ic.testMode)

	ic.TestMode()

	assert.True(t, ic.testMode)
}

func TestContainer_Flush(t *testing.T) {
	depName := "test"

	ic := New()

	if err := ic.Provide(types.Symbol(depName), dependency.New(func() int { return 10 })); err != nil {
		t.Error(err)
		return
	}

	ic.Flush()

	assert.Empty(t, ic.deps)
	assert.Empty(t, ic.solvedDeps)
}

func TestContainer_Remove(t *testing.T) {
	depName := types.Symbol("test")
	depName2 := types.Symbol("test2")

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
