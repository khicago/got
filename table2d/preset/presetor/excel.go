package presetor

import (
	"context"

	"github.com/khicago/got/table2d/preset"
	"github.com/xuri/excelize/v2"
)

func ExcelFile(filePath, sheetName string) (*preset.Preset, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *excelize.File) {
		if err != nil {
			err = file.Close()
		}
	}(f)

	if sheetName == "" {
		sheetName = "Sheet1"
	}
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	return Strings2d(context.TODO(), rows) // csv.NewReader contains a bufio.NewReader
}
