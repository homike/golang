package channeltest

import (
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"sync"
)

//---------------------------------------------------------
func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func count(nums <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range nums {
			out <- n * n
		}
		close(out)
	}()
	return out
}

// 在函数中创建一个channel 返回，同时创建一个gorutine 往 channel 中塞数据， 这是一个重要的惯用法
func RunChan1() {
	out := gen(1, 2, 3)

	for v := range out {
		fmt.Println(v)
	}
}

func RunChan11() {
	//out1 := gen(1, 2, 3)
	//out2 := count(out1)
	// for v := range out2 {
	// 	fmt.Println(v)
	// }

	// Set up the pipeline and consume the output.
	for n := range count(count(gen(2, 3))) {
		fmt.Println(n) // 16 then 81
	}
}

//---------------------------------------------------------

//Fan-out, Fan-in
//---------------------------------------------------------
func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func RunChan2() {
	trace.Start(os.Stderr)

	in := gen(2, 3)

	// Distribute the sq work across two goroutines that both read from in.
	c1 := count(in)
	c2 := count(in)

	//Consume the merged output from c1 and c2.
	for n := range merge(c1, c2) {
		fmt.Println(n) // 4 then 9, or 9 then 4
	}

	trace.Stop()
	return
}

//---------------------------------------------------------

// 加缓冲区，避免GoRutine阻塞，使发送者可以顺利退出
//---------------------------------------------------------
func gen3(nums ...int) <-chan int {
	out := make(chan int, len(nums))
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func merge3(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int, 1)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func RunChan3() {
	in := gen3(2, 3)

	c1 := count(in)
	c2 := count(in)

	// Distribute the sq work across two goroutines that both read from in.
	out := merge3(c1, c2)
	fmt.Println(<-out)

	p := pprof.Lookup("goroutine")
	p.WriteTo(os.Stdout, 1)
}

//---------------------------------------------------------

// 明确退出
//---------------------------------------------------------
func merge4(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func count4(done <-chan struct{}, nums <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range nums {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func RunChan41() {
	done := make(chan struct{})
	defer close(done)

	in := gen(2, 3)

	// Distribute the sq work across two goroutines that both read from in.
	c1 := count4(done, in)
	c2 := count4(done, in)

	out := merge4(done, c1, c2)
	fmt.Println(<-out) // 4 or 9

	p := pprof.Lookup("goroutine")
	p.WriteTo(os.Stdout, 1)
}

func RunChan4() {
	in := gen(2, 3)

	// Distribute the sq work across two goroutines that both read from in.
	c1 := count(in)
	c2 := count(in)

	// Consume the first value from output.
	done := make(chan struct{}, 2)
	out := merge4(done, c1, c2)
	fmt.Println(<-out) // 4 or 9

	// Tell the remaining senders we're leaving.
	done <- struct{}{}
	done <- struct{}{}

	p := pprof.Lookup("goroutine")
	p.WriteTo(os.Stdout, 1)
}

//---------------------------------------------------------
