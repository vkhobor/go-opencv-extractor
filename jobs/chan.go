package jobs

import (
	"reflect"
	"sync"
)

type Pipe[I any, O any] func(<-chan I) <-chan O

func AggregateChan[O any](p <-chan O) <-chan []O {
	aggr := make(chan []O)
	go func() {
		list := []O{}
		for video := range p {
			list = append(list, video)
			aggr <- list
		}
		close(aggr)
	}()

	return aggr
}

func TapChan[O any](
	p <-chan O,
	tap func(input O)) <-chan O {
	tapped := make(chan O)
	go func() {
		for val := range p {
			tap(val)
			tapped <- val
		}
		close(tapped)
	}()
	return tapped
}

func FilterChan[O any](
	p <-chan O,
	filter func(input O) bool) <-chan O {
	filtered := make(chan O)
	go func() {
		for val := range p {
			if filter(val) {
				filtered <- val
			}
		}
		close(filtered)
	}()
	return filtered
}

func MultiplexChan[O any](
	p <-chan O,
	split func(input O) string, keySpace []string) map[string]<-chan O {
	chanMapAsType := (make(map[string]chan O))
	for _, key := range keySpace {
		chanMapAsType[key] = make(chan O)
	}

	go func() {
		for val := range p {
			key := split(val)
			if _, ok := chanMapAsType[key]; ok {
				chanMapAsType[key] <- val
			}
		}
		for _, ch := range chanMapAsType {
			close(ch)
		}
	}()

	chanMap := make(map[string]<-chan O)
	for key, ch := range chanMapAsType {
		chanMap[key] = ch
	}

	return chanMap
}

func MergeChans[O any](p ...<-chan O) <-chan O {
	merged := make(chan O)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(p))

	for _, ch := range p {
		go func(ch <-chan O) {
			for val := range ch {
				merged <- val
			}
			waitGroup.Done()
		}(ch)
	}
	go func() {
		waitGroup.Wait()
		close(merged)
	}()
	return merged
}

func CountChan[O any](p <-chan O) <-chan int {
	count := make(chan int)
	go func() {
		c := 0
		for range p {
			c++
			count <- c
		}
		close(count)
	}()
	return count
}

func LatestFromChans[O any](p ...<-chan O) <-chan []O {
	latest := make(chan []O, len(p))

	go func() {
		cases := make([]reflect.SelectCase, len(p))

		for i, ch := range p {
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			}
		}
		current := make([]O, len(p))
		for {
			if len(cases) == 0 {
				break
			}

			chosen, value, ok := reflect.Select(cases)
			if ok {
				current[chosen] = value.Interface().(O)
				latest <- current
			} else {
				cases = append(cases[:chosen], cases[chosen+1:]...)
			}
		}
		close(latest)
	}()

	return latest
}

func SwallowChan[O any](p <-chan O) {
	go func() {
		for range p {
		}
	}()
}

func MapChannel[I any, O any](p <-chan I, f func(input I) O) <-chan O {
	mapped := make(chan O)
	go func() {
		for val := range p {
			mapped <- f(val)
		}
		close(mapped)
	}()
	return mapped
}
