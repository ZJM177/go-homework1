package homework01

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func Testmain(m *testing.M) {
	code := m.Run()
	num := 10
	addTen(&num)
	fmt.Println(num)
	slice1 := []int{1, 2, 3, 4, 5}
	doubleSlice(&slice1)
	fmt.Println(slice1)
	go_print_odd_even()

	r := Rectangle{10, 20}
	c := Circle{5}
	fmt.Println("Rectangle Area:", r.Area())
	fmt.Println("Rectangle Perimeter:", r.Perimeter())
	fmt.Println("Circle Area:", c.Area())
	fmt.Println("Circle Perimeter:", c.Perimeter())

	employee1 := Employee{Person{"Alice", 30}, "E12345"}
	employee1.PrintInfo()

	channel_communication()
	buffered_channel_communication()
	var wg sync.WaitGroup
	safeCounter := newSafeCounter()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			safeCounter.Incr()

		}()
	}
	wg.Wait()
	fmt.Println("SafeCounter Value:", safeCounter.getCount())

	atomic_counter()

	// 创建调度器
	scheduler := NewTaskScheduler()

	// 添加示例任务
	scheduler.AddTask("任务1", func() error {
		time.Sleep(1 * time.Second)
		fmt.Println("任务1执行完成")
		return nil
	})

	scheduler.AddTask("任务2", func() error {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("任务2执行完成")
		return nil
	})

	scheduler.AddTask("任务3", func() error {
		time.Sleep(2 * time.Second)
		fmt.Println("任务3执行完成")
		return fmt.Errorf("模拟错误")
	})

	scheduler.AddTask("任务4", func() error {
		time.Sleep(300 * time.Millisecond)
		fmt.Println("任务4执行完成")
		return nil
	})

	// 执行所有任务
	results := scheduler.Run()

	// 打印统计报告
	scheduler.PrintReport()

	// 你也可以单独处理结果
	fmt.Printf("\n=== 单独处理结果 ===\n")
	for _, r := range results {
		fmt.Printf("任务 %s 执行了 %v\n", r.Name, r.Duration)
	}

	os.Exit(code)
}

//：编写一个Go程序，
// 定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，
// 然后在主函数中调用该函数并输出修改后的值。

func addTen(num *int) {
	*num += 10
}

//实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。

func doubleSlice(s *[]int) {
	for i := range *s {
		(*s)[i] *= 2
	}
}

//编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。

func go_print_odd_even() {
	go func() {
		for i := 1; i <= 10; i += 2 {
			fmt.Println("奇数", i)
		}
	}()
	go func() {
		for i := 2; i <= 10; i += 2 {
			fmt.Println("偶数", i)
		}
	}()
	time.Sleep(1 * time.Second)

}

//设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。

// 任务调度器的实现可以参考 https://github.com/panjf2000/ants/blob/master/ants.go

// 任务类型定义
type Task struct {
	Name     string
	Function func() error
}

// 任务执行结果
type TaskResult struct {
	Name      string
	Duration  time.Duration
	Err       error
	StartTime time.Time
	EndTime   time.Time
}

// 任务调度器
type TaskScheduler struct {
	tasks     []Task
	results   []TaskResult
	wg        sync.WaitGroup
	mu        sync.Mutex
	startTime time.Time
	endTime   time.Time
}

// 创建新的调度器
func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{
		tasks:   make([]Task, 0),
		results: make([]TaskResult, 0),
	}
}

// 添加任务
func (s *TaskScheduler) AddTask(name string, fn func() error) {
	s.tasks = append(s.tasks, Task{Name: name, Function: fn})
}

// 执行单个任务
func (s *TaskScheduler) runTask(task Task) {
	defer s.wg.Done()

	startTime := time.Now()

	// 执行任务函数
	err := task.Function()

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	// 保存结果
	s.mu.Lock()
	s.results = append(s.results, TaskResult{
		Name:      task.Name,
		Duration:  duration,
		Err:       err,
		StartTime: startTime,
		EndTime:   endTime,
	})
	s.mu.Unlock()
}

