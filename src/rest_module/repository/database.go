package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "admin"
	password = "admin"
	dbname   = "hakaton_db"
)

type DBManager struct {
	database *sql.DB
}

// Конструктор БД.
func NewDBManager() *DBManager {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	manager := DBManager{}
	manager.database = db
	return &manager
}

func (manager *DBManager) CloseConnection() {
	manager.database.Close()
}

// Конструктор БД.
func ConnectDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	rows, _ := db.Query(`select id, transaction_time from transactions`)
	defer rows.Close()
	for rows.Next() {
		var id int
		var transaction_time string

		err = rows.Scan(&id, &transaction_time)
		if err != nil {
			panic(err)
		}

		fmt.Println(id, transaction_time)
	}
}
