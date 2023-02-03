package injector

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Drafteame/container/dependency"
	"github.com/Drafteame/container/types"
)

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
