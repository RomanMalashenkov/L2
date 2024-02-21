package main

import (
	"fmt"
	"sync"
	"time"
)

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func main() {
	start := time.Now()

	goroutines := []<-chan interface{}{
		sig(2 * time.Hour),
		sig(5 * time.Minute),
		sig(1 * time.Second),
		sig(1 * time.Hour),
		sig(1 * time.Minute),
	}

	endChan := or(goroutines...)
	<-endChan

	fmt.Printf("fone after %v", time.Since(start))

}

func or(channels ...<-chan interface{}) <-chan interface{} {
	if len(channels) == 0 {
		c := make(chan interface{})
		close(c)
		return c
	}

	var (
		wg   = sync.WaitGroup{}
		once = sync.Once{}
		orc  = make(chan interface{})
	)

	go func() {
		for _, channel := range channels {
			wg.Add(1)
			go func(ch <-chan interface{}) {
				defer wg.Done()

				for obj := range ch {
					orc <- obj
				}

				once.Do(func() {
					close(orc)
				})
			}(channel)
		}
		wg.Wait()
	}()

	return orc
}
