package main

import "fmt"

func main() {
    a := "hello"
    b := a // присваивание одной строки другой
    fmt.Println(&a, &b) // Ожидаем одинаковые адреса
}
