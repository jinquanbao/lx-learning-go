package main

import "fmt"

const i = 10

//批量声明
const (
	a = "a"
	b = 1
)

//const 常量结合 iota
//iota 是 golang 语言的常量计数器,只能在常量的表达式中使用
//iota 在 const 关键字出现时将被重置为 0(const 内部的第一行之前)，const 中每新增一行常量 声明将使 iota 计数一次(iota 可理解为 const 语句块中的行索引)
const name = 1
const (
	x = iota //0
	y        //1
	z        //2
)

func main() {
	fmt.Println("i=", i, ",a=", a, ",b=", b)
	fmt.Println("x=", x, ",y=", y, ",z=", z)
}
