package app

import (
	"fmt"
	"strconv"
	"sync"
)

func Square(num ...string) (string, error) {
	if len(num) == 0 {
		return "nothing to do", nil
	}
	// convert the values recieved to ints
	var err error
	nums := make([]int, len(num), len(num))
	for i, n := range num {
		nums[i], err = strconv.Atoi(n)
		if err != nil {
			return "", err
		}
	}

	// Set up a shared done channel that the pipeline uses. This channel is closed
	// to signal to the other go routines that they should exit
	done := make(chan struct{})
	defer close(done)
	// start the work
	in := gen(nums...)

	// fan out: in a program where variable number of channels are to be supported
	// this would be done differently.
	c1 := sq(done, in)
	c2 := sq(done, in)

	for n := range merge(done, c1, c2) {
		fmt.Println(n)
	}
	return "", nil
}

func gen(nums ...int) <-chan int {
	out := make(chan int, len(nums))
	for _, n := range nums {
		out <- n
	}
	close(out)
	return out
}

func sq(done <-chan struct{}, in <-chan int) <-chan int {
	// defer takes care of closing the channel when it receives the done signal.
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

// merge a list of channels into a single one. Use a WG to ensure that all sends are
// done before closing to avoid panic.
func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs. output
	// copies values from c to out until c or done is closed, then calls wg.Done.
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

	// Start a goroutine to close out once all the output goroutines are done.
	// This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
