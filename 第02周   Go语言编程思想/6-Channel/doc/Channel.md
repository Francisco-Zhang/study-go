## 1、 channe

```go
func chanDemo() {
	c := make(chan int)
	c <- 1  //执行完这句之后，会等待channe内数据被读取才返回
	c <- 2
	n := <-c
	fmt.Println(n)
}
func main() {
	chanDemo()
}
//执行会陷入死锁，all gorutines are asleep
```

```go
func chanDemo1() {
	c := make(chan int)
	go func() {
		for {
			n := <-c
			fmt.Println(n)
		}
	}()
	c <- 1
	c <- 2
}
// 如果匿名函数协程执行慢的化，有可能只打印一个1，然后主协程就结束了。
```



channel 也是一等公民，既可以当参数，也可以当返回值

```go
func createWorker(id int) chan<- int {   //此处声明返回的 channel 是用来收数据的 channel
	c := make(chan int)
	go worker(id, c)
	return c
}

func chanDemo() {
	var channels [10]chan<- int
	for i := 0; i < 10; i++ {
		channels[i] = createWorker(i)
	}

	for i := 0; i < 10; i++ {
		channels[i] <- 'a' + i
	}

	for i := 0; i < 10; i++ {
		channels[i] <- 'A' + i
	}

	time.Sleep(time.Millisecond)
}
```

```go
func bufferedChannel() {
	c := make(chan int, 3)  //缓存区大小为3，这样就不用每次都等待，发送数据的协程在发送第四个数据的时候才会等待。
	go worker(0, c)
	c <- 'a'
	c <- 'b'
	c <- 'c'
	c <- 'd'
	time.Sleep(time.Millisecond)
}
```

```go
func channelClose() {
	c := make(chan int)
	go worker(0, c)
	c <- 'a'
	c <- 'b'
	c <- 'c'
	c <- 'd'
	close(c) //close之后，外面的协程收到的是 int 类型的默认值0，外面的协程根据这个判断channel是否被关闭
	time.Sleep(time.Millisecond)
}
```

```go
func worker(id int, c chan int) {
	//_,ok:= <-c
	//if !ok{  判断是否被关闭，也可以用下面的range
	//}
	for n := range c {
		fmt.Printf("Worker %d received %c\n",
			id, n)
	}
}
```

**不要通过共享内存来通讯，通过通讯来共享内存。**  协程发送数据后，等待的是通讯。



## 2、 使用Channel等待任务结束

### 方案一

方案一是 createWorker 增加一个通知任务结束的channel，但这种方式没法提高并行，最终是顺序执行的。所以改用第二种方案。

```go
func doworker(id int, c chan int, done chan bool) {
	for n := range c {
		fmt.Printf("Worker %d received %c\n",
			id, n)
		done <- true
	}
}

type worker struct {
	in   chan int
	done chan bool
}

func createWorker(id int) worker {
	w := worker{
		in:   make(chan int),
		done: make(chan bool),
	}

	go doworker(id, w.in, w.done)
	return w
}

func chanDemo() {
	var workers [10]worker
	for i := 0; i < 10; i++ {
		workers[i] = createWorker(i)
	}

	for i := 0; i < 10; i++ {
		workers[i].in <- 'a' + i
		<-workers[i].done
	}

	for i := 0; i < 10; i++ {
		workers[i].in <- 'A' + i
		<-workers[i].done
	}
}

func main() {
	chanDemo()
}
```

### 方案一改进：

```go
func doworker(id int, c chan int, done chan bool) {
	for n := range c {
		fmt.Printf("Worker %d received %c\n",
			id, n)
		//再开一个协程，防止主协程第二个for 循环阻塞
		go func() {
			done <- true
		}()

	}
}

type worker struct {
	in   chan int
	done chan bool
}

func createWorker(id int) worker {
	w := worker{
		in:   make(chan int),
		done: make(chan bool),
	}
	go doworker(id, w.in, w.done)
	return w
}

func chanDemo() {
	var workers [10]worker
	for i := 0; i < 10; i++ {
		workers[i] = createWorker(i)
	}

	for i, worker := range workers {
		worker.in <- 'a' + i
	}

	for i, worker := range workers {
		worker.in <- 'A' + i
	}

	for _, worker := range workers {
		<-worker.done
		<-worker.done
	}

}

func main() {
	chanDemo()
}
```

### 方案二：WaitGroup

