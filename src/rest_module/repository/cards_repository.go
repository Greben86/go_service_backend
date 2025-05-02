package repository

import (
	"database/sql"
	. "rest_module/model"
)

type CardRepository struct {
	Db *DBManager // база данных
}

func InitCardRepository(db *DBManager) *CardRepository {
	repo := CardRepository{}
	repo.Db = db
	return &repo
}

func (repo *CardRepository) Database() *sql.DB {
	if repo.Db == nil {
		panic("База данных не подключена!")
	}

	return repo.Db.database
}

// Сохранение новой карты в БД
func (repo *CardRepository) InsertCard(card *Card) (int64, error) {
	insertStmt := `insert into "cards" ("number", "expiration_month", "expiration_year", "cvv", "user_id", "account_id") values($1, $2, $3, $4, $5, $6) returning "id"`

	var id int64 = 0
	err := repo.Database().QueryRow(insertStmt, card.Number, card.ExpirationMonth, card.ExpirationYear, card.CVV, card.UserId, card.AccountId).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// Поиск карты по идентификатору
func (repo *CardRepository) GetCardByID(user_id, id int64) (*Card, error) {
	selectStmt := `select "id", "number", "expiration_month", "expiration_year", "cvv", "account_id" from "cards" where "user_id" = $1 and "id" = $2`
	rows, err := repo.Database().Query(selectStmt, user_id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int64
		var number string
		var expiration_month int
		var expiration_year int
		var cvv string
		var account_id int64

		err = rows.Scan(&id, &number, &expiration_month, &expiration_year, &cvv, &account_id)
		if err != nil {
			return nil, err
		}

		return &Card{
			ID:              id,
			Number:          number,
			ExpirationMonth: expiration_month,
			ExpirationYear:  expiration_year,
			CVV:             cvv,
			AccountId:       account_id,
			UserId:          user_id,
		}, nil
	}

	return nil, nil
}

// Поиск карты по номеру
func (repo *CardRepository) GetCardByNumber(user_id int64, number string) (*Card, error) {
	selectStmt := `select "id", "number", "expiration_month", "expiration_year", "cvv", "account_id" from "cards" where "user_id" = $1 and "number" ~ $2`
	rows, err := repo.Database().Query(selectStmt, user_id, number)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int64
		var number string
		var expiration_month int
		var expiration_year int
		var cvv string
		var account_id int64

		err = rows.Scan(&id, &number, &expiration_month, &expiration_year, &cvv, &account_id)
		if err != nil {
			return nil, err
		}

		return &Card{
			ID:              id,
			Number:          number,
			ExpirationMonth: expiration_month,
			ExpirationYear:  expiration_year,
			CVV:             cvv,
			AccountId:       account_id,
			UserId:          user_id,
		}, nil
	}

	return nil, nil
}

// Список карт пользователя
func (repo *CardRepository) GetCardsByUserId(user_id int64) (*[]Card, error) {
	selectStmt := `select "id", "number", "expiration_month", "expiration_year", "cvv", "account_id" from "cards" where "user_id" = $1`
	rows, err := repo.Database().Query(selectStmt, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []Card
	for rows.Next() {
		var id int64
		var number string
		var expiration_month int
		var expiration_year int
		var cvv string
		var account_id int64

		err = rows.Scan(&id, &number, &expiration_month, &expiration_year, &cvv, &account_id)
		if err != nil {
			return nil, err
		}

		card := Card{
			ID:              id,
			Number:          number,
			ExpirationMonth: expiration_month,
			ExpirationYear:  expiration_year,
			CVV:             cvv,
			AccountId:       account_id,
			UserId:          user_id,
		}
		cards = append(cards, card)
	}

	return &cards, nil
}
