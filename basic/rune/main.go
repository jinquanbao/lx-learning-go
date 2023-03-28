package main

import "fmt"

func main() {
	var str = "第一行,xxx"
	for i := 0; i < len(str); i++ { //byte
		fmt.Printf("i=%v(%c)", str[i], str[i])
	}
	fmt.Println()
	for _, i := range str { //rune
		fmt.Printf("i=%v(%c)", i, i)
	}

	//要修改字符串，需要先将其转换成[]rune 或[]byte，完成后再转换为 string。
	//无论哪种转换， 都会重新分配内存，并复制字节数组。

	fmt.Println()

	str2 := "123"

	b1 := []byte(str2)
	b1[0] = 'a'

	fmt.Println(string(b1))

	str3 := "速度"
	b2 := []rune(str3)
	b2[0] = '啊'
	fmt.Println(string(b2))

}
