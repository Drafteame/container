package container

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Drafteame/container/dependency"
	"github.com/Drafteame/container/injector"
	"github.com/Drafteame/container/types"
)

const name = "John"

const age = 21

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

func TestInvoke(t *testing.T) {
	defer Flush()

	depName := types.Symbol("userTest")
	dep := dependency.New(newUser, name, age)

	if err := Register(depName, dep); err != nil {
		t.Fatal(err)
	}

	type args struct {
		types.In
		User *user `inject:"name=userTest"`
	}

	called := false

	invoker := func(in args) {
		if assert.NotNil(t, in.User) {
			assert.Equal(t, in.User.age, age)
			assert.Equal(t, in.User.name, name)
		}

		called = true
	}

	err := Invoke(invoker)

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestSingleton(t *testing.T) {
	t.Run("should register a raw factory singleton instance", func(t *testing.T) {
		defer Flush()

		factoryName := "test"

		err := Singleton(factoryName, newUser, name, age)

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from singleton dependency", func(t *testing.T) {
		defer Flush()

		factoryName := "test"

		dep := dependency.NewSingleton(newUser, name, age)

		err := Singleton(factoryName, dep)

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from dependency", func(t *testing.T) {
		defer Flush()

		factoryName := "test"

		dep := dependency.New(newUser, name, age)

		err := Singleton(factoryName, dep)

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from raw function and nested dependencies", func(t *testing.T) {
		defer Flush()

		depName := "db"

		if err := Singleton(depName, newDB); err != nil {
			t.Fatal(err)
		}

		factoryName := "test"

		if err := Singleton(factoryName, newUserWithDB, Inject(depName)); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("error when no dependency.Depdndendency instance or raw function is registered", func(t *testing.T) {
		defer Flush()

		err := Singleton("name", dependency.Injectable{})

		expErr := fmt.Errorf("factory parameter should be a function or a dependency.Dependency instance")

		if assert.Error(t, err) {
			assert.Equal(t, expErr, err)
		}
	})
}

func TestRegister(t *testing.T) {
	t.Run("should register a raw factory instance", func(t *testing.T) {
		defer Flush()

		factoryName := "test"

		err := Register(factoryName, newUser, name, age)

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from singleton dependency", func(t *testing.T) {
		defer Flush()

		factoryName := "test"

		dep := dependency.NewSingleton(newUser, name, age)

		err := Register(factoryName, dep)

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from dependency", func(t *testing.T) {
		defer Flush()

		factoryName := "test"

		dep := dependency.New(newUser, name, age)

		err := Register(factoryName, dep)

		assert.NoError(t, err)
	})

	t.Run("should register a singleton from raw function and nested dependencies", func(t *testing.T) {
		defer Flush()

		depName := "db"

		if err := Register(depName, newDB); err != nil {
			t.Fatal(err)
		}

		factoryName := "test"

		if err := Register(factoryName, newUserWithDB, Inject(depName)); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("error when no dependency.Dependency instance or raw function is registered", func(t *testing.T) {
		defer Flush()

		err := Register("name", dependency.Injectable{})

		expErr := fmt.Errorf("factory parameter should be a function or a dependency.Dependency instance")

		if assert.Error(t, err) {
			assert.Equal(t, expErr, err)
		}
	})
}

func TestFunctionalRegistration(t *testing.T) {
	t.Run("non singleton instance", func(t *testing.T) {
		defer Flush()
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
	const depName = "test"

	if err := Register(depName, newDB); err != nil {
		t.Fatal(err)
	}

	_, err := Get[any](depName)

	assert.NoError(t, err)

	Remove(depName)

	_, err = Get[any](depName)

	assert.Error(t, err)
}

func TestTestMode(t *testing.T) {
	t.Run("get singleton instance on test mode", func(t *testing.T) {
		TestMode()

		type test struct {
			number int
		}

		sym := types.Symbol("test")

		err := Singleton(sym, func() test {
			return test{number: rand.Int()}
		})

		if err != nil {
			t.Fatal(err)
		}

		obj1, errGet := Get[test](sym)
		if errGet != nil {
			t.Fatal(errGet)
		}

		obj2, errGet2 := Get[test](sym)
		if errGet2 != nil {
			t.Fatal(errGet2)
		}

		assert.NotEqual(t, obj1.number, obj2.number)
	})
}
