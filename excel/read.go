package excelutil

type Reader interface {
	Read() error
	Close() error
}

type FileReader interface {
	ReadCompleteCallbacks(readCompleteCallbacks ...ReadCompleteCallback) FileReader
	ReadSheets(readSheets ...*ReadSheet) Reader
	ReadSheetNo(sheetNo int, dest interface{}) *ReadSheet
	ReadSheetName(sheetName string, dest interface{}) *ReadSheet
	Close() error
}
