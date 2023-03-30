package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

//在go中使用 sync.NewCond(l Locker)可以创建一个与锁l相关的条件变量cond，
//这个锁l可以是Mutex 或者RWMutex类型的锁，该方法会返回Cond类型的指针被称为条件变量。
/***
//Cond类型主要有三个方法：
(c *Cond) Wait()：等待方法，当一个goroutine调用cond的Wait()方法后，
				  当前goroutine会被阻塞，直到有其他goroutine调用了该cond变量的Signal()或者Broadcast()方法；
				  需要注意的是调用Wait()方法前，当前goroutine要通过cond.l.Lock() 先获取与cond关联的锁l，
				  否则会抛出异常 fatal error:sync:unlock of unlocked mutex

(c *Cond) Signal()，通知方法，当一个goroutine调用了条件变量c的该方法后，就会激活一个由于调用c的Wait()方法而被阻塞的goroutine。

(c *Cond) Broadcast() 通知方法，当一个goroutine调用了条件变量c的该方法后，就会激活所有由于调用c的Wait()方法而被阻塞的goroutine。
*/
func main() {
	var m sync.Mutex
	c := sync.NewCond(&m)

	ready := make(chan struct{}, 10)
	isReady := false

	for i := 0; i < 10; i++ {
		i := i
		go func() {
			m.Lock()
			time.Sleep(time.Duration(rand.Int63n(2)) * time.Second)
			ready <- struct{}{} // 运动员i准备就绪
			//Wait 调用之前必须要加锁
			for !isReady {
				c.Wait()
			}
			log.Printf("%d started\n", i)
			m.Unlock()
		}()
	}
	// false broadcast
	c.Broadcast()

	// 裁判员检查所有的运动员是否就绪
	for i := 0; i < 10; i++ {
		<-ready
	}
	isReady = true
	// 运动员都已准备就绪，发令枪响, broadcast
	c.Broadcast()

	time.Sleep(time.Second)
}
