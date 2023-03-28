package excel

import (
	"errors"

	"github.com/xuri/excelize/v2"
)

type readWorkbook struct {
	file                  *excelize.File
	option                *readWorkbookOption
	readSheets            []*readSheet
	readCompleteCallbacks []ReadCompleteCallback
	readContext           *readContextOption
}

type ReadCompleteCallback func(ReadContext) error

func newReadWorkbook(option *readWorkbookOption) *readWorkbook {
	return &readWorkbook{
		option:     option,
		readSheets: make([]*readSheet, 0),
	}
}

func (r *readWorkbook) ReadCompleteCallbacks(readCompleteCallbacks ...ReadCompleteCallback) FileReader {
	r.readCompleteCallbacks = readCompleteCallbacks
	return r
}

func (r *readWorkbook) ReadSheets(readSheets ...*readSheet) Reader {
	r.readSheets = readSheets
	return r
}

func (r *readWorkbook) ReadSheetName(sheetName string, dest interface{}) *readSheet {
	return &readSheet{
		readWorkbook: r,
		option: &readSheetOption{
			SheetName: sheetName,
			Dest:      dest,
		},
	}
}

func (r *readWorkbook) ReadSheetNo(sheetNo int, dest interface{}) *readSheet {
	return &readSheet{
		readWorkbook: r,
		option: &readSheetOption{
			SheetNo: sheetNo,
			Dest:    dest,
		},
	}
}

func (r *readWorkbook) Read() (err error) {
	defer r.Close()

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

func (r *readWorkbook) preRead() (err error) {
	if err = r.option.Validate(); err != nil {
		return err
	}
	if len(r.readSheets) == 0 {
		return errors.New("excel reader not config doRead sheet")
	}
	if len(r.option.filePath) > 0 {
		if r.file, err = excelize.OpenFile(r.option.filePath, r.option.options...); err != nil {
			return err
		}
	} else {
		if r.file, err = excelize.OpenReader(r.option.reader, r.option.options...); err != nil {
			return err
		}
	}
	for _, v := range r.readSheets {
		if err = v.option.Validate(r.file); err != nil {
			return err
		}
	}

	r.readContext = &readContextOption{file: r.file}

	return nil
}

func (r *readWorkbook) Close() (err error) {
	if r.file != nil {
		err = r.file.Close()
	}

	return err
}
