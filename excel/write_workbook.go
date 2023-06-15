package excel

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
)

type writeWorkbook struct {
	File                   *excelize.File
	DefaultTitleStyleId    int
	DefaultValueStyleId    int
	DefaultTimeStyleId     int
	option                 *WriteWorkbookOption
	writeSheets            []*writeSheet
	writeCompleteCallbacks []WriteCompleteCallback
	writeContext           *writeContextOption
}
type (
	WriteOptions          func(*WriteWorkbookOption) error
	WriteCompleteCallback func(WriteContext) error
)

func newWriteWorkbook(option *WriteWorkbookOption) *writeWorkbook {
	w := &writeWorkbook{
		File:        excelize.NewFile(),
		option:      option,
		writeSheets: make([]*writeSheet, 0),
	}

	return w
}

func (w *writeWorkbook) WriteCompleteCallbacks(writeCompleteCallbacks ...WriteCompleteCallback) FileWriter {
	w.writeCompleteCallbacks = writeCompleteCallbacks
	return w
}

func (w *writeWorkbook) WriteSheets(writeSheets ...*writeSheet) Writer {
	w.writeSheets = writeSheets
	return w
}

//func (w *writeWorkbook) WriteSheet(sheetName string, dest interface{}) *writeSheet {
//	return &writeSheet{
//		writeWorkbook: w,
//		option: &writeSheetOption{
//			SheetName: sheetName,
//		},
//		dest: dest,
//	}
//}

func (w *writeWorkbook) StreamWriteSheet(sheetName string, dest interface{}) *writeSheet {
	return &writeSheet{
		writeWorkbook: w,
		option: &writeSheetOption{
			StreamWriter: true,
			SheetName:    sheetName,
		},
		dest: dest,
	}
}

func (w *writeWorkbook) Write() error {
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

func (w *writeWorkbook) preWrite() (err error) {
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
		_, err = os.Create(w.option.SaveFilePath)
		if err != nil {
			return err
		}
	}
	for _, v := range w.writeSheets {
		if err = v.option.Validate(w.File); err != nil {
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

func (w *writeWorkbook) Close() (err error) {
	if w.File != nil {
		err = w.File.Close()
	}

	return err
}
