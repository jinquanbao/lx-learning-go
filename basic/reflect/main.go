package main

import (
	"fmt"
	"reflect"
	"time"
)

/***

Bool
    Int
    Int8
    Int16
    Int32
    Int64
    Uint
    Uint8
    Uint16
    Uint32
    Uint64
    Uintptr
    Float32
    Float64
    Complex64
    Complex128
    Array
    Chan
    Func
    Interface
    Map
    Ptr
    Slice
    String
    Struct

*/
func main() {

	typeOf()

	valueOf := reflect.ValueOf(&FailRecordEntity{})

	methodByName := valueOf.MethodByName("TableName")
	values := methodByName.Call(nil)

	fmt.Println("tablename=", values)

	m := valueOf.MethodByName("SetInfo")

	var i []reflect.Value
	var j int64 = 1
	i = append(i, reflect.ValueOf(j))
	i = append(i, reflect.ValueOf("this is a message !"))
	call := m.Call(i)

	fmt.Println("id=", valueOf.Elem().FieldByName("Id").Int())
	fmt.Println("errorMsg=", valueOf.Elem().FieldByName("ErrorMsg").String())

	fmt.Println("call=", call)

}

func typeOf() error {
	var a float32 = 1.0
	var b int = 2.0
	var c int32 = 2.0
	var d chan int

	ofa := reflect.TypeOf(a)
	ofb := reflect.TypeOf(b)
	ofc := reflect.TypeOf(c)
	ofd := reflect.TypeOf(d)

	fmt.Printf("typeOf=%v", ofa)
	fmt.Printf("typeOf=%v", ofb.Name() == "int")
	fmt.Printf("typeOf=%v", ofc.Kind() != reflect.Map)
	fmt.Printf("typeOf=%v", ofd.Align())
	fmt.Printf("typeOf=%v", ofb.PkgPath())
	fmt.Println()
	typeOf := reflect.TypeOf(&FailRecordEntity{})

	fmt.Println("typeOf", typeOf)
	fmt.Println("Kind", typeOf.Kind())
	fmt.Println("Kind", typeOf.Elem().Kind())

	//判断是否是结构体
	fmt.Println("isStruct", typeOf.Elem().Kind() == reflect.Struct || typeOf.Kind() == reflect.Struct)
	if typeOf.Elem().Kind() == reflect.Struct {
		typeOf = typeOf.Elem()
	}

	//获取字段
	field0 := typeOf.Field(0)
	fmt.Println("fieldName=", field0.Name)
	fmt.Println("fieldType=", field0.Type)
	fmt.Println("fieldTag=", field0.Tag.Get("json"))

	//根据名称获取字段
	fieldByName, _ := typeOf.FieldByName("MsgId")
	fmt.Println("fieldByname", fieldByName.Name)
	fmt.Println("fieldType", fieldByName.Type)
	fmt.Println("fieldIndex", fieldByName.Index)
	fmt.Println("fieldTag", fieldByName.Tag.Get("db"))

	//获取方法
	method := typeOf.Method(0)
	fmt.Println("methodname", method.Name)
	fmt.Println("methodType", method.Type)
	fmt.Println("methodFunc", method.Func)
	return nil
}

func reflectSetValue(x interface{}) error {
	v := reflect.ValueOf(x)
	k := v.Kind()
	switch k {
	case reflect.Int64:
		// v.Int()从反射中获取整型的原始值，然后通过int64()强制类型转换
		fmt.Printf("type is int64, value is %d\n", int64(v.Int()))
		v.Elem().SetInt(1)
	case reflect.Float32:
		// v.Float()从反射中获取浮点型的原始值，然后通过float32()强制类型转换
		fmt.Printf("type is float32, value is %f\n", float32(v.Float()))
		v.Elem().SetFloat(1.1)
	case reflect.Float64:
		// v.Float()从反射中获取浮点型的原始值，然后通过float64()强制类型转换
		fmt.Printf("type is float64, value is %f\n", float64(v.Float()))
		v.Elem().SetFloat(2.1)

	}
	panic("not support kind")
}

type FailRecordEntity struct {
	Id         int64     `json:"id" validate:"-" db:"id"`                  // 主键id
	MsgId      string    `json:"msgId" validate:"-" db:"msg_id"`           // 消息id
	Type       string    `json:"type" validate:"-" db:"type"`              // 消息类型
	EntityId   string    `json:"entityId" validate:"-" db:"entity_id"`     // 业务entity_id
	EntityObj  string    `json:"entityObj" validate:"-" db:"entity_obj"`   // 消息体
	ErrorCode  string    `json:"errorCode" validate:"-" db:"error_code"`   // 错误编码
	ErrorMsg   string    `json:"errorMsg" validate:"-" db:"error_msg"`     // 错误原因
	CreateTime time.Time `json:"createTime" validate:"-" db:"create_time"` // 创建时间
	UpdateTime time.Time `json:"updateTime" validate:"-" db:"update_time"` // 更新时间
	Deleted    int8      `json:"deleted" validate:"-" db:"deleted"`        // 是否删除
}

func (f FailRecordEntity) TableName() string {
	return "fail_record"
}

func (f *FailRecordEntity) SetInfo(id int64, msg string) error {
	f.Id = id
	f.ErrorMsg = msg
	return nil
}
