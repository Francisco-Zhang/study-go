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

非阻塞式的获取数据方法。

底层是使用的IO多路复用技术，可以按顺序的检测channel是否有数据进入，哪个有数据，就执行哪个channel的数据读写操作。

如果都没有数据，就会出现死锁错误，所以必须有default语句在检测不到数据的情况下执行。

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

想要一致监听，就在外面加一个for循环

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





6min