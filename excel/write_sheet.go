package excelutil

import (
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/xuri/excelize/v2"
)

type writeSheet struct {
	streamWriter                *excelize.StreamWriter
	writeWorkbook               *WriteWorkbook
	option                      *writeSheetOption
	writeSheetBeforeCallbacks   []WriteSheetBeforeCallback
	writeSheetCompleteCallbacks []WriteSheetCompleteCallback
	writeCellBeforeCallbacks    []WriteCellBeforeCallback
	writeRowBeforeCallbacks     []WriteRowBeforeCallback
	writeSheetContext           *writeSheetContextOption
	writeCellContext            *writeCellContextOption
	schema                      *schema
	titles                      []*column
	dest                        interface{}
	destIndirectValue           reflect.Value
}

func (w *writeSheet) Titles(titles ...string) *writeSheet {
	w.option.Titles = titles
	return w
}

func (w *writeSheet) TitleRow(titleRow int) *writeSheet {
	w.option.TitleRow = titleRow
	return w
}

func (w *writeSheet) ContentBeginRow(contentBeginRow int) *writeSheet {
	w.option.ContentBeginRow = contentBeginRow
	return w
}

func (w *writeSheet) TitleBeginColumn(titleBeginColumn int) *writeSheet {
	w.option.TitleBeginColumn = titleBeginColumn
	return w
}

func (w *writeSheet) IncludeTitleNames(includeTitleNames ...string) *writeSheet {
	w.option.IncludeTitleNames = includeTitleNames
	return w
}

func (w *writeSheet) ExcludeTitleNames(excludeTitleNames ...string) *writeSheet {
	w.option.ExcludeTitleNames = excludeTitleNames
	return w
}

func (w *writeSheet) IncludeColumnFieldNames(includeColumnFieldNames ...string) *writeSheet {
	w.option.IncludeColumnFieldNames = includeColumnFieldNames
	return w
}

func (w *writeSheet) ExcludeColumnFieldNames(excludeColumnFieldNames ...string) *writeSheet {
	w.option.ExcludeColumnFieldNames = excludeColumnFieldNames
	return w
}

func (w *writeSheet) RegisterWriteSheetBeforeCallbacks(opts ...WriteSheetBeforeCallback) *writeSheet {
	w.writeSheetBeforeCallbacks = opts
	return w
}

func (w *writeSheet) RegisterWriteSheetCompleteCallbacks(opts ...WriteSheetCompleteCallback) *writeSheet {
	w.writeSheetCompleteCallbacks = opts
	return w
}

func (w *writeSheet) RegisterWriteCellBeforeCallbacks(opts ...WriteCellBeforeCallback) *writeSheet {
	w.writeCellBeforeCallbacks = opts
	return w
}

func (w *writeSheet) RegisterWriteRowBeforeCallbacks(opts ...WriteRowBeforeCallback) *writeSheet {
	w.writeRowBeforeCallbacks = opts
	return w
}

func (w *writeSheet) Write() error {
	return w.writeWorkbook.WriteSheets(w).Write()
}

func (w *writeSheet) getOption() *writeSheetOption {
	return w.option
}

func (w *writeSheet) preWrite() error {
	val := reflect.ValueOf(w.dest)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Slice {
		return ErrInputDest
	}

	elemType := val.Type().Elem()
	if elemType.Kind() != reflect.Struct {
		return ErrInputDestElem
	}

	schema, err := newSchema(val)
	if err != nil {
		return err
	}
	w.destIndirectValue = val

	titles := make([]*column, 0, len(schema.Columns))
	titleMap := make(map[string]*column, len(schema.Columns))

	for _, column := range schema.Columns {

		if containStr(w.option.ExcludeTitleNames, column.Name) {
			continue
		}

		if containStr(w.option.ExcludeColumnFieldNames, column.FieldName) {
			continue
		}

		writeTitle := true
		if len(w.option.IncludeTitleNames) > 0 && !containStr(w.option.IncludeTitleNames, column.Name) {
			writeTitle = false
		}
		if len(w.option.IncludeColumnFieldNames) > 0 && !containStr(w.option.IncludeColumnFieldNames, column.FieldName) {
			writeTitle = false
		}

		if writeTitle {
			if column.Index == -1 {
				column.Index = column.FieldIndex
			}
			titles = append(titles, column)
			titleMap[column.Name] = column
		}

	}

	if len(w.option.Titles) > 0 {
		w.titles = make([]*column, 0, len(w.option.Titles))
		for _, v := range w.option.Titles {
			if column, ok := titleMap[v]; ok {
				w.titles = append(w.titles, column)
			} else {
				fmt.Println("title name [" + v + "]  not exist")
				return ErrTitleNotMatch
			}
		}

	} else {
		sort.Slice(titles, func(i, j int) bool {
			return titles[i].Index < titles[j].Index
		})
		w.titles = titles
	}

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

	return nil
}

