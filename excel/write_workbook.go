package excelutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

type WriteWorkbook struct {
	File                   *excelize.File
	DefaultTitleStyleId    int
	DefaultValueStyleId    int
	DefaultTimeStyleId     int
	option                 *WriteWorkbookOption
	writeSheets            []WriteSheet
	writeCompleteCallbacks []WriteCompleteCallback
	writeContext           *writeContextOption
}
type (
	WriteOptions          func(*WriteWorkbookOption) error
	WriteCompleteCallback func(WriteContext) error
)

func newWriteWorkbook(option *WriteWorkbookOption) *WriteWorkbook {
	w := &WriteWorkbook{
		File:        excelize.NewFile(),
		option:      option,
		writeSheets: make([]WriteSheet, 0),
	}

	return w
}

func (w *WriteWorkbook) DisableAutoClose() *WriteWorkbook {
	w.option.DisableAutoClose = true
	return w
}

func (w *WriteWorkbook) WriteCompleteCallbacks(writeCompleteCallbacks ...WriteCompleteCallback) FileWriter {
	w.writeCompleteCallbacks = writeCompleteCallbacks
	return w
}

func (w *WriteWorkbook) WriteSheets(writeSheets ...WriteSheet) Writer {
	w.writeSheets = writeSheets
	return w
}

//func (w *WriteWorkbook) WriteSheet(sheetName string, dest interface{}) *writeSheet {
//	return &writeSheet{
//		WriteWorkbook: w,
//		option: &writeSheetOption{
//			SheetName: sheetName,
//		},
//		dest: dest,
//	}
//}

func (w *WriteWorkbook) StreamWriteSheet(sheetName string, dest interface{}) *writeSheet {
	return &writeSheet{
		writeWorkbook: w,
		option: &writeSheetOption{
			StreamWriter: true,
			SheetName:    sheetName,
		},
		dest: dest,
	}
}

func (w *WriteWorkbook) DynamicWriteSheet(sheetName string, titles, values [][]interface{}) *dynamicWriteSheet {
	return &dynamicWriteSheet{
		writeWorkbook: w,
		option: &writeSheetOption{
			StreamWriter: true,
			SheetName:    sheetName,
		},
		titles: titles,
		values: values,
	}
}

func (w *WriteWorkbook) Write() error {
	if !w.option.DisableAutoClose {
		defer func() {
			if err := w.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}
	if err := w.preWrite(); err != nil {
		return err
	}

	w.writeContext = &writeContextOption{file: w.File}

	for i := range w.writeSheets {
		if err := w.writeSheets[i].preWrite(); err != nil {
			return err
		}
	}

	for i := range w.writeSheets {
		if err := w.writeSheets[i].doWrite(); err != nil {
			return err
		}
	}

	for _, callback := range w.writeCompleteCallbacks {
		if err := callback(w.writeContext); err != nil {
			return err
		}
	}

	if len(w.option.SaveFilePath) > 0 {
		if err := w.File.SaveAs(w.option.SaveFilePath, w.option.saveFileOptions...); err != nil {
			return err
		}
	}

	return nil
}

func (w *WriteWorkbook) preWrite() (err error) {
	if err = w.option.Validate(); err != nil {
		return err
	}

	if len(w.writeSheets) == 0 {
		return errors.New("excel write not config doWrite sheet")
	}

	if len(w.option.SaveFilePath) > 0 {
		err := os.MkdirAll(filepath.Dir(w.option.SaveFilePath), os.ModePerm)
		if err != nil {
			return err
		}
	}
	for _, v := range w.writeSheets {
		if err = v.getOption().Validate(w.File); err != nil {
			return err
		}
	}

	w.DefaultTitleStyleId, err = defaultStyle.newDefaultTitleStyleId(w.File)
	if err != nil {
		return err
	}
	w.DefaultValueStyleId, err = defaultStyle.newDefaultValueStyleId(w.File)
	if err != nil {
		return err
	}
	w.DefaultTimeStyleId, err = defaultStyle.newDefaultTimeStyleId(w.File)
	if err != nil {
		return err
	}

	return nil
}

func (w *WriteWorkbook) Close() (err error) {
	if w.File != nil {
		err = w.File.Close()
	}

	return err
}
