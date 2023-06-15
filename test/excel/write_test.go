package excel

import (
	"github.com/jinquanbao/lx-learning-go/excel"
	"github.com/xuri/excelize/v2"
	"io"
	"io/fs"
	"io/ioutil"
	"strconv"
	"testing"
	"time"
)

func TestWriteSingleSheet(t *testing.T) {
	o1 := Order{
		OrderId:   "o1",
		OrderTime: time.Now(),
		Message:   "o1消息",
	}
	o2 := Order{
		OrderId:   "o2",
		OrderTime: time.Now(),
		Message:   "o2消息",
	}
	list := []Order{o1, o2}
	err := excel.StreamWrite(&list, "WriteSingleSheet.xlsx")
	if err != nil {
		t.Fatal(err)
	}
}

func TestWriteMultipleSheet(t *testing.T) {
	o1 := Order{
		OrderId:   "o1",
		OrderTime: time.Now(),
		Message:   "o1消息",
	}
	o2 := Order{
		OrderId:   "o2",
		OrderTime: time.Now(),
		Message:   "o2消息",
	}
	o3 := Order{
		OrderId:   "o3",
		OrderTime: time.Now(),
		Message:   "o3消息",
	}
	o4 := Order{
		OrderId:   "o4",
		OrderTime: time.Now(),
		Message:   "o4消息",
	}
	o5 := Order{
		OrderId:   "o5",
		OrderTime: time.Now(),
		Message:   "o5消息",
	}
	list1 := []Order{o1, o2, o3}
	list2 := []Order{o4, o5}

	writer := excel.NewWriter(excel.WithWriteSaveFilePath("WriteMultipleSheet.xlsx"))

	err := writer.WriteSheets(
		writer.StreamWriteSheet("sheet1", &list1),
		writer.StreamWriteSheet("sheet2", &list2),
	).Write()
	if err != nil {
		t.Fatal(err)
	}
}

func TestWriteToFileBuff(t *testing.T) {
	o1 := Order{
		OrderId:   "o1",
		OrderTime: time.Now(),
		Message:   "o1消息",
	}
	o2 := Order{
		OrderId:   "o2",
		OrderTime: time.Now(),
		Message:   "o2消息",
	}
	list1 := []Order{o1, o2}

	writer := excel.NewWriter(excel.WithWriteDisableAutoClose(true))
	writer.StreamWriteSheet("sheet1", &list1).Write()
	//需要手动关闭
	defer func() {
		writer.Close()
	}()

	piper, pipeW := io.Pipe()
	defer piper.Close()
	go func() {
		defer pipeW.Close()
		writer.File.Write(pipeW)
	}()
	data, err := ioutil.ReadAll(piper)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("WriteToFileBuff.xlsx", data, fs.ModeExclusive)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWriteCustomStyle(t *testing.T) {
	o1 := Order{
		OrderId:   "o1",
		OrderTime: time.Now(),
		Message:   "o1消息",
	}
	o2 := Order{
		OrderId:   "o2",
		OrderTime: time.Now(),
		Message:   "o2消息",
	}
	list1 := []Order{o1, o2}

	writer := excel.NewWriter(excel.WithWriteSaveFilePath("WriteCustomStyle.xlsx"))

	yellowStyleID, _ := writer.File.NewStyle(&excelize.Style{
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
		// 图案填充-黄色
		Fill: excelize.Fill{
			Pattern: 1,
			Type:    "pattern",
			Color:   []string{"FFC60A"},
		},
		// 居中
		Alignment: &excelize.Alignment{
			// 水平对齐：
			Horizontal: "center",
			// 垂直对齐
			Vertical: "center",
		},
	})

	err := writer.StreamWriteSheet("sheet1", &list1).
		RegisterWriteSheetBeforeCallbacks(func(wCtx excel.WriteSheetContext) error {
			err := wCtx.StreamWriter().SetColWidth(1, 3, 12)
			if err != nil {
				return err
			}
			err = wCtx.StreamWriter().SetColWidth(2, 2, 18)
			if err != nil {
				return err
			}
			return nil
		}).
		RegisterWriteCellBeforeCallbacks(func(wCtx excel.WriteCellContext, isTitle bool, cell *excelize.Cell) error {
			if isTitle && wCtx.TitleName() == "订单ID" {
				cell.StyleID = yellowStyleID
			}
			return nil
		}).
		Write()
	if err != nil {
		t.Fatal(err)
	}
}

func TestWriteDynamicTitle1(t *testing.T) {
	o1 := Order{
		OrderId:   "o1",
		OrderTime: time.Now(),
		Message:   "o1消息",
	}
	o2 := Order{
		OrderId:   "o2",
		OrderTime: time.Now(),
		Message:   "o2消息",
	}
	list1 := []Order{o1, o2}

	writer := excel.NewWriter(excel.WithWriteSaveFilePath("WriteDynamicTitle1.xlsx"))

	err := writer.StreamWriteSheet("sheet1", &list1).
		Titles("消息", "订单ID"). // 根据传入的标题顺序写入
		Write()
	if err != nil {
		t.Fatal(err)
	}
}

func TestWriteDynamicTitle2(t *testing.T) {
	o1 := Order{
		OrderId:   "o1",
		OrderTime: time.Now(),
		Message:   "o1消息",
	}
	o2 := Order{
		OrderId:   "o2",
		OrderTime: time.Now(),
		Message:   "o2消息",
	}
	list1 := []Order{o1, o2}

	writer := excel.NewWriter(excel.WithWriteSaveFilePath("WriteDynamicTitle2.xlsx"))

	err := writer.StreamWriteSheet("sheet1", &list1).
		IncludeTitleNames("消息", "订单ID"). //根据结构体定义的标题顺序写入
		Write()
	if err != nil {
		t.Fatal(err)
	}
}

func TestWriteDynamicTitle3(t *testing.T) {
	o1 := Order{
		OrderId:   "o1",
		OrderTime: time.Now(),
		Message:   "o1消息",
	}
	o2 := Order{
		OrderId:   "o2",
		OrderTime: time.Now(),
		Message:   "o2消息",
	}
	list1 := []Order{o1, o2}

	extraTitles := []string{"扩展的标题1", "扩展的标题2"}
	extraValues := []string{"扩展的内容1", "扩展的内容2"}

	writer := excel.NewWriter(excel.WithWriteSaveFilePath("WriteDynamicTitle3.xlsx"))

	err := writer.StreamWriteSheet("sheet1", &list1).
		TitleBeginColumn(1).
		IncludeTitleNames("消息", "订单ID").
		RegisterWriteRowBeforeCallbacks(func(wCtx excel.WriteRowContext, isTitle bool, row *[]interface{}) error {
			if isTitle {
				for _, v := range extraTitles {
					*row = append(*row, excelize.Cell{StyleID: writer.DefaultTitleStyleId, Value: v})
				}
			} else {
				for _, v := range extraValues {
					*row = append(*row, excelize.Cell{StyleID: writer.DefaultValueStyleId, Value: "row" + strconv.Itoa(wCtx.RowIndex()) + "-" + v})
				}
			}
			return nil
		}).
		Write()
	if err != nil {
		t.Fatal(err)
	}
}
