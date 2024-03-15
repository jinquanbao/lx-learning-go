package excelutil

import (
	"errors"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ReadWorkbook struct {
	File                  *excelize.File
	NumFmtTimeStyleId     int
	option                *readWorkbookOption
	readSheets            []*ReadSheet
	readCompleteCallbacks []ReadCompleteCallback
	readContext           *readContextOption
}

type ReadCompleteCallback func(ReadContext) error

func newReadWorkbook(option *readWorkbookOption) *ReadWorkbook {
	return &ReadWorkbook{
		option:     option,
		readSheets: make([]*ReadSheet, 0),
	}
}

func (r *ReadWorkbook) DisableAutoClose() *ReadWorkbook {
	r.option.DisableAutoClose = true
	return r
}

func (r *ReadWorkbook) ReadCompleteCallbacks(readCompleteCallbacks ...ReadCompleteCallback) FileReader {
	r.readCompleteCallbacks = readCompleteCallbacks
	return r
}

func (r *ReadWorkbook) ReadSheets(readSheets ...*ReadSheet) Reader {
	r.readSheets = readSheets
	return r
}

func (r *ReadWorkbook) ReadSheetName(sheetName string, dest interface{}) *ReadSheet {
	return &ReadSheet{
		readWorkbook: r,
		option: &readSheetOption{
			SheetName: sheetName,
			Dest:      dest,
		},
	}
}

func (r *ReadWorkbook) ReadSheetNo(sheetNo int, dest interface{}) *ReadSheet {
	return &ReadSheet{
		readWorkbook: r,
		option: &readSheetOption{
			SheetNo: sheetNo,
			Dest:    dest,
		},
	}
}

func (r *ReadWorkbook) Read() (err error) {
	if !r.option.DisableAutoClose {
		defer func() {
			if err := r.Close(); err != nil {
				fmt.Println(err)
			}
		}()
	}

	if err = r.preRead(); err != nil {
		return err
	}

	for i := range r.readSheets {
		if err = r.readSheets[i].preRead(); err != nil {
			return err
		}
	}

	for i := range r.readSheets {
		if err = r.readSheets[i].doRead(); err != nil {
			return err
		}
	}

	for _, callback := range r.readCompleteCallbacks {
		if err = callback(r.readContext); err != nil {
			return err
		}
	}

	return nil
}

func (r *ReadWorkbook) initializationFile() (err error) {
	if r.File != nil {
		return nil
	}
	if len(r.option.filePath) > 0 {
		if r.File, err = excelize.OpenFile(r.option.filePath, r.option.options...); err != nil {
			return err
		}
	} else {
		if r.File, err = excelize.OpenReader(r.option.reader, r.option.options...); err != nil {
			return err
		}
	}
	return nil
}

func (r *ReadWorkbook) preRead() (err error) {
	if err = r.option.Validate(); err != nil {
		return err
	}
	if len(r.readSheets) == 0 {
		return errors.New("excel reader not config doRead sheet")
	}
	if err = r.initializationFile(); err != nil {
		return err
	}
	for _, v := range r.readSheets {
		if err = v.option.Validate(r.File); err != nil {
			return err
		}
	}

	r.NumFmtTimeStyleId, err = r.File.NewStyle(&excelize.Style{NumFmt: 0})
	if err != nil {
		return err
	}

	r.readContext = &readContextOption{file: r.File}

	return nil
}

func (r *ReadWorkbook) Close() (err error) {
	if r.File != nil {
		err = r.File.Close()
	}

	return err
}