func (w *writeSheet) doWrite() (err error) {
	for _, opt := range w.writeSheetBeforeCallbacks {
		if err = opt(w.writeSheetContext); err != nil {
			return err
		}
	}

	if w.option.StreamWriter {
		return w.streamWrite()
	} else {
		panic("not support yet")
	}
}

func (w *writeSheet) streamWrite() error {
	writeRowContext := &writeRowContextOption{
		writeSheetContext: w.writeSheetContext,
		rowIndex:          w.option.TitleRow,
	}

	writeCellContext := &writeCellContextOption{
		writeSheetContext: w.writeSheetContext,
		rowIndex:          w.option.TitleRow,
	}

	// 写入标题
	titleCells := make([]interface{}, w.option.TitleBeginColumn, len(w.titles)+w.option.TitleBeginColumn)
	for i, v := range w.titles {
		cell := excelize.Cell{StyleID: w.writeWorkbook.DefaultTitleStyleId, Value: v.Name}

		columnIndex := i + w.option.TitleBeginColumn
		writeCellContext.columnIndex = columnIndex
		writeCellContext.fieldName = v.FieldName
		writeCellContext.titleName = v.Name

		for _, opt := range w.writeCellBeforeCallbacks {
			if err := opt(writeCellContext, true, &cell); err != nil {
				return err
			}
		}
		titleCells = append(titleCells, cell)
	}

	for _, opt := range w.writeRowBeforeCallbacks {
		if err := opt(writeRowContext, true, &titleCells); err != nil {
			return err
		}
	}
	cell, err := excelize.CoordinatesToCellName(1, w.option.TitleRow+1)
	if err != nil {
		return err
	}
	if err := w.streamWriter.SetRow(cell, titleCells, excelize.RowOpts{Height: 20}); err != nil {
		return err
	}

	contentBeginRow := w.option.TitleRow + 1
	if contentBeginRow < w.option.ContentBeginRow {
		contentBeginRow = w.option.ContentBeginRow
	}

	for rowID := 0; rowID < w.destIndirectValue.Len(); rowID++ {
		row := make([]interface{}, w.option.TitleBeginColumn, len(w.titles)+w.option.TitleBeginColumn)
		elem := w.destIndirectValue.Index(rowID)

		rowIndex := rowID + contentBeginRow
		writeRowContext.rowIndex = rowIndex

		for i, column := range w.titles {
			field := elem.Field(column.FieldIndex)
			var fieldVal interface{}
			if !field.IsZero() {
				fieldVal = reflect.Indirect(field).Interface()
			}
			cell := excelize.Cell{StyleID: w.writeWorkbook.DefaultValueStyleId, Value: fieldVal}
			if timeV, ok := fieldVal.(time.Time); ok {
				cell.StyleID = w.writeWorkbook.DefaultTimeStyleId
				if timeV.Year() < 1976 {
					cell.Value = nil
				}
			}

			columnIndex := i + w.option.TitleBeginColumn
			writeCellContext.rowIndex = rowIndex
			writeCellContext.columnIndex = columnIndex
			writeCellContext.fieldName = column.FieldName
			writeCellContext.titleName = column.Name
			writeCellContext.fieldIsTime = column.FieldIsTime
			for _, opt := range w.writeCellBeforeCallbacks {
				if err = opt(writeCellContext, false, &cell); err != nil {
					return err
				}
			}

			row = append(row, cell)
		}

		for _, opt := range w.writeRowBeforeCallbacks {
			if err = opt(writeRowContext, false, &row); err != nil {
				return err
			}
		}
		cell, err := excelize.CoordinatesToCellName(1, rowIndex+1)
		if err != nil {
			return err
		}
		if err := w.streamWriter.SetRow(cell, row); err != nil {
			return err
		}
	}

	for _, opt := range w.writeSheetCompleteCallbacks {
		if err = opt(w.writeSheetContext); err != nil {
			return err
		}
	}

	if err := w.streamWriter.Flush(); err != nil {
		return err
	}

	return nil
}
