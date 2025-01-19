package main

import (
	"errors"
	"fmt"
	"strings"
)

var e error

func main() {
	e := era()
	fmt.Printf("err %+v type %T () %v\n", e, e, e.Error())
	//ew := errors.New("wtf")

	//	if errors.Is(era(), ew) {
	if strings.Contains(e.Error(), "wtf") {
		fmt.Printf("the same %+v %+v", era(), e)
	}
	// fmt.Printf("different %+v %+v", era(1), e)

}

func era() error {
	e := errors.New("wtf")
	//	fm := fmt.Errorf("%w", e)
	return e
}
