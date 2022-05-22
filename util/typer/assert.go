package typer

func AssertType[TAssert, TVal any](v TVal) bool {
	_, ok := any(v).(TAssert)
	return ok
}
