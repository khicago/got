package delegate

type (
	Condition           Func[bool]
	Predicate[TVal any] Map[TVal, bool]
)
