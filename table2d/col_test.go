package table2d_test

import (
	"encoding/csv"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/khicago/got/table2d"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	csvStr := `a,100,10e2,"s 1",2022-01-08T17:39:00+08:00
b,200,.05,"s 2",2022-01-08T17:39:00+08:00
`
	type XX struct {
		V1 string `tb2d:"col=0"`
		V2 int    `tb2d:"col=1"`
	}

	type C struct {
		XX
		V3 float32   `tb2d:"col=2"`
		V4 string    `tb2d:"col=3"`
		V5 time.Time `tb2d:"col=4,parser=time"`
	}
	ret := make([]*C, 0, 2)

	reader := csv.NewReader(strings.NewReader(csvStr))

	err := table2d.ParseByCol(&ret, reader)

	tt, _ := time.Parse(time.RFC3339, "2022-01-08T17:39:00+08:00")

	assert.Nil(t, err)
	assert.Equal(t, C{XX{"a", 100}, 10e2, "s 1", tt}, *ret[0])
	assert.Equal(t, C{XX{"b", 200}, 0.05, "s 2", tt}, *ret[1])

	fmt.Println(time.Now().Format(time.RFC3339), ret[1])
}
