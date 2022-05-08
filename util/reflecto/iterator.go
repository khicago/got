package reflecto

import (
	"fmt"
	"reflect"

	"github.com/khicago/got/util/procast"
)

type (
	Iterator func() (any, error)

	ItrMapper        func(iv any) (any, error)
	ItrReducer       func(iv any, in any) (any, error)
	ItrExitValidator func(iv any, err error) (bool, error)
)

func (itr Iterator) Next() (any, error) {
	return itr()
}

func defaultExitValidator(iv any, err error) (bool, error) {
	return iv == nil, err
}

func (itr Iterator) WriteTo(slicePointer any, handler ...any) (err error) {
	defer procast.Recover(func(e error) { err = e })

	vSlice, err := ToWriteableSliceValue(slicePointer)
	if err != nil {
		return fmt.Errorf("input error, %w", err)
	}

	vNewSlice := reflect.MakeSlice(vSlice.Type(), 0, vSlice.Cap())

	var mapper ItrMapper
	var reducer ItrReducer
	// default exit validator
	var exitValidator ItrExitValidator = defaultExitValidator
	for _, h := range handler {
		switch t := h.(type) {
		case ItrMapper:
			mapper = t
		case ItrReducer:
			reducer = t
		case ItrExitValidator:
			exitValidator = t
		}
	}

	for {
		v, err := itr.Next()

		exit, err := exitValidator(v, err)
		if err != nil {
			return fmt.Errorf("itr failed, %w", err)
		}
		if exit {
			break
		}

		if mapper != nil {
			if v, err = mapper(v); err != nil {
				return fmt.Errorf("mapping failed, %w", err)
			}
		}

		if reducer != nil {
			if vNewSlice.Len() == 0 {
				vNewSlice = reflect.Append(vNewSlice, reflect.ValueOf(v))
			} else {
				v0 := vNewSlice.Index(0)
				if v, err = reducer(v0.Interface(), v); err != nil {
					return fmt.Errorf("reducing failed, %w", err)
				}
				v0.Set(reflect.ValueOf(v))
			}
		} else {
			vNewSlice = reflect.Append(vNewSlice, reflect.ValueOf(v))
		}

	}
	vSlice.Set(vNewSlice)
	return nil
}
