package excelutil

import "github.com/xuri/excelize/v2"

type WriteWorkbookOption struct {
	SaveFilePath     string
	saveFileOptions  []excelize.Options
	DisableAutoClose bool
}

func (r *WriteWorkbookOption) Validate() error {
	return nil
}

type writeSheetOption struct {
	SheetNo                 int
	SheetName               string
	StreamWriter            bool
	Titles                  []string // dynamic Titles
	TitleRow                int
	ContentBeginRow         int
	TitleBeginColumn        int
	IncludeTitleNames       []string
	ExcludeTitleNames       []string
	IncludeColumnFieldNames []string
	ExcludeColumnFieldNames []string
}

func (w *writeSheetOption) Validate(file *excelize.File) (err error) {
	if w.TitleRow < 0 {
		return ErrTitleRow
	}
	if w.TitleBeginColumn < 0 {
		return ErrTitleBeginColumn
	}
	w.SheetNo, err = file.NewSheet(w.SheetName)
	if err != nil {
		return err
	}
	return nil
}
