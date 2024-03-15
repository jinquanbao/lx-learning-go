package excelutil

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

const (
	excelTag     = "excel"
	indexTag     = "index"
	cellTag      = "cell"
	nameTag      = "name"
	converterTag = "converter"
	ignoredTag   = "-"
)

type schema struct {
	Value          reflect.Value
	IndirectType   reflect.Type
	ElemType       reflect.Type
	Columns        []*column
	Parent         *schema
	FieldIndex     int
	IndexColumnMap map[int]*column
	IndexElemMap   map[int]elem
	TitleElemMap   map[elem][]int
}

type column struct {
	FieldIndex  int
	FieldName   string
	FieldIsTime bool
	Index       int
	Cell        string
	Name        string
	Converter   string
	Schema      *schema
}

type elem struct {
	elemType  reflect.Type
	elemIndex int
	titleName string
}

func newSchema(v reflect.Value) (s *schema, err error) {
	indirectType := reflect.Indirect(v).Type()
	elemType := indirectType

	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}
	if elemType.Kind() == reflect.Slice {
		elemType = elemType.Elem()
	}
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	if elemType.Kind() != reflect.Struct {
		return nil, err
	}
	s = &schema{
		Value:        v,
		IndirectType: indirectType,
		ElemType:     elemType,
	}

	err = s.parse()
	return s, err
}

func (s *schema) parse() (err error) {
	columns := make([]*column, 0, s.ElemType.NumField())
	for i := 0; i < s.ElemType.NumField(); i++ {
		field := s.ElemType.Field(i)
		if value, ok := field.Tag.Lookup(excelTag); ok {
			if value != ignoredTag {
				c, err := s.parseTagValue(value)
				if err != nil {
					fmt.Println("reflect type of " + s.ElemType.Name() + "'s field[" + field.Name + "] parse tag error:" + err.Error())
					return ErrParseTag
				}
				c.FieldIndex = i
				c.FieldName = field.Name
				c.FieldIsTime = isTime(field.Type)
				if len(c.Name) == 0 {
					c.Name = field.Name
				}
				columns = append(columns, c)
			}
		} else {
			children, err := newSchema(reflect.New(field.Type))
			if err != nil {
				return err
			}
			if children != nil {
				children.Parent = s
				children.FieldIndex = i
				columns = append(columns, children.Columns...)
			} else {
				c := s.newColumn()
				c.FieldIndex = i
				c.FieldName = field.Name
				c.FieldIsTime = isTime(field.Type)
				c.Name = field.Name
				columns = append(columns, c)
			}
		}
	}
	s.Columns = columns
	return nil
}

func (s *schema) newColumn() (c *column) {
	c = &column{
		Index:  -1,
		Schema: s,
	}
	return c
}

func (s *schema) parseTagValue(value string) (c *column, err error) {
	c = s.newColumn()
	params := strings.Split(value, ";")

	for _, param := range params {

		tagKeyWithValue := strings.Split(param, ":")
		if len(tagKeyWithValue) != 2 {
			continue
		}

		value := strings.TrimSpace(tagKeyWithValue[1])

		switch strings.TrimSpace(tagKeyWithValue[0]) {
		case nameTag:
			c.Name = value
		case cellTag:
			c.Cell = value
		case indexTag:
			if c.Index, err = strconv.Atoi(value); err != nil {
				return nil, err
			}
		case converterTag:
			c.Converter = value
		}

	}
	return c, nil
}

func (s *schema) initialization() error {
	indexColumnMap := make(map[int]*column)
	indexElemMap := make(map[int]elem)
	titleElemMap := make(map[elem][]int)
	for _, v := range s.Columns {

		index := v.Index
		if len(v.Cell) > 0 {
			col, _, err := excelize.CellNameToCoordinates(v.Cell + "1")
			if err != nil {
				return ErrInvalidExcelTagCell
			}
			index = col - 1
		}

		if v.Index != -1 {
			if _, ok := indexColumnMap[index]; ok {
				fmt.Println("excel column name [" + v.Name + "] index of " + strconv.Itoa(index) + " duplicate")
				return ErrExcelTagIndexDuplicate
			}
			indexColumnMap[index] = v
			indexElemMap[index] = elem{
				elemType:  v.Schema.ElemType,
				elemIndex: 0,
			}
			titleElemMap[elem{
				elemType:  v.Schema.ElemType,
				titleName: v.Name,
			}] = []int{index}
		}
	}
	s.IndexColumnMap = indexColumnMap
	s.IndexElemMap = indexElemMap
	s.TitleElemMap = titleElemMap

	return nil
}

func (s *schema) getNextElemValue(dest reflect.Value) reflect.Value {
	if dest.Type().Kind() == reflect.Ptr {
		dest = dest.Elem()
	}
	elemType := dest.Type().Elem()

	var elem reflect.Value

	if elemType.Kind() == reflect.Ptr {
		elemPtr := reflect.New(elemType.Elem())
		dest.Set(reflect.Append(dest, elemPtr))
		elem = elemPtr.Elem()
		return elem
	} else {
		elem = reflect.Zero(elemType)
		dest.Set(reflect.Append(dest, elem))
		return dest.Index(dest.Len() - 1)
	}
}

func (s *schema) getElemValue(source reflect.Value, schema *schema) reflect.Value {
	if schema.Parent == nil {
		return source
	}
	parentValue := s.getElemValue(source, schema.Parent)

	elemValue := parentValue.Field(schema.FieldIndex)

	if elemValue.Kind() == reflect.Ptr {
		if elemValue.IsNil() {
			elemValue.Set(reflect.New(schema.IndirectType.Elem()))
		}
		elemValue = elemValue.Elem()
	}

	if isSlice(schema.IndirectType) {
		elemValue = s.getNextElemValue(elemValue)
	}

	return elemValue
}

func (s *schema) MakeSlice(elemType reflect.Type, cap int) reflect.Value {
	slice := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(elemType)), 0, cap)
	results := reflect.New(slice.Type())
	results.Elem().Set(slice)
	return results
}
