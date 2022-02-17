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

