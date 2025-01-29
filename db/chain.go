package main

import "fmt"

type chaina struct {
	ch chan int
}

func (s chaina) add2() chaina {
	chacha := chaina{ch: make(chan int)}
	go func() {
		defer close(chacha.ch)
		for value := range s.ch {
			chacha.ch <- value + 2
		}
	}()
	return chacha
}
func (s chaina) add3() chaina {
	chacha := chaina{ch: make(chan int)}
	go func() {
		defer close(chacha.ch)
		for value := range s.ch {
			chacha.ch <- value + 3
		}
	}()
	return chacha
}
func (s chaina) summ() chaina {
	chacha := chaina{ch: make(chan int)}
	go func() {
		defer close(chacha.ch)
		sum := 0
		for value := range s.ch {
			sum += value
		}
		chacha.ch <- sum
	}()
	return chacha
}

func main() {
	numbers := []int{1, 2, 3, 4, 5}
	s := chaina{ch: make(chan int)}

	go func() {
		defer close(s.ch)
		for _, num := range numbers {
			s.ch <- num
		}
	}()

	addCh := s.add2().add3().summ()

	for res := range addCh.ch {
		fmt.Print(res, " ")
	}
}
