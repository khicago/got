package main

import (
	"context"
	"github.com/khicago/got/table2d/preset"
	"github.com/khicago/got/table2d/preset/pcol"
	"github.com/khicago/got/table2d/tablety"
	"os"
)

func main() {

	var data = tablety.WarpLineReader([][]string{
		{"@", "ID", "INT", "Float", "[", "ID", "]", "{", "ID", "}"},
		{" ", "link(@)", "test($>1,$<50)", "test($%2)", "link(item)", "", "", "select", "", ""},
		{"PID", "LvUp", "Power", "Magic", "InitItems", "", "", "InnerLvUpItem", "LvUp", ""},
		{"10001", "10002", "12", "1.2", "", "1010001", "", "", "1010003", ""},
	})

	p, _ := preset.ReadLines(context.TODO(), data)

	file, err := os.OpenFile("./generated.go", os.O_WRONLY, os.ModePerm)
	if err != nil {
		file, err = os.OpenFile("./table2d/preset/cmd/generated.go", os.O_WRONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	defer func(file *os.File) {
		if err != nil {
			err = file.Close()
		}
	}(file)

	if err = file.Truncate(0); err != nil {
		panic(err)
	}
	if err = pcol.GenerateCode(p.Headline, "main", "the test class", file); err != nil {
		panic(err)
	}

}
