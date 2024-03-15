package excelutil

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type ReadSheet struct {
	readWorkbook               *ReadWorkbook
	option                     *readSheetOption
	readSheetCompleteCallbacks []ReadSheetCompleteCallback
	readCellCompleteCallbacks  []ReadCellCompleteCallback
	listeners                  []ReadListener
	readSheetContext           *readSheetContextOption
	readCellContext            *readCellContextOption
	destIndirectValue          reflect.Value
	elemType                   reflect.Type
	elemMapValueIsInterface    bool
	sourceTitles               []string
	titleIndexMap              map[int]string
	schema                     *schema
	indexElemMap               map[int]elem
	indexTimeColumnMap         map[int]*column
	importTitles               []string
	importFieldNames           []string
}

type (
	// interface is option's dest without ptr
	ReadSheetCompleteCallback func(rCtx ReadSheetContext, dest interface{}) error
	// interface is option's dest elem with ptr
	ReadCellCompleteCallback func(rCtx ReadCellContext, destElem interface{}, err error) error
)

func (r *ReadSheet) TitleRow(titleRow int) *ReadSheet {
	r.option.TitleRow = titleRow
	return r
}

func (r *ReadSheet) TitleBeginColumn(titleBeginColumn int) *ReadSheet {
	r.option.TitleBeginColumn = titleBeginColumn
	return r
}

func (r *ReadSheet) ContentBeginRow(contentBeginRow int) *ReadSheet {
	r.option.ContentBeginRow = contentBeginRow
	return r
}

func (r *ReadSheet) DisableAutoMatchTitleLength(disableAutoMatchTitleLength bool) *ReadSheet {
	r.option.DisableAutoMatchTitleLength = disableAutoMatchTitleLength
	return r
}

func (r *ReadSheet) TimeTitles(timeTitles ...string) *ReadSheet {
	r.option.TimeTitles = timeTitles
	return r
}

func (r *ReadSheet) RegisterReadSheetCompleteCallbacks(opts ...ReadSheetCompleteCallback) *ReadSheet {
	r.readSheetCompleteCallbacks = opts
	return r
}

func (r *ReadSheet) RegisterReadCellCompleteCallbacks(opts ...ReadCellCompleteCallback) *ReadSheet {
	r.readCellCompleteCallbacks = opts
	return r
}

func (r *ReadSheet) Read() (err error) {
	return r.readWorkbook.ReadSheets(r).Read()
}

