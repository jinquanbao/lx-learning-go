package excel

type Writer interface {
	Write() error
}

type FileWriter interface {
	WriteCompleteCallbacks(writeCompleteCallbacks ...WriteCompleteCallback) FileWriter
	WriteSheets(writeSheets ...*writeSheet) Writer
	//WriteSheet(sheetName string, dest interface{}) *writeSheet
	StreamWriteSheet(sheetName string, dest interface{}) *writeSheet
	Close() error
}
