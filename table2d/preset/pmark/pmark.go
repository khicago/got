package pmark

type (
	Mark string
)

const (
	NIL            Mark = ""
	BucketsAngleL  Mark = "<"
	BucketsAngleR  Mark = ">"
	BucketsRoundL  Mark = "("
	BucketsRoundR  Mark = ")"
	BucketsSquareL Mark = "["
	BucketsSquareR Mark = "]"
	BucketsCurlyL  Mark = "{"
	BucketsCurlyR  Mark = "}"
)

var (
	pairLR = map[Mark]Mark{
		BucketsAngleL:  BucketsAngleR,
		BucketsRoundL:  BucketsRoundR,
		BucketsSquareL: BucketsSquareR,
		BucketsCurlyL:  BucketsCurlyR,
	}
	pairRL = map[Mark]Mark{
		BucketsAngleR:  BucketsAngleL,
		BucketsRoundR:  BucketsRoundL,
		BucketsSquareR: BucketsSquareL,
		BucketsCurlyR:  BucketsCurlyL,
	}
	names = map[Mark]string{
		BucketsAngleL:  "left angle bucket",
		BucketsAngleR:  "right angle bucket",
		BucketsRoundL:  "left round bucket",
		BucketsRoundR:  "right round bucket",
		BucketsSquareL: "left square bucket",
		BucketsSquareR: "right square bucket",
		BucketsCurlyL:  "left curly bucket",
		BucketsCurlyR:  "right curly bucket",
	}
)

func (p Mark) Name() string {
	return names[p]
}

func (p Mark) Registered() (registered bool) {
	_, registered = names[p]
	return
}

func (p Mark) IsLeft() (isLeft bool) {
	_, isLeft = pairLR[p]
	return
}

func (p Mark) Pairing() Mark {
	r, ok := pairLR[p]
	if ok {
		return r
	}
	l, ok := pairRL[p]
	if ok {
		return l
	}
	return NIL
}

func IsPair[T ~string](begin, end T) bool {
	r, ok := pairLR[Mark(begin)]
	if !ok {
		return false
	}
	return r == Mark(end)
}
