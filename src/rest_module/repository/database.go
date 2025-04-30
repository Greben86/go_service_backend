package repository

import (
	"database/sql"
	"fmt"

	"github.com/qustavo/dotsql"

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

// Миграция БД.
func (manager *DBManager) InitDB() error {
	dot, err := dotsql.LoadFromFile("repository/init_database.sql")
	if err != nil {
		fmt.Errorf("Файл миграции БД не найден %s", err.Error())
		return err
	}

	if _, err := dot.Exec(manager.database, "create-users-table"); err != nil {
		fmt.Errorf("Ошибка создания таблицы пользователей %s", err.Error())
		return err
	}

	if _, err := dot.Exec(manager.database, "create-accounts-table"); err != nil {
		fmt.Errorf("Ошибка создания таблицы счетов %s", err.Error())
		return err
	}

	fmt.Println("База данных обновлена!")

	return nil
}
