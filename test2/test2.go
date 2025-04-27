package main

import (
	"fmt"
	"time"
)

func main() {
	// Пример времени с локальным часовым поясом
	monthTime := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.Local)

	// Вызов без UTC (используется локальное время)
	localMonth := monthTime.Month()
	fmt.Println("Месяц в локальном времени:", localMonth)

	// Вызов с UTC (переводим в UTC)
	utcMonth := monthTime.UTC().Month()
	fmt.Println("Месяц в UTC:", utcMonth)
}
