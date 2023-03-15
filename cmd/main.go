package main

import (
	"fmt"
	"github.com/Drafteame/container"
	"reflect"
)

type user struct {
	ID string
}

type factory struct {
	model any
}

func newFactory(model any) *factory {
	return &factory{model}
}

func (f *factory) MustCreateWithOption(opts map[string]any) any {
	model := f.model.(user)

	if val, ok := opts["ID"]; ok {
		model.ID = val.(string)
	}

	return model
}

func RegisterFactory[T any](model any) {
	t := reflect.TypeOf(*new(T))
	key := fmt.Sprintf("testFactory-%s", t.Name())

	if err := container.Register(key, newFactory, model); err != nil {
		panic(err)
	}
}

func MustGetWithOptions[T any](opts map[string]any) T {
	t := reflect.TypeOf(*new(T))
	key := fmt.Sprintf("testFactory-%s", t.Name())

	f, err := container.Get[*factory](key)
	if err != nil {
		panic(err)
	}

	return f.MustCreateWithOption(opts).(T)
}

// nolint
func main() {
	RegisterFactory[user](user{ID: "some"})
	m := MustGetWithOptions[user](map[string]any{"ID": "other"})

	fmt.Println(m)
}
