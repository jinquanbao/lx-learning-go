package excelutil

import "github.com/xuri/excelize/v2"

var defaultStyle = &style{}

type style struct{}

func (s *style) newDefaultTitleStyleId(file *excelize.File) (int, error) {
	styleID, err := file.NewStyle(&excelize.Style{
		// 标题加粗
		Font: &excelize.Font{
			Bold:  true,
			Color: "FFFFFF",
		},
		// 设置边框
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		// 图案填充-绿色
		Fill: excelize.Fill{
			Pattern: 1,
			Type:    "pattern",
			Color:   []string{"2EA121"},
		},
		// 居中
		Alignment: &excelize.Alignment{
			// 水平对齐：
			Horizontal: "center",
			// 垂直对齐
			Vertical: "center",
		},
	})
	return styleID, err
}

func (s *style) newDefaultValueStyleId(file *excelize.File) (int, error) {
	styleID, err := file.NewStyle(&excelize.Style{
		// 设置边框
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		// 居中
		Alignment: &excelize.Alignment{
			// 水平对齐：
			Horizontal: "center",
			// 垂直对齐
			Vertical: "center",
		},
	})
	return styleID, err
}

func (s *style) newDefaultTimeStyleId(file *excelize.File) (int, error) {
	styleID, err := file.NewStyle(&excelize.Style{
		// 设置边框
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		// 居中
		Alignment: &excelize.Alignment{
			// 水平对齐：
			Horizontal: "center",
			// 垂直对齐
			Vertical: "center",
		},
		NumFmt: 22,
	})
	return styleID, err
}
