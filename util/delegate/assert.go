package delegate

type (
	Condition            Func[bool]
	Predicate[TVal any]  Convert[TVal, bool]
	PredicateE[TVal any] func(TVal) (bool, error)
)
