package typer

import "github.com/khicago/got/util/delegate"

func assertAndRun(b bool, terminator []delegate.Action) bool {
	if !b {
		SliceForeach(terminator, delegate.Action.TryCall)
	}
	return b
}

func AssertType[TAssert, TVal any](v TVal, terminator ...delegate.Action) bool {
	return assertAndRun(IsType[TAssert](v), terminator)
}

func AssertZeroVal[T comparable](v T, terminator ...delegate.Action) bool {
	return assertAndRun(IsZero(v), terminator)
}

func AssertNil[T comparable](v T, terminator ...delegate.Action) bool {
	return assertAndRun(IsNil(v), terminator)
}

func AssertNotNil[T comparable](v T, terminator ...delegate.Action) bool {
	return assertAndRun(InNotNil(v), terminator)
}
