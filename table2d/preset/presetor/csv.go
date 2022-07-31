package presetor

import (
	"context"
	"encoding/csv"
	"github.com/khicago/got/table2d/preset"
	"io"
	"os"
	"strings"
)

func CSV(r io.Reader) (*preset.Preset, error) {
	return preset.ReadLines(context.TODO(), csv.NewReader(r))
}

func CSVFile(path string) (*preset.Preset, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		if err != nil {
			err = file.Close()
		}
	}(file)

	return CSV(file) // csv.NewReader contains a bufio.NewReader
}

func CSVStr(str string) (*preset.Preset, error) {
	return CSV(strings.NewReader(str))
}
