package excel

import "github.com/xuri/excelize/v2"

type ReadContext interface {
	File() *excelize.File
}

type ReadSheetContext interface {
	File() *excelize.File
	SheetNo() int
	SheetName() string
	TitleRow() int
	TitleBeginColumn() int
}

type ReadCellContext interface {
	File() *excelize.File
	Sheet() ReadSheetContext
	RowIndex() int
	ColumnIndex() int
	TitleName() string
	CellValue() string
}

type readContextOption struct {
	file *excelize.File
}

type readSheetContextOption struct {
	file             *excelize.File
	sheetNo          int
	sheetName        string
	titleRow         int
	titleBeginColumn int
}

type readCellContextOption struct {
	file             *excelize.File
	readSheetContext *readSheetContextOption
	rowIndex         int
	columnIndex      int
	titleName        string
	cellValue        string
}

func (r *readContextOption) File() *excelize.File {
	return r.file
}

func (r *readSheetContextOption) File() *excelize.File {
	return r.file
}

func (r *readSheetContextOption) SheetNo() int {
	return r.sheetNo
}

func (r *readSheetContextOption) SheetName() string {
	return r.sheetName
}

func (r *readSheetContextOption) TitleRow() int {
	return r.titleRow
}

func (r *readSheetContextOption) TitleBeginColumn() int {
	return r.titleBeginColumn
}

func (r *readCellContextOption) File() *excelize.File {
	return r.file
}

func (r *readCellContextOption) Sheet() ReadSheetContext {
	return r.readSheetContext
}

func (r *readCellContextOption) RowIndex() int {
	return r.rowIndex
}

func (r *readCellContextOption) ColumnIndex() int {
	return r.columnIndex
}

func (r *readCellContextOption) TitleName() string {
	return r.titleName
}

func (r *readCellContextOption) CellValue() string {
	return r.cellValue
}
