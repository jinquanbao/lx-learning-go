package excelutil

import "github.com/xuri/excelize/v2"

type WriteContext interface {
	File() *excelize.File
}

type WriteSheetContext interface {
	WriteContext() WriteContext
	File() *excelize.File
	StreamWriter() *excelize.StreamWriter
	SheetNo() int
	SheetName() string
	TitleRow() int
	TitleBeginColumn() int
}

type WriteRowContext interface {
	File() *excelize.File
	Sheet() WriteSheetContext
	RowIndex() int
}

type WriteCellContext interface {
	File() *excelize.File
	Sheet() WriteSheetContext
	RowIndex() int
	ColumnIndex() int
	FieldName() string
	TitleName() string
	FieldIsTime() bool
}

type writeContextOption struct {
	file *excelize.File
}

func (r *writeContextOption) File() *excelize.File {
	return r.file
}

type writeSheetContextOption struct {
	writeContext     *writeContextOption
	streamWriter     *excelize.StreamWriter
	sheetNo          int
	sheetName        string
	titleRow         int
	titleBeginColumn int
}

func (r *writeSheetContextOption) WriteContext() WriteContext {
	return r.writeContext
}

func (r *writeSheetContextOption) File() *excelize.File {
	return r.writeContext.File()
}

func (r *writeSheetContextOption) StreamWriter() *excelize.StreamWriter {
	return r.streamWriter
}

func (r *writeSheetContextOption) SheetNo() int {
	return r.sheetNo
}

func (r *writeSheetContextOption) SheetName() string {
	return r.sheetName
}

func (r *writeSheetContextOption) TitleRow() int {
	return r.titleRow
}

func (r *writeSheetContextOption) TitleBeginColumn() int {
	return r.titleBeginColumn
}

type writeRowContextOption struct {
	writeSheetContext *writeSheetContextOption
	rowIndex          int
}

func (r *writeRowContextOption) File() *excelize.File {
	return r.writeSheetContext.File()
}

func (r *writeRowContextOption) Sheet() WriteSheetContext {
	return r.writeSheetContext
}

func (r *writeRowContextOption) RowIndex() int {
	return r.rowIndex
}

type writeCellContextOption struct {
	writeSheetContext *writeSheetContextOption
	rowIndex          int
	columnIndex       int
	fieldName         string
	titleName         string
	fieldIsTime       bool
}

func (r *writeCellContextOption) File() *excelize.File {
	return r.writeSheetContext.File()
}

func (r *writeCellContextOption) Sheet() WriteSheetContext {
	return r.writeSheetContext
}

func (r *writeCellContextOption) RowIndex() int {
	return r.rowIndex
}

func (r *writeCellContextOption) ColumnIndex() int {
	return r.columnIndex
}

func (r *writeCellContextOption) FieldName() string {
	return r.fieldName
}

func (r *writeCellContextOption) TitleName() string {
	return r.titleName
}

func (r *writeCellContextOption) FieldIsTime() bool {
	return r.fieldIsTime
}
