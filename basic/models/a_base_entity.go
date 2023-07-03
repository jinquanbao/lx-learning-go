package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

type BaseEntity interface {
	TableName() string
}

func NewMyTime(t time.Time) MyTime {
	return MyTime{Time: t}
}

type MyTime struct {
	time.Time
}

// UnmarshalJSON ..
func (d *MyTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "0" {
		d.Time, err = time.Parse("2006-01-02 15:04:05.000", "0001-01-01 00:00:00.000")
		return err
	}

	if len(s) > 0 && len(s) <= 13 {
		millSec, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		d.Time = time.UnixMilli(millSec)
	} else {
		return errors.New("invalid timestamp")
	}

	return nil
}

// MarshalJSON ..
func (d MyTime) MarshalJSON() ([]byte, error) {
	if d.Time.Year() < 1800 {
		return json.Marshal(0)
	}
	return json.Marshal(d.Time.UnixMilli())
}

// Scan 从数据库读取进行绑定时用到
func (d *MyTime) Scan(b interface{}) (err error) {
	switch x := b.(type) {
	case time.Time:
		d.Time = x
	case []byte:
		t, err := time.Parse("2006-01-02 15:04:05.000", string(b.([]byte)))
		if err != nil {
			return err
		}
		d.Time = t
	default:
		d.Time = time.Time{}
	}
	return nil
}

// Value 写入数据库时用到
func (d MyTime) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return `0001-01-01 00:00:00`, nil
	}
	v := d.Time.Format("2006-01-02 15:04:05.000")
	return v, nil
}

func (d *MyTime) GetMillSec() int64 {
	if d.Time.IsZero() {
		return 0
	}

	return d.Time.UnixMilli()
}

// 是否 在dao层中使用，其他层建议使用constant.NO 其他层建议使用constant.YES
const (
	Yes = 1
	No  = 2
)

// 是否删除
const (
	DeleteYes = 1
	DeleteNo  = 2
)

// 是否启用
const (
	Enabled   = 1
	UnEnabled = 2
)
