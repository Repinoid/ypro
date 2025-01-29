package main

import "fmt"

func add(inputCh chan int) chan int {
	resultCh := make(chan int)
	go func() {
		defer close(resultCh)
		for value := range inputCh {
			resultCh <- value
		}
	}()
	return resultCh
}

func generator(numbers []int) chan int {
	outputCh := make(chan int)
	go func() {
		defer close(outputCh)
		for _, num := range numbers {
			outputCh <- num
		}
	}()
	return outputCh
}

func main() {
	numbers := []int{1, 2, 3, 4, 5}
	inputCh := make(chan int)
	go func() {
		defer close(inputCh)
		for _, num := range numbers {
			inputCh <- num
		}
	}()

	//inputCh := generator(numbers)

	addCh := add(inputCh)
	for res := range addCh {
		fmt.Print(res, " ")
	}
}
