package container

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Drafteame/container/dependency"
)

func TestMustRegister(t *testing.T) {
	t.Run("should register a raw factory instance without panic", func(t *testing.T) {
		c := New()

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustRegister panicked unexpectedly: %v", r)
			}
		}()

		MustRegister(factoryName, newUser, WithArgs(name, age), WithContainer(c))

		// Verify registration was successful by getting the instance
		instance, err := Get[*user](factoryName, WithContainer(c))
		assert.NoError(t, err)
		assert.NotNil(t, instance)
		assert.Equal(t, name, instance.name)
		assert.Equal(t, age, instance.age)
	})

	t.Run("should panic when registration fails", func(t *testing.T) {
		c := New()

		defer func() {
			r := recover()
			assert.NotNil(t, r)

			expErr := fmt.Errorf("factory parameter should be a function or a dependency.Dependency instance")
			assert.Equal(t, expErr.Error(), fmt.Sprintf("%v", r))
		}()

		// This should panic because dependency.Injectable{} is not a valid factory
		MustRegister("name", dependency.Injectable{}, WithContainer(c))
	})
}

func TestMustSingleton(t *testing.T) {
	t.Run("should register a singleton without panic", func(t *testing.T) {
		c := New()

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustSingleton panicked unexpectedly: %v", r)
			}
		}()

		MustSingleton(factoryName, newUser, WithArgs(name, age), WithContainer(c))

		// Verify registration was successful by getting the instance
		instance1, err := Get[*user](factoryName, WithContainer(c))
		assert.NoError(t, err)
		assert.NotNil(t, instance1)

		// Get a second instance to verify it's a singleton (same instance)
		instance2, err := Get[*user](factoryName, WithContainer(c))
		assert.NoError(t, err)
		assert.NotNil(t, instance2)

		// Both instances should be the same object (pointer equality)
		assert.Same(t, instance1, instance2)
	})

	t.Run("should panic when registration fails", func(t *testing.T) {
		c := New()

		defer func() {
			r := recover()
			assert.NotNil(t, r)

			expErr := fmt.Errorf("factory parameter should be a function or a dependency.Dependency instance")
			assert.Equal(t, expErr.Error(), fmt.Sprintf("%v", r))
		}()

		// This should panic because dependency.Injectable{} is not a valid factory
		MustSingleton("name", dependency.Injectable{}, WithContainer(c))
	})
}

func TestMustOverride(t *testing.T) {
	t.Run("should override a dependency without panic", func(t *testing.T) {
		c := New()

		// First, register a dependency
		err := Register(depName, func() int { return 10 }, WithContainer(c))
		require.NoError(t, err)

		// Verify initial value
		v, err := Get[int](depName, WithContainer(c))
		assert.NoError(t, err)
		assert.Equal(t, 10, v)

		// Now override it without a panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustOverride panicked unexpectedly: %v", r)
			}
		}()

		MustOverride(depName, func() int { return 20 }, WithContainer(c))

		// Verify the value was overridden
		v, err = Get[int](depName, WithContainer(c))
		assert.NoError(t, err)
		assert.Equal(t, 20, v)
	})

	t.Run("should panic when override fails with invalid factory", func(t *testing.T) {
		c := New()

		defer func() {
			r := recover()
			assert.NotNil(t, r)

			// The error message will be about the factory not returning a value
			expErr := fmt.Errorf("inject: dependency factory should return at least one return type: dependency.Dependency{Factory: func(), Args: []}")
			assert.Equal(t, expErr.Error(), fmt.Sprintf("%v", r))
		}()

		// This should panic because the factory doesn't return a value
		MustOverride("test_override", func() {}, WithContainer(c))
	})
}
