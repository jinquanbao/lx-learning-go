package excel

import (
	"fmt"
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

func StreamWrite(value interface{}, saveFilePath string, options ...excelize.Options) error {
	fileWriter := NewWriter(WithWriteSaveFilePath(saveFilePath))
	defer func() {
		if err := fileWriter.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	return fileWriter.StreamWriteSheet("Sheet1", value).Write()
}

func NewWriter(opts ...WriteOptions) *writeWorkbook {
	option := &WriteWorkbookOption{}
	for _, opt := range opts {
		opt(option)
	}
	return newWriteWorkbook(option)
}

func WithWriteSaveFilePath(saveFilePath string, options ...excelize.Options) WriteOptions {
	return func(option *WriteWorkbookOption) error {
		option.SaveFilePath = saveFilePath
		option.saveFileOptions = options
		return nil
	}
}

func WithWriteDisableAutoClose(disableAutoClose bool) WriteOptions {
	return func(option *WriteWorkbookOption) error {
		option.DisableAutoClose = disableAutoClose
		return nil
	}
}