```go
func doWork(id int,
	w worker) {
	for n := range w.in {
		fmt.Printf("Worker %d received %c\n",
			id, n)
		w.done()
	}
}

type worker struct {
	in   chan int
	done func()
}

func createWorker(
	id int, wg *sync.WaitGroup) worker {
	w := worker{
		in: make(chan int),
		done: func() {
			wg.Done()
		},
	}
	go doWork(id, w)
	return w
}

func chanDemo() {
	var wg sync.WaitGroup

	var workers [10]worker
	for i := 0; i < 10; i++ {
		workers[i] = createWorker(i, &wg)
	}

	wg.Add(20)
	for i, worker := range workers {
		worker.in <- 'a' + i
	}
	for i, worker := range workers {
		worker.in <- 'A' + i
	}

	wg.Wait()
}
func main() {
	chanDemo()
}
```



## 3、 使用Channel进行树的遍历

```go
c := root.TraverseWithChannel()
	maxNodeValue := 0
	for node := range c {
		if node.Value > maxNodeValue {
			maxNodeValue = node.Value
		}
	}


func (node *Node) TraverseWithChannel() chan *Node {
	out := make(chan *Node)
	go func() {
		node.TraverseFunc(func(node *Node) {
			out <- node
		})
		close(out)
	}()
	return out
}
```



## 4、 Select 

一、select简介

        1、Go的select语句是一种仅能用于channl发送和接收消息的专用语句，此语句运行期间是阻塞的；当select中没有case语句的时候，会阻塞当前groutine。
        
        2、select是Golang在语言层面提供的I/O多路复用的机制，其专门用来检测多个channel是否准备完毕：可读或可写。
    
        3、select语句中除default外，每个case操作一个channel，要么读要么写
    
        4、select语句中除default外，各case执行顺序是随机的
    
        5、select语句中如果没有default语句，则会阻塞等待任一case
    
        6、select语句中读操作要判断是否成功读取，关闭的channel也可以读取




Select：多个channel的复用

非阻塞式的获取数据方法。

底层是使用的IO多路复用技术，可以按顺序的检测channel是否有数据进入，哪个有数据，就执行哪个channel的数据读写操作。

如果都没有数据，就会出现死锁错误，所以必须有default语句在检测不到数据的情况下执行。

原理：

select底层会转换成for循环顺序执行，如果判断出channel没有数据就会进行下一次循环，不会阻塞。

![1](img/1.png)



```go
func main() {
	var c1, c2 chan int
	//想同时判断 c1,c2而不会被阻塞
	select {
	case n := <-c1:
		fmt.Println("Received from c1", n)
	case n := <-c2:
		fmt.Println("Received from c2", n)
	default:
		fmt.Println("No Value Received ")
	}
}
```

想要一直监听，就在外面加一个for循环

```go
func generator() chan int {
	out := make(chan int)
	go func() {
		i := 0
		for {
			time.Sleep(
				time.Duration(rand.Intn(1500)) *
					time.Millisecond)
			out <- i
			i++
		}
	}()
	return out
}

func main() {
	var c1, c2 chan int = generator(), generator()

	for {
		select {
		case n := <-c1:
			fmt.Println("Received from c1", n)
		case n := <-c2:
			fmt.Println("Received from c2", n)
		}
	}
}
```



```go
func generator() chan int {
	out := make(chan int)
	go func() {
		i := 0
		for {
			time.Sleep(
				time.Duration(rand.Intn(1500)) *
					time.Millisecond)
			out <- i
			i++
		}
	}()
	return out
}

func worker(id int, c chan int) {
	for n := range c {
		time.Sleep(time.Second)
		fmt.Printf("Worker %d received %d\n",
			id, n)
	}
}

func createWorker(id int) chan<- int {
	c := make(chan int)
	go worker(id, c)
	return c
}

func main() {
	var c1, c2 chan int = generator(), generator()
	worker := createWorker(0)
	n := 0
	hasValue := false
	for {
		var activeWorker chan<- int  // nil 也可以被select 阻塞
		if hasValue == true {
			activeWorker = worker
		}
		select {
		case n = <-c1:
			hasValue = true
		case n = <-c2:
			hasValue = true
		case activeWorker <- n:   //消耗数据太慢，n会被冲掉
			hasValue = false
		}
	}
}

```



```go
//解决n被冲掉
var values []int
for {
    var activeWorker chan<- int
    var activeValue int
    if len(values) > 0 {
        activeWorker = worker
        activeValue = values[0]
    }

    select {
        case n := <-c1:
        values = append(values, n)
        case n := <-c2:
        values = append(values, n)
        case activeWorker <- activeValue:
        values = values[1:]
    }
}
```



```go
tm := time.After(10 * time.Second)  //10s后通过return结束程序
tick := time.Tick(time.Second)    //每隔1s送出一份数据
for {
    var activeWorker chan<- int
    var activeValue int
    if len(values) > 0 {
        activeWorker = worker
        activeValue = values[0]
    }

    select {
        case n := <-c1:
        	values = append(values, n)
        case n := <-c2:
        	values = append(values, n)
        case activeWorker <- activeValue:
        	values = values[1:]
        case <-time.After(800 * time.Millisecond): //800ms的时间内没有生出数据
			fmt.Println("timeout")
        case <-tick:
			fmt.Println("queue len =", len(values))
        case <-tm:
        	fmt.Println("bye")
        return
    }
}
```



