package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	_, err := sql.Open("pgx", "postgres://latte:latte@db:5432/frappuccino")
	if err != nil {
		log.Fatalf("%v", err)
		fmt.Println("qyzyq")
	}
	// for {
	// 	var now string
	// 	err = db.QueryRow("SELECT NOW()").Scan(&now)
	// 	if err != nil {
	// 		log.Fatalf("Ошибка запроса: %v", err)
	// 	}
	// 	// fmt.Println("hello")

	// 	fmt.Println("⏰ Время в БД:", now)
	// 	time.Sleep(time.Second * 5)
	// }
	// log.Fatalln(http.ListenAndServe(":"+"8080", router.Allrouter(db)))
}
