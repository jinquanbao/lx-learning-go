package excel

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/jinquanbao/lx-learning-go/excel"
)

func TestSingleSheet(t *testing.T) {
	dest := &[]User{}
	err := excel.ReadFile("./example.xlsx", dest)
	if err != nil {
		println(err.Error())
	}
	r, err := json.Marshal(dest)
	if err != nil {
		println(err.Error())
	}
	println(string(r))
}

func TestMultipleSheet(t *testing.T) {
	dest := &[]ExcelExample{}
	destUser := &[]User{}
	fileReader := excel.OpenFile("./example.xlsx")

	err := fileReader.ReadSheets(
		fileReader.ReadSheetNo(0, destUser).TitleRow(0),
		fileReader.ReadSheetNo(1, dest).TitleRow(1),
	).Read()
	if err != nil {
		println(err.Error())
	}
	r, err := json.Marshal(dest)
	if err != nil {
		println(err.Error())
	}
	println(string(r))
}

func TestReflect(t *testing.T) {
	var s1 []int
	s2 := make([]int, 0)
	s4 := make([]int, 0)
	fmt.Printf("s1 pointer:%+v, s2 pointer:%+v, s4 pointer:%+v, \n", *(*reflect.SliceHeader)(unsafe.Pointer(&s1)), *(*reflect.SliceHeader)(unsafe.Pointer(&s2)), *(*reflect.SliceHeader)(unsafe.Pointer(&s4)))
	fmt.Printf("%v\n", (*(*reflect.SliceHeader)(unsafe.Pointer(&s1))).Data == (*(*reflect.SliceHeader)(unsafe.Pointer(&s2))).Data)
	fmt.Printf("%v\n", (*(*reflect.SliceHeader)(unsafe.Pointer(&s2))).Data == (*(*reflect.SliceHeader)(unsafe.Pointer(&s4))).Data)
	s1 = append(s1, 1)
	fmt.Printf("s1 pointer:%+v, s2 pointer:%+v, s4 pointer:%+v, \n", *(*reflect.SliceHeader)(unsafe.Pointer(&s1)), *(*reflect.SliceHeader)(unsafe.Pointer(&s2)), *(*reflect.SliceHeader)(unsafe.Pointer(&s4)))

	//for i := 0; i < 10; i++ {
	//	defer println(strconv.Itoa(i))
	//}
}

type ExcelExample struct {
	Standard Standard
	User     *User
	Orders   *[]Order
}

type Standard struct {
	Message string `json:"message" excel:"name:消息"`
}

type User struct {
	Id      int       `json:"id" excel:"index:1 ; name:用户id; "`
	Name    string    `json:"name" excel:"name:用户名; "`
	Time    time.Time `json:"time" excel:"name:创建时间; "`
	Roles   []string  `json:"time" excel:"name:角色; converter:ConvertRoles; "`
	Message string    `json:"message" excel:"index:5; name:消息; "`
}

func (u *User) ConvertRoles(cellVal string) error {
	if cellVal == "" {
		return nil
	}
	u.Roles = strings.Split(cellVal, "|")
	return nil
}

type Order struct {
	OrderId   string    `json:"orderId" excel:" index:6; name:订单ID;" `
	OrderTime time.Time `json:"orderTime" excel:" index:7; name:订单时间; "`
	Message   string    `json:"message" excel:"index:8; name:消息; "`
}
