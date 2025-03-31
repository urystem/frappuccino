package main

import (
	// "database/sql" // Import the database/sql package to interact with SQL databases
	"fmt"      // Import the fmt package for formatted I/O (printing messages, etc.)
	"log"      // Import the log package for logging errors
	"net/http" // listen and serve
	"os"       // Import the os package to access environment variables and other OS functions

	"hot-coffee/internal/router" // for mux

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib" // Import the pq PostgreSQL driver (side-effect import, it registers itself with database/sql)
	"github.com/jmoiron/sqlx"
)

func main() {
	// Data Source Name
	dsn := fmt.Sprintf("postgres://%[1]s:%s@%s:%s/%s", // index starts at 1 in formatting
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"), // you can skip this env, if you use default posgresql port. 5432
		os.Getenv("DB_NAME"))

	db, err := sqlx.Open("pgx", dsn) // Attempt to set up the database connection
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	routes := router.Allrouter(db)

	log.Fatal(http.ListenAndServe(":8080", routes))
}

func GetMenuItems(db *sqlx.DB) ([]MenuItem, error) {
	query := `SELECT id, name, description, tags, allergens, price FROM menu_items`
	rows, err := db.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menuItems []MenuItem
	for rows.Next() {
		var item MenuItem
		var tags pgx.TextArray // 👈 Используем pgx.TextArray
		var allergens pgx.TextArray

		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&tags,      // 👈 pgx.TextArray автоматически конвертируется в []string
			&allergens, // 👈 pgx.TextArray автоматически конвертируется в []string
			&item.Price,
		)
		if err != nil {
			return nil, err
		}

		item.Tags = tags.Elements // 👈 Присваиваем []string
		item.Allergens = allergens.Elements
		menuItems = append(menuItems, item)
	}

	return menuItems, nil
}
