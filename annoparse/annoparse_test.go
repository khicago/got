package annoparse

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Parser is the Annotation target
type Parser struct {
	Field string `anno:"name=f"`
}

type Parent struct {
	B float32 `parser:"f=b"`
}

type Target struct {
	A int `parser:"f=a"`
	Parent
	C  string `parser:"Field=c"`
	DD string `parser:"field=dd"`
}

func (*Parser) TagName() string {
	return "parser"
}

func TestKVStr_ReflectTo(t *testing.T) {
	targetObj := Target{}
	table, err := ExtractAnno(&targetObj, &Parser{})
	assert.Nil(t, err)

	fmt.Printf("table %#v \n", table)
	assert.Equal(t, "a", table[".A"].Field)
	assert.Equal(t, "b", table[".Parent.B"].Field)
	assert.Equal(t, "c", table[".C"].Field)
	assert.Equal(t, "dd", table[".DD"].Field)
}
