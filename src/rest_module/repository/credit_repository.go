package repository

import (
	"database/sql"
	. "rest_module/model"
	"time"
)

type CreditRepository struct {
	Db *DBManager // база данных
}

func InitCreditRepository(db *DBManager) *CreditRepository {
	repo := CreditRepository{}
	repo.Db = db
	return &repo
}

func (repo *CreditRepository) Database() *sql.DB {
	if repo.Db == nil {
		panic("База данных не подключена!")
	}

	return repo.Db.database
}

// Сохранение нового кредита
func (repo *CreditRepository) InsertCredit(credit *Credit) (int64, error) {
	insertStmt := `insert into "credits" ("amount", "rate", "month_count", "start_date", "account_id", "user_id") values($1, $2, $3, $4, $5, $6) returning "id"`

	var id int64 = 0
	err := repo.Database().QueryRow(insertStmt, credit.Amount, credit.Rate, credit.MonthCount, credit.StartDate, credit.AccountId, credit.UserId).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// Поиск кредита по идентификатору
func (repo *CreditRepository) GetCreditByID(user_id, id int64) (*Credit, error) {
	selectStmt := `select "id", "amount", "rate", "month_count", "start_date", "account_id" from "credits" where "user_id" = $1 and "id" = $2`
	rows, err := repo.Database().Query(selectStmt, user_id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int64
		var amount float64
		var rate float64
		var month_count int
		var start_date time.Time
		var account_id int64

		err = rows.Scan(&id, &amount, &rate, &month_count, &start_date, &account_id)
		if err != nil {
			return nil, err
		}

		return &Credit{
			ID:         id,
			Amount:     amount,
			Rate:       rate,
			MonthCount: month_count,
			StartDate:  start_date,
			AccountId:  account_id,
			UserId:     user_id,
		}, nil
	}

	return nil, nil
}

// Список кредитов пользователя
func (repo *CreditRepository) GetCreditsByUserId(user_id int64) (*[]Credit, error) {
	selectStmt := `select "id", "amount", "rate", "month_count", "start_date", "account_id" from "credits" where "user_id" = $1`
	rows, err := repo.Database().Query(selectStmt, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credits []Credit
	for rows.Next() {
		var id int64
		var amount float64
		var rate float64
		var month_count int
		var start_date time.Time
		var account_id int64

		err = rows.Scan(&id, &amount, &rate, &month_count, &start_date, &account_id)
		if err != nil {
			return nil, err
		}

		credit := Credit{
			ID:         id,
			Amount:     amount,
			Rate:       rate,
			MonthCount: month_count,
			StartDate:  start_date,
			AccountId:  account_id,
			UserId:     user_id,
		}
		credits = append(credits, credit)
	}

	return &credits, nil
}

// Список всех кредитов
func (repo *CreditRepository) GetCredits() (*[]Credit, error) {
	selectStmt := `select "id", "amount", "rate", "month_count", "start_date", "account_id", "user_id" from "credits"`
	rows, err := repo.Database().Query(selectStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credits []Credit
	for rows.Next() {
		var id int64
		var amount float64
		var rate float64
		var month_count int
		var start_date time.Time
		var account_id int64
		var user_id int64

		err = rows.Scan(&id, &amount, &rate, &month_count, &start_date, &account_id, &user_id)
		if err != nil {
			return nil, err
		}

		credit := Credit{
			ID:         id,
			Amount:     amount,
			Rate:       rate,
			MonthCount: month_count,
			StartDate:  start_date,
			AccountId:  account_id,
			UserId:     user_id,
		}
		credits = append(credits, credit)
	}

	return &credits, nil
}
