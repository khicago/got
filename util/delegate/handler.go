package delegate

type (
	Handler                  Func[error]
	Handler1[TIn any]        Func1[TIn, error]
	Handler2[TIn1, TIn2 any] Func2[TIn1, TIn2, error]
)
