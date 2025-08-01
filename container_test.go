package container

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Drafteame/container/dependency"
	"github.com/Drafteame/container/injector"
)

const name = "John"
const age = 21
const factoryName = "test"

type user struct {
	name string
	age  int
	db   *sql.DB
}

func newUser(name string, age int) *user {
	return &user{
		age:  age,
		name: name,
	}
}

func newUserError(_ string, _ int) (*user, error) {
	return nil, errors.New("some error")
}

func newUserWithDB(db *sql.DB) *user {
	return &user{db: db}
}

func newDB() *sql.DB {
	return &sql.DB{}
}

func TestNew(t *testing.T) {
	ic := New()

	assert.IsType(t, &injector.Container{}, ic)
	assert.Implements(t, new(Container), ic)
}

func TestSingleton(t *testing.T) {
	t.Run("should register a raw factory singleton instance", func(t *testing.T) {
		t.Cleanup(Flush)

		err := Singleton(factoryName, newUser, WithArgs(name, age))

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from singleton dependency", func(t *testing.T) {
		t.Cleanup(Flush)

		dep := dependency.NewSingleton(newUser, name, age)

		err := Singleton(factoryName, dep)

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from dependency", func(t *testing.T) {
		t.Cleanup(Flush)

		dep := dependency.New(newUser, name, age)

		err := Singleton(factoryName, dep)

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from raw function and nested dependencies", func(t *testing.T) {
		t.Cleanup(Flush)

		if err := Singleton("db", newDB); err != nil {
			t.Fatal(err)
		}

		if err := Singleton(factoryName, newUserWithDB, WithArgs(Inject(depName))); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("error when no dependency.Depdndendency instance or raw function is registered", func(t *testing.T) {
		t.Cleanup(Flush)

		err := Singleton("name", dependency.Injectable{})

		expErr := fmt.Errorf("factory parameter should be a function or a dependency.Dependency instance")

		if assert.Error(t, err) {
			assert.Equal(t, expErr, err)
		}
	})
}

func TestRegister(t *testing.T) {
	t.Run("should register a raw factory instance", func(t *testing.T) {
		t.Cleanup(Flush)

		err := Register(factoryName, newUser, WithArgs(name, age))

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from singleton dependency", func(t *testing.T) {
		t.Cleanup(Flush)

		dep := dependency.NewSingleton(newUser, name, age)

		err := Register(factoryName, dep)

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from dependency", func(t *testing.T) {
		t.Cleanup(Flush)

		dep := dependency.New(newUser, name, age)

		err := Register(factoryName, dep)

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from raw function and nested dependencies", func(t *testing.T) {
		t.Cleanup(Flush)

		if err := Register("db", newDB); err != nil {
			t.Fatal(err)
		}

		if err := Register(factoryName, newUserWithDB, WithArgs(Inject(depName))); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("error when no dependency.Dependency instance or raw function is registered", func(t *testing.T) {
		t.Cleanup(Flush)

		err := Register("name", dependency.Injectable{})

		expErr := fmt.Errorf("factory parameter should be a function or a dependency.Dependency instance")

		if assert.Error(t, err) {
			assert.Equal(t, expErr, err)
		}
	})
}

func TestFunctionalRegistration(t *testing.T) {
	t.Run("non singleton instance", func(t *testing.T) {
		t.Cleanup(Flush)
		defer func() {
			if r := recover(); r != nil {
				t.Fatal(r)
			}
		}()

		const databaseSymbol = "database"
		err := Register(databaseSymbol, newDB)

		if err != nil {
			t.Fatal(err)
		}

		const userSymbol = "user"
		err = Register(userSymbol, func() *user {
			db := MustGet[*sql.DB](databaseSymbol)
			return newUserWithDB(db)
		})
		if err != nil {
			t.Fatal(err)
		}

		userInstance := MustGet[*user](userSymbol)

		assert.NotEmpty(t, userInstance)
		assert.NotNil(t, userInstance.db)
	})
}

func TestRemove(t *testing.T) {
	c := New()
	if err := Register(depName, newDB, WithContainer(c)); err != nil {
		t.Fatal(err)
	}

	_, err := Get[any](depName, WithContainer(c))

	assert.NoError(t, err)

	Remove(depName)

	_, err = Get[any](depName)

	assert.Error(t, err)
}

func TestOverride(t *testing.T) {
	t.Cleanup(Flush)

	if err := Register(depName, func() int { return 10 }); err != nil {
		t.Fatal(err)
	}

	v, err := Get[any](depName)

	assert.NoError(t, err)
	assert.Equal(t, 10, v)

	if err = Override(depName, func() int { return 20 }); err != nil {
		t.Fatal(err)
	}

	v, err = Get[any](depName)

	assert.NoError(t, err)
	assert.Equal(t, 20, v)
}