// 并发执行所有任务
func (s *TaskScheduler) Run() []TaskResult {
	fmt.Printf("开始执行 %d 个任务...\n", len(s.tasks))

	s.startTime = time.Now()

	// 为每个任务启动一个goroutine
	for _, task := range s.tasks {
		s.wg.Add(1)
		go s.runTask(task)
	}

	// 等待所有任务完成
	s.wg.Wait()

	s.endTime = time.Now()

	return s.results
}

// 打印统计报告
func (s *TaskScheduler) PrintReport() {
	fmt.Printf("\n=== 任务执行统计报告 ===\n")
	fmt.Printf("总任务数: %d\n", len(s.tasks))
	fmt.Printf("总执行时间: %v\n", s.endTime.Sub(s.startTime))

	fmt.Printf("\n详细结果:\n")
	for _, result := range s.results {
		status := "✓ 成功"
		if result.Err != nil {
			status = fmt.Sprintf("✗ 失败: %v", result.Err)
		}
		fmt.Printf("任务: %-20s 耗时: %-12v 状态: %s\n",
			result.Name, result.Duration, status)
	}

	// 计算平均时间
	if len(s.results) > 0 {
		var totalDuration time.Duration
		for _, r := range s.results {
			totalDuration += r.Duration
		}
		avgDuration := totalDuration / time.Duration(len(s.results))
		fmt.Printf("\n平均执行时间: %v\n", avgDuration)
	}
}

//题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
//考察点 ：接口的定义与实现、面向对象编程风格。

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	width  float64
	height float64
}

type Circle struct {
	radius float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

func (c Circle) Area() float64 {
	return 3.14 * c.radius * c.radius
}

func (c Circle) Perimeter() float64 {
	return 2 * 3.14 * c.radius
}

// 题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
// 考察点 ：组合的使用、方法接收者。
type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeID string
}

func (e Employee) PrintInfo() {
	fmt.Printf("Name: %s, Age: %d, EmployeeID: %s\n", e.Name, e.Age, e.EmployeeID)
}

// 题目 ：编写一个程序，使用通道实现两个协程之间的通信。
// 一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
func channel_communication() {
	ch := make(chan int)

	go func() {
		for i := 1; i <= 10; i++ {
			ch <- i
		}
		close(ch)
	}()
	go func() {
		for num := range ch {
			fmt.Println("通道：", num)
		}
	}()
	time.Sleep(1 * time.Second)
}

// 题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
func buffered_channel_communication() {
	ch := make(chan int, 100)
	go func() {
		for i := 0; i < 100; i++ {
			ch <- i
		}
		close(ch)
	}()
	go func() {
		for num := range ch {
			fmt.Println("缓冲通道：", num)
		}
	}()
	time.Sleep(2 * time.Second)
}

//	题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，
//
// 每个协程对计数器进行1000次递增操作，最后输出计数器的值。
// 考察点 ： sync.Mutex 的使用、并发数据安全。
type SafeCounter struct {
	count int
	mutex sync.Mutex
}

func newSafeCounter() *SafeCounter {
	return &SafeCounter{0, sync.Mutex{}}
}
func (sc *SafeCounter) Incr() {
	defer sc.mutex.Unlock()
	sc.mutex.Lock()
	sc.count++
}
func (sc *SafeCounter) getCount() int {
	defer sc.mutex.Unlock()
	sc.mutex.Lock()
	return sc.count
}

// 题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
// 考察点 ：原子操作、并发数据安全。
type AtomicCounter struct {
	count int64
}

func newAtomicCounter() *AtomicCounter {
	return &AtomicCounter{0}
}
func (ac *AtomicCounter) Incr() {
	atomic.AddInt64(&ac.count, 1)
}
func (ac *AtomicCounter) getCount() int64 {
	return atomic.LoadInt64(&ac.count)
}
func (ac *AtomicCounter) add(n int64) {
	atomic.AddInt64(&ac.count, n)
}
func atomic_counter() {
	counter := newAtomicCounter()
	var wg sync.WaitGroup

	startTime := time.Now()

	// 启动10个goroutine
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			counter.Incr()
			fmt.Printf("goroutine %d 完成\n", id)
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(startTime)

	fmt.Printf("\n最终计数: %d\n", counter.getCount())
	fmt.Printf("执行时间: %v\n", elapsed)
}
