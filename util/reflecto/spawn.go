package reflecto

import (
	"reflect"
)

type (
	Spawner[T any] func() T
)

func NewAnySpawner(model any) (spawner Spawner[any]) {
	return NewAnySpawnerFromType(reflect.TypeOf(model))
}

func NewTSpawner[T any](model T) (spawner Spawner[T]) {
	anySpawner := NewAnySpawnerFromType(reflect.TypeOf(model))
	return func() T {
		return anySpawner.Spawn().(T)
	}
}

func NewAnySpawnerFromType(ty reflect.Type) (spawner Spawner[any]) {
	if ty.Kind() == reflect.Ptr {
		return func() any { return reflect.New(ty.Elem()).Interface() } // todo: .Addr() ? test this
	} else {
		return func() any { return reflect.New(ty).Interface() }
	}
}

func (sp Spawner[T]) Spawn() T {
	return sp()
}
