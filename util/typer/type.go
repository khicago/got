package typer

type (
	Predicate[TVal any] func(val TVal) bool
)
