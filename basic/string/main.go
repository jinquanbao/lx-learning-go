package main

import (
	"fmt"
	"strings"
)

func main() {
	var str = "第一行,xxx"
	var str2 = "第二行"
	fmt.Println("str内容=", str)
	fmt.Println("str长度=", len(str))
	fmt.Println("str分割数组=", strings.Split(str, ","))
	fmt.Println("str是否包含XXX：", strings.Contains(str, "xxx"))
	fmt.Println("str中XXX出现的位置：", strings.Index(str, "xxx"))
	fmt.Println("合并字符串数组：", strings.Join([]string{"A", "B"}, ","))
	fmt.Println("字符串相加：", str+str2)

}
