package repository

import (
	"database/sql"
	. "rest_module/model"
)

type AccountRepository struct {
	Db *DBManager // база данных
}

func InitAccountRepository(db *DBManager) *AccountRepository {
	repo := AccountRepository{}
	repo.Db = db
	return &repo
}

func (repo *AccountRepository) Database() *sql.DB {
	if repo.Db == nil {
		panic("База данных не подключена!")
	}

	return repo.Db.database
}

// Сохранение нового счета в БД
func (repo *AccountRepository) InsertAccount(account *Account) (int64, error) {
	insertStmt := `insert into "accounts" ("name", "bank", "balance", "user_id") values($1, $2, $3, $4) returning "id"`

	var id int64 = 0
	err := repo.Database().QueryRow(insertStmt, account.Name, account.Bank, account.Balance, account.UserId).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// Обновление счета в БД
func (repo *AccountRepository) UpdateAccount(account *Account) error {
	updateStmt := `update "accounts" set "name" = $1, "bank" = $2, "balance" = $3 where "id" = $4`

	_, err := repo.Database().Exec(updateStmt, account.Name, account.Bank, account.Balance, account.ID)
	if err != nil {
		return err
	}

	return nil
}

// Поиск счета по идентификатору
func (repo *AccountRepository) GetAccountByIDAndUserID(user_id, id int64) (*Account, error) {
	selectStmt := `select "id", "name", "bank", "balance" from "accounts" where "user_id" = $1 and "id" = $2`
	rows, err := repo.Database().Query(selectStmt, user_id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int64
		var name string
		var bank string
		var balance float64

		err = rows.Scan(&id, &name, &bank, &balance)
		if err != nil {
			return nil, err
		}

		return &Account{
			ID:      id,
			Name:    name,
			Bank:    bank,
			Balance: balance,
			UserId:  user_id,
		}, nil
	}

	return nil, nil
}

// Поиск счета по идентификатору
func (repo *AccountRepository) GetAccountByID(id int64) (*Account, error) {
	selectStmt := `select "id", "name", "bank", "balance", "user_id" from "accounts" where "id" = $1`
	rows, err := repo.Database().Query(selectStmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int64
		var name string
		var bank string
		var balance float64
		var user_id int64

		err = rows.Scan(&id, &name, &bank, &balance, &user_id)
		if err != nil {
			return nil, err
		}

		return &Account{
			ID:      id,
			Name:    name,
			Bank:    bank,
			Balance: balance,
			UserId:  user_id,
		}, nil
	}

	return nil, nil
}

// Поиск счета по названию
func (repo *AccountRepository) GetAccountByName(user_id int64, name string) (*Account, error) {
	selectStmt := `select "id", "name", "bank", "balance" from "accounts" where "user_id" = $1 and "name" ~ $2`
	rows, err := repo.Database().Query(selectStmt, user_id, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int64
		var name string
		var bank string
		var balance float64

		err = rows.Scan(&id, &name, &bank, &balance)
		if err != nil {
			return nil, err
		}

		return &Account{
			ID:      id,
			Name:    name,
			Bank:    bank,
			Balance: balance,
			UserId:  user_id,
		}, nil
	}

	return nil, nil
}

// Список счетов пользователя
func (repo *AccountRepository) GetAccountsByUserId(user_id int64) (*[]Account, error) {
	selectStmt := `select "id", "name", "bank", "balance" from "accounts" where "user_id" = $1`
	rows, err := repo.Database().Query(selectStmt, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var id int64
		var name string
		var bank string
		var balance float64

		err = rows.Scan(&id, &name, &bank, &balance)
		if err != nil {
			return nil, err
		}

		account := Account{
			ID:      id,
			Name:    name,
			Bank:    bank,
			Balance: balance,
			UserId:  user_id,
		}
		accounts = append(accounts, account)
	}

	return &accounts, nil
}
