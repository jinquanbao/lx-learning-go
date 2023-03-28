package main

//定义一个长度为3的数组
var a [3]int

//定义一个长度为3的数组并赋值
var b = [3]int{1, 2, 3}

//定义一个长度为2的数组并赋值，让编译器 根据初始值的个数自行推断数组的长度
var c = [...]int{1, 2}

//定义一个长度为6的数组，并给第2个和第6个赋值，让编译器 根据初始值的个数自行推断数组的长度
var d = [...]int{1: 2, 5: 3}

func main() {
	for i, value := range d {
		println("i=", i, "value=", value)
	}
}
