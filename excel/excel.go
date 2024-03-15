package excelutil

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

// OpenFile 延迟初始化 *excelize.File，调用 reader.Read() 方法时才会初始会*excelize.File
func OpenFile(filePath string, opts ...excelize.Options) *ReadWorkbook {
	return newReadWorkbook(&readWorkbookOption{
		filePath: filePath,
		options:  opts,
	})
}

// SmartOpenFile 禁用自动关闭，并优先初始化 *excelize.File
//
//	 reader,err := SmartOpenFile(f,opts)
//	 defer func() {
//			 if err := reader.Close(); err != nil {
//				log.FromContext(ctx).Errorf("%v", err)
//			 }
//		}()
//		if err != nil {
//			 return err
//		}
func SmartOpenFile(filePath string, opts ...excelize.Options) (*ReadWorkbook, error) {
	r := OpenFile(filePath, opts...).DisableAutoClose()
	if err := r.initializationFile(); err != nil {
		return r, err
	}
	return r, nil
}

// OpenReader 延迟初始化 *excelize.File，调用 reader.Read() 方法时才会初始会*excelize.File
func OpenReader(reader io.Reader, opts ...excelize.Options) *ReadWorkbook {
	return newReadWorkbook(&readWorkbookOption{
		reader:  reader,
		options: opts,
	})
}

// SmartOpenReader 禁用自动关闭，并优先初始化 *excelize.File，外部需要自己调用Close()方法
//
//	 reader,err := SmartOpenReader(r,opts)
//	 defer func() {
//			 if err := reader.Close(); err != nil {
//				log.FromContext(ctx).Errorf("%v", err)
//			 }
//		}()
//		if err != nil {
//			 return err
//		}
func SmartOpenReader(reader io.Reader, opts ...excelize.Options) (*ReadWorkbook, error) {
	r := OpenReader(reader, opts...).DisableAutoClose()
	if err := r.initializationFile(); err != nil {
		return r, err
	}
	return r, nil
}

func StreamWrite(value interface{}, saveFilePath string, options ...excelize.Options) error {
	fileWriter := NewWriter(WithWriteSaveFilePath(saveFilePath, options...))
	defer func() {
		if err := fileWriter.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	return fileWriter.StreamWriteSheet("Sheet1", value).Write()
}

func DynamicWrite(titles, values [][]interface{}, saveFilePath string, options ...excelize.Options) error {
	fileWriter := NewWriter(WithWriteSaveFilePath(saveFilePath, options...))
	defer func() {
		if err := fileWriter.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	return fileWriter.DynamicWriteSheet("Sheet1", titles, values).Write()
}

func NewWriter(opts ...WriteOptions) *WriteWorkbook {
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

func WithWriteDisableAutoClose() WriteOptions {
	return func(option *WriteWorkbookOption) error {
		option.DisableAutoClose = true
		return nil
	}
}
