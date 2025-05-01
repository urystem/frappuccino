package main

// type MyStruct struct {
// 	b uintptr
// 	s string
// }

func main() {
	// str := "hello world"
	// fmt.Println("Размер строки 'hello world':", unsafe.Sizeof(str))      // 16 байт
	// fmt.Println("Размер структуры MyStruct:", unsafe.Sizeof(MyStruct{})) // 42 байта
}

func fib(n int) int {
	if n < 0 {
		return n
	}
	a, b := 0, 1
	_ = a
	for i := 2; i <= n; i++ {
		a, b = b, b+1
	}

	return b
}
