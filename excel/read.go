package excel

type Reader interface {
	Read() error
	Close() error
}

type FileReader interface {
	ReadCompleteCallbacks(readCompleteCallbacks ...ReadCompleteCallback) FileReader
	ReadSheets(readSheets ...*readSheet) Reader
	ReadSheetNo(sheetNo int, dest interface{}) *readSheet
	ReadSheetName(sheetName string, dest interface{}) *readSheet
	Close() error
}
