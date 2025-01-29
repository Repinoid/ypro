package main

import "fmt"

type chaina struct {
	ch chan int
}


func (s chaina)  add() chan int {
	chacha := make(chan int)
	go func() {
		defer close(chacha)
		for value := range s.ch {
			chacha <- value
		}
	}()
	return chacha
}

func main() {
	numbers := []int{1, 2, 3, 4, 5}
	s := chaina{ch:make(chan int) }

	go func() {
		defer close(s.ch)
		for _, num := range numbers {
			s.ch <- num
		}
	}()

	addCh := s.add()
	
	for res := range addCh {
		fmt.Print(res, " ")
	}
}
