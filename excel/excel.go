package excel

import (
	"io"

	"github.com/xuri/excelize/v2"
)

func ReadFile(filePath string, dest interface{}, opts ...excelize.Options) error {
	fileReader := OpenFile(filePath, opts...)
	defer fileReader.Close()
	return fileReader.ReadSheetNo(0, dest).Read()
}

func Read(reader io.Reader, dest interface{}, opts ...excelize.Options) error {
	fileReader := OpenReader(reader, opts...)
	defer fileReader.Close()
	return fileReader.ReadSheetNo(0, dest).Read()
}

func OpenFile(filePath string, opts ...excelize.Options) FileReader {
	return newReadWorkbook(&readWorkbookOption{
		filePath: filePath,
		options:  opts,
	})
}

func OpenReader(reader io.Reader, opts ...excelize.Options) FileReader {
	return newReadWorkbook(&readWorkbookOption{
		reader:  reader,
		options: opts,
	})
}
