package main

import "fmt"

/***
 Go 语言在声明变量的时候，会自动对变量对应的内存区域进行初始化操作。
每个变量会被 初始化成其类型的默认值，例如： 整型和浮点型变量的默认值为 0。
字符串变量的默认值 为空字符串。
布尔型变量默认为 false。
切片、函数、指针变量的默认为 nil。
*/

func main() {
	//var i  = 10
	//var j int = 8
	//一次定义多个变量
	var i, j int
	i = 10
	j = 8
	//一次定义多个变量 并赋值
	var x, y int = 10, 20

	var user = "a"

	fmt.Printf("user=%v", user)
	fmt.Println()
	fmt.Printf("i=%d,j=%x", i, j)
	fmt.Println()
	fmt.Printf("x=%v,y=%v", x, y)
	fmt.Println()

	//短变量声明 ：短变量只能用于声明局部变量，不能用于全局变量的声明

	n := 10
	fmt.Println("n=", n)

	//批量定义
	var (
		a string
		b int
		c bool
	)
	a = "a"
	b = 1
	c = true
	fmt.Printf("a=%v,b=%v,c=%v", a, b, c)

}
