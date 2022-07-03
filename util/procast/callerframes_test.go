package procast

import (
	"fmt"
	"testing"
)

func TestCallFrames(t *testing.T) {
	fns := GetFrameNodes(0, "runtime.")
	fmt.Printf("%#v", fns)
}