func (r *ReadSheet) preRead() (err error) {
	dest := r.option.Dest

	val := reflect.ValueOf(dest)
	typ := reflect.Indirect(val).Type()

	if val.Kind() != reflect.Ptr {
		return ErrInputDest
	}
	if typ.Kind() != reflect.Slice {
		return ErrInputDest
	}

	elemTyp := typ.Elem()
	if elemTyp.Kind() == reflect.Ptr {
		elemTyp = elemTyp.Elem()
	}

	r.elemType = elemTyp
	r.destIndirectValue = val.Elem()

	switch elemTyp.Kind() {
	case reflect.Struct:
		schema, err := newSchema(r.destIndirectValue)
		if err != nil {
			return err
		}
		if err = schema.initialization(); err != nil {
			return err
		}
		r.schema = schema
		if err = r.readTitle(); err != nil {
			return err
		}
		if err = r.computeStructTitle(); err != nil {
			return err
		}
		break
	case reflect.Map:
		if elemTyp.Key().Kind() != reflect.String {
			return ErrInputDestMapElem
		}
		elemMapValueType := elemTyp.Elem()
		if elemMapValueType.Kind() != reflect.String && elemMapValueType.Kind() != reflect.Interface {
			return ErrInputDestMapElem
		}
		r.elemMapValueIsInterface = elemMapValueType.Kind() == reflect.Interface
		if err = r.readTitle(); err != nil {
			return err
		}
		break
	default:
		return ErrInputDestElem
	}

	r.readSheetContext = &readSheetContextOption{
		file:             r.readWorkbook.File,
		sheetNo:          r.option.SheetNo,
		sheetName:        r.option.SheetName,
		titleRow:         r.option.TitleRow,
		titleBeginColumn: r.option.TitleBeginColumn,
		fieldNames:       r.importFieldNames,
	}

	r.readCellContext = &readCellContextOption{
		file:             r.readWorkbook.File,
		readSheetContext: r.readSheetContext,
		lastColumnIndex:  len(r.sourceTitles) - 1,
	}

	if elemTyp.Kind() == reflect.Map {
		for i := r.option.TitleBeginColumn; i < len(r.sourceTitles); i++ {
			r.readCellContext.rowIndex = r.option.TitleRow
			r.readCellContext.columnIndex = i
			r.readCellContext.titleName = r.sourceTitles[i]
			r.readCellContext.cellValue = r.sourceTitles[i]
			for _, callback := range r.readCellCompleteCallbacks {
				if err = callback(r.readCellContext, make(map[string]interface{}), err); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *ReadSheet) readTitle() (err error) {
	file := r.readWorkbook.File

	rows, err := file.Rows(r.option.SheetName)
	if err != nil {
		return err
	}

	titles, cur := make([]string, 0), 0
	for rows.Next() {
		if cur == r.option.TitleRow {
			titles, err = rows.Columns(r.readWorkbook.option.options...)
			if err != nil {
				return err
			}
			break
		}
		cur++
	}

	titleIndexMap := make(map[int]string)
	for i := r.option.TitleBeginColumn; i < len(titles); i++ {
		v := strings.TrimSpace(titles[i])
		titles[i] = v
		if v == "" {
			continue
		}
		r.importTitles = append(r.importTitles, v)
		titleIndexMap[i] = v
	}
	r.sourceTitles = titles
	r.titleIndexMap = titleIndexMap

	return nil
}

func (r *ReadSheet) computeStructTitle() error {
	var (
		titleIndexMap = r.titleIndexMap
		schema        = r.schema
		titleCount    = len(r.titleIndexMap)
	)
	titleMap := make(map[string][]int)
	for i, v := range r.sourceTitles {
		if v == "" {
			continue
		}
		if index, ok := titleMap[v]; ok {
			index = append(index, i)
			titleMap[v] = index
		} else {
			titleMap[v] = []int{i}
		}
	}

	indexColumnMap := schema.IndexColumnMap
	indexElemMap := schema.IndexElemMap
	titleElemMap := schema.TitleElemMap
	indexTimeColumnMap := make(map[int]*column)
	for i := range schema.Columns {

		column := schema.Columns[i]
		if v, ok := titleMap[column.Name]; ok {

			if column.Index != -1 {
				match := false
				for _, index := range v {
					if _, yes := indexColumnMap[index]; yes {
						match = true
						break
					}
				}
				if !match {
					fmt.Println("excel column name [" + column.Name + "] title index [" + strconv.Itoa(v[0]) + "] not match struct index [" + strconv.Itoa(column.Index) + "]")
					return ErrTitleNotMatch
				}
			}

			for _, index := range v {

				if val, yes := indexColumnMap[index]; yes {
					column = val
				} else {

					// 1. struct field not config index
					// 2. field of struct is slice or ptr to slice
					indexColumnMap[index] = column

					// 兼容相同的标题名称配置在不同的结构体上
					titleElem := elem{
						elemType:  column.Schema.ElemType,
						titleName: column.Name,
					}

					if indexList, yes := titleElemMap[titleElem]; yes {

						if !isSlice(column.Schema.IndirectType) {
							fmt.Println("duplicate title name [" + column.Name + "]  which is not slice type")
							return ErrTitleNotMatch
						}
						indexElemMap[index] = elem{
							elemType:  column.Schema.ElemType,
							elemIndex: len(indexList),
						}
						indexList = append(indexList, index)
						titleElemMap[titleElem] = indexList

					} else {

						titleElemMap[titleElem] = []int{index}
						indexElemMap[index] = elem{
							elemType:  column.Schema.ElemType,
							elemIndex: 0,
						}
					}
				}

				if column.FieldIsTime {
					indexTimeColumnMap[index] = column
				}
			}

		} else if !r.option.DisableAutoMatchTitleLength {
			fmt.Println("excel column name [" + column.Name + "] not exist")
			return ErrTitleNotMatch
		}
	}

	if titleCount > len(indexColumnMap) || titleCount == 0 {
		return ErrTitleNotMatch
	} else if titleCount != len(indexColumnMap) && !r.option.DisableAutoMatchTitleLength {
		fmt.Println("excel title cells not match struct fields")
		return ErrTitleNotMatch
	}

	schema.IndexColumnMap = indexColumnMap
	r.indexElemMap = indexElemMap
	r.indexTimeColumnMap = indexTimeColumnMap

	for index := range titleIndexMap {
		fileName := indexColumnMap[index].FieldName
		if len(fileName) > 0 && !containStr(r.importFieldNames, fileName) {
			r.importFieldNames = append(r.importFieldNames, fileName)
		}
	}
	return nil
}

func (r *ReadSheet) doRead() (err error) {
	file := r.readWorkbook.File

	rows, err := file.GetRows(r.option.SheetName, r.readWorkbook.option.options...)
	if err != nil {
		return err
	}

	r.destIndirectValue.Set(reflect.MakeSlice(r.destIndirectValue.Type(), 0, len(rows)))
	contentBeginRow := r.option.ContentBeginRow
	if contentBeginRow <= r.option.TitleRow {
		contentBeginRow = r.option.TitleRow + 1
	}

	switch r.elemType.Kind() {
	case reflect.Struct:
		if err = r.readToSliceStruct(contentBeginRow, rows); err != nil {
			return err
		}
		break
	case reflect.Map:
		if err = r.readToSliceMap(contentBeginRow, rows); err != nil {
			return err
		}
		break
	default:
		return ErrInputDestElem
	}

	for i := range r.listeners {
		if err = r.listeners[i].ReadCompleteTrigger(r.readSheetContext, r.destIndirectValue.Interface()); err != nil {
			return err
		}
	}

	for _, callback := range r.readSheetCompleteCallbacks {
		if err = callback(r.readSheetContext, r.destIndirectValue.Interface()); err != nil {
			return err
		}
	}

	return err
}

func (r *ReadSheet) readToSliceStruct(contentBeginRow int, rows [][]string) (err error) {
	file := r.readWorkbook.File
	// excel 日期格式数字值
	numFmtTimeStyle, err := file.NewStyle(&excelize.Style{NumFmt: 0})
	if err != nil {
		return err
	}

	for i := contentBeginRow; i < len(rows); i++ {

		r.readCellContext.rowIndex = i

		columns := rows[i]

		timeStyleMap := make(map[int]int)
		for col := range r.indexTimeColumnMap {
			if len(columns) > col {
				cellName, _ := excelize.CoordinatesToCellName(col+1, i+1)

				style, _ := file.GetCellStyle(r.option.SheetName, cellName)
				timeStyleMap[col] = style
				file.SetCellStyle(r.option.SheetName, cellName, cellName, numFmtTimeStyle)

				if cellTime, err := file.GetCellValue(r.option.SheetName, cellName); err == nil {
					columns[col] = cellTime
				}
			}
		}

		if err = r.readToStruct(columns); err != nil {
			return err
		}

		// 恢复原有样式
		for col := range r.indexTimeColumnMap {
			if len(columns) > col {
				if style, ok := timeStyleMap[col]; ok {
					cellName, _ := excelize.CoordinatesToCellName(col+1, i+1)
					file.SetCellStyle(r.option.SheetName, cellName, cellName, style)
				}
			}
		}
	}
	return nil
}

func (r *ReadSheet) readToStruct(columns []string) (err error) {
	sourceVal := r.schema.getNextElemValue(r.destIndirectValue)

	elemValueMap := make(map[elem]reflect.Value)

	for i, v := range columns {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		r.readCellContext.columnIndex = i
		if column, ok := r.schema.IndexColumnMap[i]; ok {
			if elem, yes := r.indexElemMap[i]; yes {

				var elemValue reflect.Value

				if value, yes := elemValueMap[elem]; yes {
					elemValue = value
				} else {
					elemValue = r.schema.getElemValue(sourceVal, column.Schema)
					elemValueMap[elem] = elemValue
				}

				if column.Converter != "" {
					if method := elemValue.Addr().MethodByName(column.Converter); method.IsValid() {
						result := method.Call([]reflect.Value{reflect.ValueOf(v)})
						for _, v := range result {
							if val, yes := v.Interface().(error); yes {
								err = val
								break
							}
						}
					}
				} else {
					fieldValue := elemValue.Field(column.FieldIndex)
					err = scan(fieldValue, v)
				}

				r.readCellContext.titleName = column.Name
				r.readCellContext.cellValue = v

				for i := range r.listeners {
					if err = r.listeners[i].ReadCellCompleteTrigger(r.readCellContext, sourceVal.Addr().Interface(), err); err != nil {
						return err
					}
				}

				for _, callback := range r.readCellCompleteCallbacks {
					if err = callback(r.readCellContext, sourceVal.Addr().Interface(), err); err != nil {
						return err
					}
				}

				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *ReadSheet) readToSliceMap(contentBeginRow int, rows [][]string) (err error) {
	for i := contentBeginRow; i < len(rows); i++ {
		r.readCellContext.rowIndex = i
		if err = r.readToMap(i, rows[i]); err != nil {
			return err
		}
	}
	return nil
}

func (r *ReadSheet) readToMap(rowIndex int, columns []string) (err error) {
	var indexTimeList []int
	for i, titleName := range r.sourceTitles {
		if containStr(r.option.TimeTitles, titleName) {
			indexTimeList = append(indexTimeList, i)
		}
	}

	sourceVal := r.schema.getNextElemValue(r.destIndirectValue)
	sourceVal.Set(reflect.MakeMapWithSize(r.elemType, len(r.titleIndexMap)))

	for i := r.option.TitleBeginColumn; i < len(columns); i++ {
		v := strings.TrimSpace(columns[i])
		titleName, ok := r.titleIndexMap[i]
		if ok {
			if containInt(indexTimeList, i) && r.elemMapValueIsInterface {
				timeVal, err := r.parseTimeValue(rowIndex, i)
				if err != nil {
					sourceVal.SetMapIndex(reflect.ValueOf(titleName), reflect.ValueOf(v))
				} else {
					sourceVal.SetMapIndex(reflect.ValueOf(titleName), reflect.ValueOf(timeVal))
				}
			} else {
				sourceVal.SetMapIndex(reflect.ValueOf(titleName), reflect.ValueOf(v))
			}
		}

		r.readCellContext.columnIndex = i
		r.readCellContext.titleName = titleName
		r.readCellContext.cellValue = v

		for i := range r.listeners {
			if err = r.listeners[i].ReadCellCompleteTrigger(r.readCellContext, sourceVal.Addr().Interface(), err); err != nil {
				return err
			}
		}
		for _, callback := range r.readCellCompleteCallbacks {
			if err = callback(r.readCellContext, sourceVal.Addr().Interface(), err); err != nil {
				return err
			}
		}
	}

	for i := len(columns); i < len(r.sourceTitles); i++ {
		titleName, ok := r.titleIndexMap[i]
		if !ok {
			continue
		}
		// fill default value to column not read.
		sourceVal.SetMapIndex(reflect.ValueOf(titleName), reflect.ValueOf(""))
	}

	return nil
}

func (r *ReadSheet) parseTimeValue(rowIndex, col int) (res time.Time, err error) {
	var cellTime string
	cellName, _ := excelize.CoordinatesToCellName(col+1, rowIndex+1)
	styleId, _ := r.readWorkbook.File.GetCellStyle(r.option.SheetName, cellName)
	_ = r.readWorkbook.File.SetCellStyle(r.option.SheetName, cellName, cellName, r.readWorkbook.NumFmtTimeStyleId)

	defer func() {
		// 恢复原有样式
		_ = r.readWorkbook.File.SetCellStyle(r.option.SheetName, cellName, cellName, styleId)
	}()

	if cellTime, err = r.readWorkbook.File.GetCellValue(r.option.SheetName, cellName); err != nil {
		return res, err
	}

	if cellTime == "" {
		return res, nil
	}
	timeNum, e := strconv.ParseFloat(cellTime, 10)
	if e != nil {
		return convertTime(cellTime)
	} else {
		return excelize.ExcelDateToTime(timeNum, false)
	}
}

func (r *ReadSheet) Close() (err error) {
	return nil
}
