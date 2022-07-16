package procast

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/khicago/got/util/typer"
)

type (
	FrameNode struct {
		File string  `json:"file"`
		Func string  `json:"func"`
		Line int     `json:"line"`
		PC   uintptr `json:"pc"`
	}
)

func (n FrameNode) FuncName() string {
	i := strings.LastIndex(n.Func, "/")
	if i != -1 {
		return n.Func[i+1:]
	}
	return n.Func
}

func (n FrameNode) LocString() string {
	switch {
	case n.Func != "":
		return fmt.Sprintf("[func:%v:%v]", n.FuncName(), n.Line)
	case n.File != "":
		return fmt.Sprintf("[file:%v:%v]", n.File, n.Line)
	case n.PC != 0:
		return fmt.Sprintf("[pc:%x]", n.PC)
	}
	return "[empty-frame]"
}

func GetCallersPCLst(callerSkip int, depthMax uint) []uintptr {
	if depthMax < 1 {
		depthMax = 1
	}
	pcs := make([]uintptr, depthMax)
	n := runtime.Callers(2+callerSkip, pcs)
	return pcs[:n]
}

func GetFramesOfPCs(pcs []uintptr, prefixSkip string) []FrameNode {
	frames := runtime.CallersFrames(pcs)
	frame := runtime.Frame{}

	ret := make([]FrameNode, 0, len(pcs))
	for more := true; more; {
		frame, more = frames.Next()
		if prefixSkip != "" && strings.HasPrefix(frame.Function, prefixSkip) {
			continue
		}
		ret = append(ret, FrameNode{
			PC:   frame.PC,
			File: frame.File,
			Func: frame.Function,
			Line: frame.Line,
		})
	}
	return ret
}

func GetFrameNodes(callerSkip int, prefixSkip string) []FrameNode {
	pc := GetCallersPCLst(1+callerSkip, 16)
	// fmt.Printf("pc %v\n", pc)
	return GetFramesOfPCs(pc, prefixSkip)
}

func GetFrameNode(callerSkip int, prefixSkip string) FrameNode {
	pc := GetCallersPCLst(1+callerSkip, 16)
	// fmt.Printf("pc %v\n", pc)
	nodes := GetFramesOfPCs(pc, prefixSkip)
	fmt.Printf("nodes %v\n", nodes)
	return typer.IfThen(len(nodes) > 0, nodes[0], typer.ZeroVal[FrameNode]())
}
