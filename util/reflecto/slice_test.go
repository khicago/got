package reflecto_test

import (
	"io"
	"reflect"
	"testing"

	"github.com/khicago/got/util/reflecto"
	"github.com/stretchr/testify/assert"
)

func TestIterator_WriteTo(t *testing.T) {
	i := 0
	itr := func() (any, error) {
		i++
		if i > 10 {
			return nil, nil
		}
		return i, nil
	}

	ret := make([]int, 0, 10)
	err := reflecto.Iterator(itr).WriteTo(ret)
	assert.NotNil(t, err, "set to slice should be failed")
	err = reflecto.Iterator(itr).WriteTo(&ret)
	assert.Nil(t, err, "set to pointer of slice should be ok")
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, ret, "set to pointer of slice should be ok")
}

func TestIterator_WriteToPtr(t *testing.T) {
	i := 0

	type c struct {
		Val int
	}

	itr := func() (any, error) {
		i++
		if i > 10 {
			return nil, nil
		}
		return &c{i}, nil
	}

	ret := make([]*c, 0, 10)
	err := reflecto.Iterator(itr).WriteTo(ret)
	assert.NotNil(t, err, "set to slice should be failed")
	err = reflecto.Iterator(itr).WriteTo(&ret)
	assert.Nil(t, err, "set to pointer of slice should be ok")
	for i := 0; i < 10; i++ {
		val := c{i + 1}
		assert.Equal(t, val, *ret[i], "value of item %d should be %v", i, val)
	}
}

func TestIterator_WriteToWithMap(t *testing.T) {
	i := 0

	type c struct {
		Val int
	}

	itr := func() (any, error) {
		i++
		if i > 10 {
			return nil, nil
		}
		return i, nil
	}
	ret := make([]*c, 0, 10)
	err := reflecto.Iterator(itr).WriteTo(&ret, reflecto.ItrMapper(func(in any) (any, error) { return &c{in.(int)}, nil }))
	assert.Nil(t, err, "set to pointer of slice should be ok")
	for i := 0; i < 10; i++ {
		val := c{i + 1}
		assert.Equal(t, val, *ret[i], "value of item %d should be %v", i, val)
	}
}

func TestIterator_WriteToWithReduce(t *testing.T) {
	i := 0

	itr := func() (any, error) {
		i++
		if i > 10 {
			return nil, nil
		}
		return i, nil
	}
	ret := make([]int, 0, 10)
	err := reflecto.Iterator(itr).WriteTo(&ret, reflecto.ItrReducer(func(a any, b any) (any, error) { return a.(int) + b.(int), nil }))
	assert.Nil(t, err, "set to pointer of slice should be ok")
	assert.Equal(t, 55, ret[0], "value should be sum of values")
}

func TestIterator_WriteToWithMapAndReduce(t *testing.T) {
	i := 0

	type c struct {
		Val int
	}

	itr := func() (any, error) {
		i++
		if i > 10 {
			return nil, nil
		}
		return &c{i}, nil
	}
	ret := make([]int, 0, 10)
	err := reflecto.Iterator(itr).WriteTo(&ret,
		reflecto.ItrMapper(func(in any) (any, error) { return in.(*c).Val, nil }),
		reflecto.ItrReducer(func(a any, b any) (any, error) { return a.(int) + b.(int), nil }),
	)
	assert.Nil(t, err, "set to pointer of slice should be ok")
	assert.Equal(t, 55, ret[0], "value should be sum of values")
}

func TestIterator_WriteToWithExitValidator(t *testing.T) {
	i := 0

	itr := func() (any, error) {
		i++
		if i > 10 {
			return nil, io.EOF
		}
		return i, nil
	}
	ret := make([]int, 0, 10)
	err := reflecto.Iterator(itr).WriteTo(&ret, reflecto.ItrExitValidator(func(iv any, err error) (bool, error) {
		if err == io.EOF {
			return true, nil
		}
		return false, err
	}))
	assert.Nil(t, err, "should exit correctly when io.EOF got")
}

func TestGetSliceElementType(t *testing.T) {
	a := make([]int, 0, 10)
	ty, err := reflecto.GetSliceElementType(a)
	assert.Nil(t, err)
	assert.Equal(t, ty, reflect.TypeOf(0))
	ty, err = reflecto.GetSliceElementType(&a)
	assert.Nil(t, err)
	assert.Equal(t, ty, reflect.TypeOf(0))
}

func TestNewSlicePtrReflector(t *testing.T) {
	// receive plain value
	a := []int{1, 2, 3, 4, 5, 6}
	sp, err := reflecto.NewSlicePtrReflector(&a)
	assert.Nil(t, err)

	assert.Equal(t, sp.ItemType(), reflect.TypeOf(1))

	xx := 0
	err = sp.Read(1, &xx)
	assert.Nil(t, err)

	assert.Equal(t, 2, xx)

	// receive elem
	type St struct{ v int }
	b := []St{{1}, {2}}
	sp, err = reflecto.NewSlicePtrReflector(&b)
	assert.Nil(t, err)
	assert.Equal(t, sp.ItemType(), reflect.TypeOf(St{}))
	yy := St{}
	err = sp.Read(1, &yy)
	assert.Nil(t, err)
	assert.Equal(t, St{2}, yy)
	yyy := &St{}
	err = sp.Read(1, yyy)
	assert.Nil(t, err)
	assert.Equal(t, St{2}, *yyy)

	// receive pointer
	c := []*St{{1}, {2}}
	sp, err = reflecto.NewSlicePtrReflector(&c)
	assert.Nil(t, err)
	assert.Equal(t, sp.ItemType(), reflect.TypeOf(&St{}))
	zz := St{}
	err = sp.Read(1, &zz)
	assert.Nil(t, err)
	assert.Equal(t, St{2}, zz)
	zzz := &St{}
	err = sp.Read(1, zzz)
	assert.Nil(t, err)
	assert.Equal(t, c[1], zzz)
}
