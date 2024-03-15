package excelutil

import "github.com/xuri/excelize/v2"

type Writer interface {
	Write() error
}

type FileWriter interface {
	WriteCompleteCallbacks(writeCompleteCallbacks ...WriteCompleteCallback) FileWriter
	WriteSheets(writeSheets ...WriteSheet) Writer
	// WriteSheet(sheetName string, dest interface{}) *writeSheet
	StreamWriteSheet(sheetName string, dest interface{}) *writeSheet
	DynamicWriteSheet(sheetName string, titles, values [][]interface{}) *dynamicWriteSheet
	Close() error
}

type WriteSheet interface {
	Writer
	getOption() *writeSheetOption
	preWrite() error
	doWrite() error
}

type (
	WriteSheetBeforeCallback   func(wCtx WriteSheetContext) error
	WriteSheetCompleteCallback func(wCtx WriteSheetContext) error
	WriteCellBeforeCallback    func(wCtx WriteCellContext, isTitle bool, cell *excelize.Cell) error
	WriteRowBeforeCallback     func(wCtx WriteRowContext, isTitle bool, row *[]interface{}) error
)
