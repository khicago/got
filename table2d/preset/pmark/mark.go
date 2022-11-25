package pmark

type (
	// Mark is a string type that represents a mark
	// it is used to pair marks
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

// Name returns the name of the mark
func (p Mark) Name() string {
	return names[p]
}

// Registered returns true if the mark is registered
func (p Mark) Registered() (registered bool) {
	// NIL mark is not considered as registered
	_, registered = names[p]
	return
}

// IsLeft returns true if the mark is a left mark
func (p Mark) IsLeft() (isLeft bool) {
	_, isLeft = pairLR[p]
	return
}

// Pairing returns the pairing mark
// - if the mark is not registered, it returns NIL
// - if the mark is a left mark, it returns the right mark
// - if the mark is a right mark, it returns the left mark
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

// PairedWith returns true if the mark is paired with the given mark
func (p Mark) PairedWith(other Mark) bool {
	if !p.Registered() || !other.Registered() {
		return false
	}
	return p.Pairing() == other
}

// IsPair returns true if the two marks are a pair
func IsPair[T ~string](begin, end T) bool {
	r, ok := pairLR[Mark(begin)]
	if !ok {
		return false
	}
	return r == Mark(end)
}
