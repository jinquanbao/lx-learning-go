package excel

import "errors"

var (
	ErrInputFile              = errors.New("reader or filePath must one not be null")
	ErrInputDest              = errors.New("excel reader dest must be ptr to slice")
	ErrInputDestElem          = errors.New("excel reader dest elem must be struct")
	ErrTitleRow               = errors.New("title row must grater than or equal to 0")
	ErrTitleBeginColumn       = errors.New("title begin column must grater than or equal to 0")
	ErrTitleNotMatch          = errors.New("title not match")
	ErrTitleDuplicate         = errors.New("duplicate title name which is not slice type")
	ErrInvalidSheetNo         = errors.New("sheet name is invalid or sheet doesn't exist")
	ErrInvalidSheetName       = errors.New("sheet index is invalid or sheet doesn't exist")
	ErrParseTag               = errors.New("parse excel tag err")
	ErrExcelTagIndexDuplicate = errors.New("excel tag's index duplicate")
	ErrInvalidExcelTagCell    = errors.New("excel tag's cell is invalid")
	ErrReflectValueAddr       = errors.New("reflect Value must be can addr")
	ErrReflectValueType       = errors.New("reflect Value type not support")
)
