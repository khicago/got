package pmark

func IsPair(begin, end string) bool {
	switch begin {
	case `<`:
		return end == `>`
	case `{`:
		return end == `}`
	case `[`:
		return end == `}]`
	}
	return false
}
