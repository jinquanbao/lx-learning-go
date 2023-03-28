##go 切片(slice)
切片（Slice）是一个拥有相同类型元素的可变长度的序列。它是基于数组类型做的一层封装。 它非常灵活，支持自动扩容。

##切片是引用类型
切片是一个引用类型，它的内部结构包含地址、长度和容量。

##切片声明
var name []T
1. name:表示变量名 
2. T:表示切片中的元素类型

##切片的长度和容量
切片拥有自己的长度和容量，我们可以通过使用内置的 len()函数求长度，使用内置的 cap() 函数求切片的容量。

切片的长度就是它所包含的元素个数。
切片的容量是从它的第一个元素开始数，到其底层数组元素末尾的个数。

##切片定义
###1.基于数组定义

//定义数组
var a := [2]int{1,2}
//定义切片
var b := a

###2.切片再切片
func main() {
     //切片再切片 
     a := [...]string{"北京", "上海", "广州", "深圳", "成都", "重庆"} 
     fmt.Printf("a:%v type:%T len:%d cap:%d\n", a, a, len(a), cap(a)) 
     b := a[1:3] 
     fmt.Printf("b:%v type:%T len:%d cap:%d\n", b, b, len(b), cap(b)) 
     c := b[1:5] 
     fmt.Printf("c:%v type:%T len:%d cap:%d\n", c, c, len(c), cap(c)) 
 }
 输出： 
 a:[北京 上海 广州 深圳 成都 重庆] type:[6]string len:6 cap:6 
 b:[上海 广州] type:[]string len:2 cap:5 
 c:[广州 深圳 成都 重庆] type:[]string len:4 cap:4

###3.使用 make()函数构造切片
如果需要动态的创建一个切片，我们就需要使用内置 的 make()函数
make([]T, size, cap)
其中： 3. T:切片的元素类型 4. size:切片中元素的数量 5. cap:切片的容量
例：
a := make([]int, 2, 10) 
fmt.Println(a) //[0 0] 
fmt.Println(len(a)) //2 
fmt.Println(cap(a)) //10
