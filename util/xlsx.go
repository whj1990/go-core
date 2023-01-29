package util

import (
	"bytes"
	"github.com/whj1990/go-core/handler"
	"github.com/tealeg/xlsx"
)

func ReadXlsx(xlsxBytes []byte) ([][]string, error) {
	var xlData [][]string
	xlFile, err := xlsx.OpenBinary(xlsxBytes)
	if err != nil {
		return xlData, err
	}
	for index, sheet := range xlFile.Sheets {
		//第一个sheet
		if index == 0 {
			temp := make([][]string, len(sheet.Rows))
			for k, row := range sheet.Rows {
				var data []string
				for _, cell := range row.Cells {
					data = append(data, cell.Value)
				}
				temp[k] = data
			}
			xlData = append(xlData, temp...)
		}
	}
	return xlData, nil
}

func WriteXlsx(xlData [][]string) ([]byte, error) {
	resultFile := xlsx.NewFile()
	sheet, err := resultFile.AddSheet("sheet1")
	if err != nil {
		return nil, handler.HandleError(err)
	}
	if xlData != nil && len(xlData) > 0 {
		for _, v := range xlData {
			row := sheet.AddRow()
			row.SetHeightCM(1)
			for _, c := range v {
				cell := row.AddCell()
				cell.Value = c
			}
		}
	}
	buf := new(bytes.Buffer)
	if err = resultFile.Write(buf); err != nil {
		return nil, handler.HandleError(err)
	}
	return buf.Bytes(), nil
}
