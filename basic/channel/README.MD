#go Channel 管道
管道是 Golang 在语言级别上提供的 goroutine 间的通讯方式，我们可以使用 channel 在 多个 goroutine 之间传递消息
channel 是可以让一个 goroutine 发送特定值到另一个 goroutine 的通信机制

Golang 的并发模型是 CSP（Communicating Sequential Processes），提倡通过通信共享内 存而不是通过共享内存而实现通信。

Go 语言中的管道（channel）是一种特殊的类型。管道像一个传送带或者队列，
总是遵 循先入先出（First In First Out）的规则，保证收发数据的顺序。每一个管道都是一个具体类型的导管

##channel 类型
var 变量 chan 元素类型
var ch1 chan int // 声明一个传递整型的管道 
var ch3 chan []int // 声明一个传递 int

##创建channel
声明的管道后需要使用 make 函数初始化之后才能使用。
make(chan 元素类型, 容量)

//创建一个能存储 10 个 int 类型数据的管道 
ch1 := make(chan int, 10) 
//创建一个能存储 3 个[]int 切片类型数据的管道 
ch3 :=  make(chan []int, 3)

##管道发送
ch <- 10 // 把 10 发送到 ch 中

##管道接收
x := <- ch // 从 ch 中接收值并赋值给变量 
x <-ch // 从 ch 中接收值，忽略结果

##管道关闭
close(ch)

只有在通知接收方 goroutine 所有的数据都发送完毕的时 候才需要关闭管道。
管道是可以被垃圾回收机制回收的，它和关闭文件是不一样的，在结束 操作之后关闭文件是必须要做的，但关闭管道不是必须的。

关闭后的管道有以下特点： 
1. 对一个关闭的管道再发送值就会导致 panic。 
2. 对一个关闭的管道进行接收会一直获取值直到管道为空。
3. 对一个关闭的并且没有值的管道执行接收操作会得到对应类型的零值。 
4. 关闭一个已经关闭的管道会导致 panic。


for range 遍历管道，当管道被关闭的时候就会退出 for range
for val := range ch1 { fmt.Println(val) }

##管道阻塞
如果创建管道的时候没有指定容量，那么我们可以叫这个管道为无缓冲的管道 无缓冲的管道又称为阻塞的管道。
ch := make(chan int)
ch <- 10  //能够通过编译，但是执行的时候会出错




