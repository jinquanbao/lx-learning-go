package excel

import (
	"fmt"
	"io"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type readWorkbookOption struct {
	filePath string
	reader   io.Reader
	options  []excelize.Options
}

func (r *readWorkbookOption) Validate() error {
	if r.reader == nil && len(r.filePath) == 0 {
		return ErrInputFile
	}
	return nil
}

type readSheetOption struct {
	SheetNo          int
	SheetName        string
	TitleRow         int
	TitleBeginColumn int
	Dest             interface{}
}

func (r *readSheetOption) Validate(file *excelize.File) error {
	if r.Dest == nil {
		return ErrInputDest
	}
	if len(r.SheetName) > 0 {
		sheetIndex, err := file.GetSheetIndex(r.SheetName)
		if err != nil {
			return err
		}
		if sheetIndex == -1 {
			fmt.Println("excel sheet name [" + r.SheetName + "] is invalid or sheet doesn't exist")
			return ErrInvalidSheetName
		}
		r.SheetNo = sheetIndex
	} else {
		r.SheetName = file.GetSheetName(r.SheetNo)
		if len(r.SheetName) == 0 {
			fmt.Println("excel sheet index [" + strconv.Itoa(r.SheetNo) + "] is invalid or sheet doesn't exist")
			return ErrInvalidSheetNo
		}
	}
	if r.TitleRow < 0 {
		return ErrTitleRow
	}
	if r.TitleBeginColumn < 0 {
		return ErrTitleBeginColumn
	}
	return nil
}
