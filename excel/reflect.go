package excelutil

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func scan(v reflect.Value, value string) (err error) {
	if value == "" {
		return nil
	}
	if !v.CanAddr() {
		return ErrReflectValueAddr
	}
	ptr := v.Addr().Interface()
	switch p := ptr.(type) {
	case *string:
		*p = value
	case *[]byte:
		*p = []byte(value)
	case *[]rune:
		*p = []rune(value)
	case *int:
		*p, err = strconv.Atoi(value)
	case *int8:
		val, e := strconv.ParseInt(value, 10, 8)
		err = e
		*p = int8(val)
	case *int16:
		val, e := strconv.ParseInt(value, 10, 16)
		err = e
		*p = int16(val)
	case *int32:
		val, e := strconv.ParseInt(value, 10, 32)
		err = e
		*p = int32(val)
	case *int64:
		*p, err = strconv.ParseInt(value, 10, 64)
	case *uint:
		val, e := strconv.ParseUint(value, 10, 64)
		err = e
		*p = uint(val)
	case *uint8:
		val, e := strconv.ParseUint(value, 10, 8)
		err = e
		*p = uint8(val)
	case *uint16:
		val, e := strconv.ParseUint(value, 10, 16)
		err = e
		*p = uint16(val)
	case *uint32:
		val, e := strconv.ParseUint(value, 10, 32)
		err = e
		*p = uint32(val)
	case *uint64:
		*p, err = strconv.ParseUint(value, 10, 64)
	case *float32:
		val, e := strconv.ParseFloat(value, 32)
		err = e
		*p = float32(val)
	case *float64:
		*p, err = strconv.ParseFloat(value, 64)
	case *bool:
		*p, err = strconv.ParseBool(value)
	case *time.Time:
		timeNum, e := strconv.ParseFloat(value, 10)
		if e != nil {
			// return e
			*p, err = convertTime(value)
		} else {
			*p, err = excelize.ExcelDateToTime(timeNum, false)
		}

	default:
		err = ErrReflectValueType
	}
	return err
}

func convertTime(timeStrInput string) (res time.Time, err error) {
	if len(timeStrInput) > 0 {
		timeStr := strings.ReplaceAll(timeStrInput, "/", "")
		timeStr = strings.ReplaceAll(timeStr, "-", "")
		if len(timeStr) == 8 {
			res, err = parseTime("20060102", timeStr)
		} else if len(timeStr) == 11 {
			res, err = parseTime("20060102 15", timeStr)
		} else if len(timeStr) == 14 {
			res, err = parseTime("20060102 15:04", timeStr)
		} else if len(timeStr) == 17 {
			res, err = parseTime("20060102 15:04:05", timeStr)
		} else if len(timeStr) == 21 {
			res, err = parseTime("20060102 15:04:05.000", timeStr)
		} else {
			err = errors.New("bad value for timeStr")
		}
		if err != nil {
			timeStr = strings.ReplaceAll(timeStrInput, "-", "/")
			if len(timeStr) == 8 {
				res, err = parseTime("2006/1/2", timeStr)
			} else if len(timeStr) == 9 {
				res, err = parseTime("2006/01/2", timeStr)
				if err != nil {
					res, err = parseTime("2006/1/02", timeStr)
				}
			} else if len(timeStr) == 11 {
				res, err = parseTime("2006/1/2 15", timeStr)
			} else if len(timeStr) == 12 {
				res, err = parseTime("2006/01/2 15", timeStr)
				if err != nil {
					res, err = parseTime("2006/1/02 15", timeStr)
				}
			} else if len(timeStr) == 14 {
				res, err = parseTime("2006/1/2 15:04", timeStr)
			} else if len(timeStr) == 15 {
				res, err = parseTime("2006/1/02 15:04", timeStr)
				if err != nil {
					res, err = parseTime("2006/01/2 15:04", timeStr)
				}
			} else if len(timeStr) == 17 {
				res, err = parseTime("2006/1/2 15:04:05", timeStr)
			} else if len(timeStr) == 18 {
				res, err = parseTime("2006/1/02 15:04:05", timeStr)
				if err != nil {
					res, err = parseTime("2006/01/2 15:04:05", timeStr)
				}
			}
		}
	}
	//else {
	//	res, err = time.Parse("2006-01-02 15:04:05.000", "0001-01-01 00:00:00.000")
	//}
	return res, err
}

func isSlice(typ reflect.Type) bool {
	if typ.Kind() == reflect.Slice || (typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Slice) {
		return true
	}
	return false
}

var timeKind = reflect.TypeOf(time.Time{}).Kind()

func isTime(typ reflect.Type) bool {
	if typ.Kind() == timeKind || (typ.Kind() == reflect.Ptr && typ.Elem().Kind() == timeKind) {
		return true
	}
	return false
}

func parseTime(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, time.Local)
}
