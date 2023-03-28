package main

import "fmt"

func main() {
	var i int = 10
	var j int = 10

	fmt.Println("i=%v", i)
	fmt.Println("i=%v", i)

	fmt.Printf("i=%v", i)
	fmt.Println()
	fmt.Printf("i=%X,j=%d", i, j)
}
