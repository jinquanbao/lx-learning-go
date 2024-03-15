package excelutil

import (
	"strings"
	"testing"
	"time"
)

func TestReadSingleSheet(t *testing.T) {
	dest := &[]User{}
	err := ReadFile("./example.xlsx", dest)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dest)
}

func TestReadMultipleSheet(t *testing.T) {
	dest := &[]ExcelExample{}
	destUser := &[]User{}
	fileReader := OpenFile("./example.xlsx")

	err := fileReader.ReadSheets(
		fileReader.ReadSheetNo(0, destUser).TitleRow(0).
			RegisterReadCellCompleteCallbacks(func(rCtx ReadCellContext, destElem interface{}, err error) error {
				t.Logf("ReadCellComplete : %v", destElem)
				return nil
			}).
			RegisterReadSheetCompleteCallbacks(func(rCtx ReadSheetContext, dest interface{}) error {
				t.Logf("ReadSheetComplete: %v  ", dest)
				return nil
			}),
		fileReader.ReadSheetNo(1, dest).TitleRow(1),
	).Read()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*dest)
	t.Log(destUser)
}

func TestReadToMap(t *testing.T) {
	var dest []map[string]interface{}

	reader := OpenFile("./example.xlsx")

	err := reader.ReadSheets(reader.
		ReadSheetNo(0, &dest).
		RegisterReadCellCompleteCallbacks(func(rCtx ReadCellContext, destElem interface{}, err error) error {
			if rCtx.RowIndex() == 0 {
				t.Log(rCtx.TitleName())
			}
			return nil
		}).
		TimeTitles("创建时间")).
		Read()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dest)
}

func TestSmartOpenFile(t *testing.T) {
	reader, err := SmartOpenFile("./example.xlsx")
	defer func() {
		if err := reader.Close(); err != nil {
			t.Logf("%v", err)
		}
	}()
	if err != nil {
		t.Fatal(err)
	}
	sheetList := reader.File.GetSheetList()
	t.Logf("%v", sheetList)
	var dest []map[string]interface{}
	for i := range sheetList {
		err := reader.ReadSheetNo(i, &dest).RegisterReadSheetCompleteCallbacks(func(rCtx ReadSheetContext, dest interface{}) error {
			t.Logf("%v", dest)
			return nil
		}).Read()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestReadDisableAutoClose(t *testing.T) {
	dest := &[]User{}
	reader := OpenFile("./example.xlsx").DisableAutoClose()
	// 禁用读取自动关闭需手动关闭
	defer func() {
		if err := reader.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	err := reader.ReadSheetNo(0, dest).
		Read()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dest)
	// 可以拿到 excelize.File 做后续的操作
	t.Log(reader.File.SheetCount)
}

//func TestReflect(t *testing.T) {
//	var s1 []int
//	s2 := make([]int, 0)
//	s4 := make([]int, 0)
//	fmt.Printf("s1 pointer:%+v, s2 pointer:%+v, s4 pointer:%+v, \n", *(*reflect.SliceHeader)(unsafe.Pointer(&s1)), *(*reflect.SliceHeader)(unsafe.Pointer(&s2)), *(*reflect.SliceHeader)(unsafe.Pointer(&s4)))
//	fmt.Printf("%v\n", (*(*reflect.SliceHeader)(unsafe.Pointer(&s1))).Data == (*(*reflect.SliceHeader)(unsafe.Pointer(&s2))).Data)
//	fmt.Printf("%v\n", (*(*reflect.SliceHeader)(unsafe.Pointer(&s2))).Data == (*(*reflect.SliceHeader)(unsafe.Pointer(&s4))).Data)
//	s1 = append(s1, 1)
//	fmt.Printf("s1 pointer:%+v, s2 pointer:%+v, s4 pointer:%+v, \n", *(*reflect.SliceHeader)(unsafe.Pointer(&s1)), *(*reflect.SliceHeader)(unsafe.Pointer(&s2)), *(*reflect.SliceHeader)(unsafe.Pointer(&s4)))
//
//	//for i := 0; i < 10; i++ {
//	//	defer println(strconv.Itoa(i))
//	//}
//}

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
	Roles   []string  `json:"roles" excel:"name:角色; converter:ConvertRoles; "`
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
