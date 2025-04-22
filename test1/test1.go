package main

import (
	"fmt"
	"unsafe"
)

type MyStruct struct {
	b uintptr
	s string
}

func main() {
	str := "hello world"
	fmt.Println("Размер строки 'hello world':", unsafe.Sizeof(str))      // 16 байт
	fmt.Println("Размер структуры MyStruct:", unsafe.Sizeof(MyStruct{})) // 42 байта
}
