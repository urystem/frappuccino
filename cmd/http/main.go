package main

import (
	"cafeteria/internal/handlers"
	"cafeteria/pkg/config"
	"cafeteria/pkg/lib/logger"
	"database/sql"
	"log"
	"os"
)

func main() {
	cfg := config.LoadConfig()
	logger := logger.SetupPrettySlog(os.Stdout)

	db, err := sql.Open("postgres", cfg.MakeConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	httpSrv := handlers.NewAPIServer(
		"0.0.0.0:8080",
		db,
		logger,
	)
	httpSrv.Run()
}
