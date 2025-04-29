package main

import (
	"fmt"
	"time"
)

type Point struct {
	X, Y int
}

// Конструктор возвращает просто `Point`, без указателя
func NewPoint(x, y int) *Point {
	p := Point{X: x, Y: y}
	return &p
}

func main() {
	r := NewPoint(1, 2)

	// Пример времени с локальным часовым поясом
	monthTime := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.Local)

	// Вызов без UTC (используется локальное время)
	localMonth := monthTime.Month()
	fmt.Println("Месяц в локальном времени:", localMonth)

	// Вызов с UTC (переводим в UTC)
	utcMonth := monthTime.UTC().Month()
	fmt.Println("Месяц в UTC:", utcMonth)
}
