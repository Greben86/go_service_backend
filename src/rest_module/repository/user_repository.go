package repository

import (
	"database/sql"
	. "rest_module/domain/model"
)

type UserRepository struct {
	db *DBManager // база данных
}

func InitUserRepository(db *DBManager) *UserRepository {
	repo := UserRepository{}
	repo.db = db
	return &repo
}

func (repo *UserRepository) Database() *sql.DB {
	if repo.db == nil {
		panic("База данных не подключена!")
	}
	return repo.db.database
}

// Сохранение нового пользователя в БД
func (repo *UserRepository) InsertUser(user *User) int {
	insertStmt := `insert into "users" ("username", "password") values($1, $2) returning "id"`

	id := 0
	err := repo.Database().QueryRow(insertStmt, user.Username, user.Password).Scan(&id)
	if err != nil {
		panic(err)
	}

	return id
}

// Поиск пользователя по идентификатору
func (repo *UserRepository) GetUserByID(id int) *User {
	selectStmt := `select "id", "username", "password" from "users" where "id"=$1`
	rows, err := repo.Database().Query(selectStmt, id)
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
func (repo *UserRepository) GetUserByName(name string) *User {
	selectStmt := `select "id", "username", "password" from "users" where "username"=$1`
	rows, err := repo.Database().Query(selectStmt, name)
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