## 5、 传统同步机制

传统的同步模型要少用，Go语言习惯用CSA通讯模型，通过channel达到同步，而不是用锁。

- WaitGroup
- Mutex
- Cond



如果变量有读写冲突等非安全的使用， 使用go -race 命令可以检测出来，会有提示信息。



```go
type atomicInt struct {
	value int
	lock  sync.Mutex
}

func (a *atomicInt) increment() {
	fmt.Println("safe increment")
	func() {  //通过匿名函数来对一块代码区域进行保护。
		a.lock.Lock()
		defer a.lock.Unlock()

		a.value++
	}()
}

func (a *atomicInt) get() int {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.value
}

func main() {
	var a atomicInt
	a.increment()
	go func() {
		a.increment()
	}()
	time.Sleep(time.Millisecond)
	fmt.Println(a.get())
}
```



## 6、 并发模式（上）

## 7、 并发模式（下）





```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func msgGen(name string) chan string {
	c := make(chan string)
	go func() {
		i := 0
		for {
			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
			c <- fmt.Sprintf("service %s: message %d", name, i)
			i++
		}
	}()
	return c
}

func fanIn(chs ...chan string) chan string {
	c := make(chan string)
	for _, ch := range chs {
		go func(in chan string) {  //要防止ch在循环中被覆盖，所以需要使用参数传入，而不能直接使用。
			for {
				c <- <-in
			}
		}(ch)
	}
	return c
}

//fanIn的另一种写法，这种方式可以减少goroutine使用，但是必须明确的知道 channel 的个数
func fanInBySelect(c1, c2 chan string) chan string {
	c := make(chan string)
	go func() {
		for {
			select {
			case m := <-c1:
				c <- m
			case m := <-c2:
				c <- m
			}
		}
	}()
	return c
}

func main() {
	m1 := msgGen("service1")
	m2 := msgGen("service2")
	m3 := msgGen("service3")
    //将多个channel集中到一个，可以防止从多个channel取数据。
    //使用场景是多个任务，哪个任务先执行完就处理哪个任务的返回结果
	m := fanIn(m1, m2, m3) 
	for {
		fmt.Println(<-m)
	}
}
```



## 8、 并发任务的控制

### 处理机制

- 非阻塞等待
- 超时机制
- 任务中断/退出
- 优雅退出



### 非阻塞等待



```go
func nonBlockingWait(c chan string) (string, bool) {
	select {
	case m := <-c:
		return m, true
	default:  //select 中一旦引入了default就变成了非阻塞，会直接执行default。
		return "", false
	}
}
```



### 超时机制



```go
func timeoutWait(c chan string, timeout time.Duration) (string, bool) {
	select {
	case m := <-c:
		return m, true
	case <-time.After(timeout):
		return "", false
	}
}
```



### 任务中断/退出

```go
func msgGen(name string,done chan struct{}) chan string {  //空的struct更省空间
	c := make(chan string)
	go func() {
		i := 0
		for {
            select{
                case <-time.After(time.Duration(rand.Intn(2000)) * time.Millisecond):
               		 c <- fmt.Sprintf("service %s: message %d", name, i)
                case <-done:
                	 fmt.Println("clean up")
                	 return
            }
			i++
		}
	}()
	return c
}

func main() {
    done:=make(chan struct{})
	m1 := msgGen("service1",done)
    done<-struct{}{}	//第一个{}定义一个没有字段的struct,第二个{}初始化为空
    time.Sleep(time.Second)
}

```



### 优雅退出

```go
func msgGen(name string,done chan struct{}) chan string {  //空的struct更省空间
	c := make(chan string)
	go func() {
		i := 0
		for {
            select{
                case <-time.After(time.Duration(rand.Intn(2000)) * time.Millisecond):
               		 c <- fmt.Sprintf("service %s: message %d", name, i)
                case <-done:
                	 fmt.Println("clean up")
                	 time.Sleep(time.Second)
                	 fmt.Println("clean done")
                	 done<-struct{}{}   //双向channel不便于理解，一般不用，这里暂时用下，处理结束后再返回一个数据
                	 return
            }
			i++
		}
	}()
	return c
}


func main() {
    done:=make(chan struct{})
	m1 := msgGen("service1",done)
    done<-struct{}{}	//第一个{}定义一个没有字段的struct,第二个{}初始化为空
    <-done
}
```

