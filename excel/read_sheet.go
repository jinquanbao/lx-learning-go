package excel

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type readSheet struct {
	readWorkbook               *readWorkbook
	option                     *readSheetOption
	readSheetCompleteCallbacks []ReadSheetCompleteCallback
	readCellCompleteCallbacks  []ReadCellCompleteCallback
	listeners                  []ReadListener
	readSheetContext           *readSheetContextOption
	readCellContext            *readCellContextOption
	destIndirectValue          reflect.Value
	schema                     *schema
	indexElemMap               map[int]elem
	indexTimeColumnMap         map[int]*column
}

type (
	// interface is option's dest without ptr
	ReadSheetCompleteCallback func(rCtx ReadSheetContext, dest interface{}) error
	// interface is option's dest elem with ptr
	ReadCellCompleteCallback func(rCtx ReadCellContext, destElem interface{}, err error) error
)

func (r *readSheet) TitleRow(titleRow int) *readSheet {
	r.option.TitleRow = titleRow
	return r
}

func (r *readSheet) TitleBeginColumn(titleBeginColumn int) *readSheet {
	r.option.TitleBeginColumn = titleBeginColumn
	return r
}

func (r *readSheet) RegisterReadSheetCompleteCallbacks(opts ...ReadSheetCompleteCallback) *readSheet {
	r.readSheetCompleteCallbacks = opts
	return r
}

func (r *readSheet) RegisterReadCellCompleteCallbacks(opts ...ReadCellCompleteCallback) *readSheet {
	r.readCellCompleteCallbacks = opts
	return r
}

//func (r *readSheet) RegisterReadListener(listeners ...ReadListener) *readSheet {
//	r.listeners = listeners
//	return r
//}

func (r *readSheet) Read() (err error) {
	return r.readWorkbook.ReadSheets(r).Read()
}

func (r *readSheet) preRead() (err error) {
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

	if elemTyp.Kind() != reflect.Struct {
		return ErrInputDestElem
	}

	r.destIndirectValue = val.Elem()

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

	r.readSheetContext = &readSheetContextOption{
		file:             r.readWorkbook.file,
		sheetNo:          r.option.SheetNo,
		sheetName:        r.option.SheetName,
		titleRow:         r.option.TitleRow,
		titleBeginColumn: r.option.TitleBeginColumn,
	}

	r.readCellContext = &readCellContextOption{
		file:             r.readWorkbook.file,
		readSheetContext: r.readSheetContext,
	}

	return nil
}

func (r *readSheet) readTitle() (err error) {
	file := r.readWorkbook.file

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
	titleCount := 0
	for i := r.option.TitleBeginColumn; i < len(titles); i++ {
		v := strings.TrimSpace(titles[i])
		titles[i] = v
		if v == "" {
			continue
		}
		titleCount++
	}

	titleMap := make(map[string][]int)
	for i, v := range titles {
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

	schema := r.schema

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

		} else {
			fmt.Println("excel column name [" + column.Name + "] not exist")
			return ErrTitleNotMatch
		}
	}

	if titleCount != len(indexColumnMap) {
		fmt.Println("excel title cells not match struct fields")
		return ErrTitleNotMatch
	}

	schema.IndexColumnMap = indexColumnMap
	r.indexElemMap = indexElemMap
	r.indexTimeColumnMap = indexTimeColumnMap

	return nil
}

func (r *readSheet) doRead() (err error) {
	file := r.readWorkbook.file

	rows, err := file.GetRows(r.option.SheetName, r.readWorkbook.option.options...)
	if err != nil {
		return err
	}
	// fmt.Println("excel read complete " + time.Now().String())
	r.destIndirectValue.Set(reflect.MakeSlice(r.destIndirectValue.Type(), 0, len(rows)))

	// excel 日期格式数字值
	numFmtTimeStyle, err := file.NewStyle(&excelize.Style{NumFmt: 0})
	timeStyleMap := make(map[int]int)
	for col := range r.indexTimeColumnMap {
		cellName, _ := excelize.CoordinatesToCellName(col+1, r.option.TitleRow+2)
		style, _ := file.GetCellStyle(r.option.SheetName, cellName)
		timeStyleMap[col] = style
		vCell, _ := excelize.CoordinatesToCellName(col+1, len(rows))
		file.SetCellStyle(r.option.SheetName, cellName, vCell, numFmtTimeStyle)
	}

	for i := r.option.TitleRow + 1; i < len(rows); i++ {

		r.readCellContext.rowIndex = i

		columns := rows[i]

		for col := range r.indexTimeColumnMap {
			if len(columns) > col {
				cellName, _ := excelize.CoordinatesToCellName(col+1, i+1)
				if cellTime, err := file.GetCellValue(r.option.SheetName, cellName); err == nil {
					columns[col] = cellTime
				}
			}
		}

		if err = r.readRow(columns); err != nil {
			return err
		}
	}

	// 恢复原有样式
	for col := range r.indexTimeColumnMap {
		if style, ok := timeStyleMap[col]; ok {
			cellName, _ := excelize.CoordinatesToCellName(col+1, r.option.TitleRow+2)
			vCell, _ := excelize.CoordinatesToCellName(col+1, len(rows))
			file.SetCellStyle(r.option.SheetName, cellName, vCell, style)
		}
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

	return nil
}

func (r *readSheet) readRow(columns []string) (err error) {
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

func (r *readSheet) Close() (err error) {
	return nil
}
