package main

import (
	"database/sql"
	"fmt"
	"sync"

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
	m        sync.Mutex //мьютекс для синхронизации доступа
	database *sql.DB
	idUser   int // текущее значение ID для нового заказа
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
	manager.idUser = 1
	return &manager
}

func (manager *DBManager) CloseConnection() {
	manager.database.Close()
}

// Сохранение нового пользователя в БД
func (manager *DBManager) InsertUser(user *User) int {
	manager.m.Lock()
	defer manager.m.Unlock()

	insertStmt := `insert into "users" ("username", "password") values($1, $2) returning "id"`

	id := 0
	err := manager.database.QueryRow(insertStmt, user.Username, user.Password).Scan(&id)
	if err != nil {
		panic(err)
	}

	return id
}

// Поиск пользователя по идентификатору
func (manager *DBManager) GetUserByID(id int) *User {
	manager.m.Lock()
	defer manager.m.Unlock()

	selectStmt := `select "id", "username", "password" from "users" where "id"=$1`
	rows, err := manager.database.Query(selectStmt, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		var username string
		var password string

		err = rows.Scan(&id, &username, &password)
		if err != nil {
			panic(err)
		}

		return &User{
			ID:       id,
			Username: username,
			Password: password,
		}
	}

	return nil
}

// Поиск пользователя по имени
func (manager *DBManager) GetUserByName(name string) *User {
	manager.m.Lock()
	defer manager.m.Unlock()

	selectStmt := `select "id", "username", "password" from "users" where "username"=$1`
	rows, err := manager.database.Query(selectStmt, name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		var username string
		var password string

		err = rows.Scan(&id, &username, &password)
		if err != nil {
			panic(err)
		}

		return &User{
			ID:       id,
			Username: username,
			Password: password,
		}
	}

	return nil
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
