package main

import "fmt"

type X struct {
	val int
}

func (x X) S() {
	fmt.Println(x.val)
}

func main() {
	x := X{123}
	defer (&x).S()           // 123
	defer x.S()              // 123
	defer func() { x.S() }() // 456
	x.val = 456
}
