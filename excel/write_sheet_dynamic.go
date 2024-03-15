package excelutil

import (
	"time"

	"github.com/xuri/excelize/v2"
)

type dynamicWriteSheet struct {
	streamWriter                *excelize.StreamWriter
	writeWorkbook               *WriteWorkbook
	option                      *writeSheetOption
	writeSheetBeforeCallbacks   []WriteSheetBeforeCallback
	writeSheetCompleteCallbacks []WriteSheetCompleteCallback
	writeCellBeforeCallbacks    []WriteCellBeforeCallback
	writeRowBeforeCallbacks     []WriteRowBeforeCallback
	writeSheetContext           *writeSheetContextOption
	writeRowContext             *writeRowContextOption
	writeCellContext            *writeCellContextOption
	titles                      [][]interface{}
	values                      [][]interface{}
	indexTitleMap               map[int]string
}

func (w *dynamicWriteSheet) TitleRow(titleRow int) *dynamicWriteSheet {
	w.option.TitleRow = titleRow
	return w
}

func (w *dynamicWriteSheet) ContentBeginRow(contentBeginRow int) *dynamicWriteSheet {
	w.option.ContentBeginRow = contentBeginRow
	return w
}

func (w *dynamicWriteSheet) TitleBeginColumn(titleBeginColumn int) *dynamicWriteSheet {
	w.option.TitleBeginColumn = titleBeginColumn
	return w
}

func (w *dynamicWriteSheet) RegisterWriteSheetBeforeCallbacks(opts ...WriteSheetBeforeCallback) *dynamicWriteSheet {
	w.writeSheetBeforeCallbacks = opts
	return w
}

func (w *dynamicWriteSheet) RegisterWriteSheetCompleteCallbacks(opts ...WriteSheetCompleteCallback) *dynamicWriteSheet {
	w.writeSheetCompleteCallbacks = opts
	return w
}

func (w *dynamicWriteSheet) RegisterWriteCellBeforeCallbacks(opts ...WriteCellBeforeCallback) *dynamicWriteSheet {
	w.writeCellBeforeCallbacks = opts
	return w
}

func (w *dynamicWriteSheet) RegisterWriteRowBeforeCallbacks(opts ...WriteRowBeforeCallback) *dynamicWriteSheet {
	w.writeRowBeforeCallbacks = opts
	return w
}

func (w *dynamicWriteSheet) Write() error {
	return w.writeWorkbook.WriteSheets(w).Write()
}

func (w *dynamicWriteSheet) getOption() *writeSheetOption {
	return w.option
}

func (w *dynamicWriteSheet) preWrite() (err error) {
	w.streamWriter, err = w.writeWorkbook.File.NewStreamWriter(w.option.SheetName)
	if err != nil {
		return err
	}

	w.writeSheetContext = &writeSheetContextOption{
		writeContext:     &writeContextOption{file: w.writeWorkbook.File},
		streamWriter:     w.streamWriter,
		sheetName:        w.option.SheetName,
		sheetNo:          w.option.SheetNo,
		titleRow:         w.option.TitleRow,
		titleBeginColumn: w.option.TitleBeginColumn,
	}
	w.writeRowContext = &writeRowContextOption{
		writeSheetContext: w.writeSheetContext,
		rowIndex:          w.option.TitleRow,
	}
	w.writeCellContext = &writeCellContextOption{
		writeSheetContext: w.writeSheetContext,
		rowIndex:          w.option.TitleRow,
	}

	if w.option.ContentBeginRow < w.option.TitleRow+len(w.titles) {
		w.option.ContentBeginRow = w.option.TitleRow + len(w.titles)
	}

	indexTitleMap := make(map[int]string, len(w.titles[len(w.titles)-1]))
	for i, v := range w.titles[len(w.titles)-1] {
		title, ok := v.(string)
		if !ok {
			continue
		}
		indexTitleMap[i] = title
	}
	w.indexTitleMap = indexTitleMap
	return nil
}

func (w *dynamicWriteSheet) doWrite() error {
	for _, opt := range w.writeSheetBeforeCallbacks {
		if err := opt(w.writeSheetContext); err != nil {
			return err
		}
	}
	if err := w.writeContent(w.titles, w.option.TitleRow, true, w.writeWorkbook.DefaultTitleStyleId); err != nil {
		return err
	}
	if err := w.writeContent(w.values, w.option.ContentBeginRow, false, w.writeWorkbook.DefaultValueStyleId); err != nil {
		return err
	}
	for _, opt := range w.writeSheetCompleteCallbacks {
		if err := opt(w.writeSheetContext); err != nil {
			return err
		}
	}
	if err := w.streamWriter.Flush(); err != nil {
		return err
	}
	return nil
}

func (w *dynamicWriteSheet) writeContent(values [][]interface{}, startRow int, isTitle bool, styleId int) (err error) {
	var (
		writeRowContext  = w.writeRowContext
		writeCellContext = w.writeCellContext
	)
	writeRowContext.rowIndex = startRow
	writeCellContext.rowIndex = startRow

	for rowID := 0; rowID < len(values); rowID++ {
		row := values[rowID]
		rowIndex := startRow + rowID
		for i, v := range row {
			cell := excelize.Cell{}
			if val, ok := v.(excelize.Cell); ok {
				cell = val
			} else if _, ok := v.(time.Time); ok {
				cell = excelize.Cell{StyleID: w.writeWorkbook.DefaultTimeStyleId, Value: v}
			} else {
				cell = excelize.Cell{StyleID: styleId, Value: v}
			}

			columnIndex := i + w.option.TitleBeginColumn
			writeCellContext.columnIndex = columnIndex
			writeCellContext.rowIndex = rowIndex
			if isTitle {
				writeCellContext.titleName = v.(string)
			} else {
				writeCellContext.titleName = w.indexTitleMap[i]
			}
			for _, opt := range w.writeCellBeforeCallbacks {
				if err = opt(writeCellContext, isTitle, &cell); err != nil {
					return err
				}
			}
			row[i] = cell
		}

		writeRowContext.rowIndex = rowIndex
		for _, opt := range w.writeRowBeforeCallbacks {
			if err := opt(writeRowContext, isTitle, &row); err != nil {
				return err
			}
		}
		cell, _ := excelize.CoordinatesToCellName(w.option.TitleBeginColumn+1, rowIndex+1)
		if err = w.streamWriter.SetRow(cell, row); err != nil {
			return err
		}
	}
	return nil
}
