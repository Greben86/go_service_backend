package repository

import (
	"database/sql"
	. "rest_module/model"
	"time"
)

type PaymentRepository struct {
	Db *DBManager // база данных
}

func InitPaymentRepository(db *DBManager) *PaymentRepository {
	repo := PaymentRepository{}
	repo.Db = db
	return &repo
}

func (repo *PaymentRepository) Database() *sql.DB {
	if repo.Db == nil {
		panic("База данных не подключена!")
	}

	return repo.Db.database
}

// Сохранение нового платежа по кредиту
func (repo *PaymentRepository) InsertPayment(payment *PaymentSchedule) (int64, error) {
	insertStmt := `insert into "payment_schedules" ("expiration_time", "amount", "payment_status", "user_id", "credit_id") values($1, $2, $3, $4, $5) returning "id"`

	var id int64 = 0
	err := repo.Database().QueryRow(insertStmt, payment.ExpirationTime, payment.Amount, payment.PaymentStatus, payment.UserId, payment.CreditId).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// Обновление платежа по кредиту
func (repo *PaymentRepository) UpdatePayment(payment *PaymentSchedule) error {
	insertStmt := `update "payment_schedules" set "expiration_time"=$1, "amount"=$2, "payment_status"=$3 where "id" = $1`

	_, err := repo.Database().Exec(insertStmt, payment.ExpirationTime, payment.Amount, payment.PaymentStatus, payment.ID)
	if err != nil {
		return err
	}

	return nil
}

// Поиск платежа по идентификатору
func (repo *PaymentRepository) GetPaymentByID(user_id, id int64) (*PaymentSchedule, error) {
	selectStmt := `select "id", "expiration_time", "amount", "payment_status", "credit_id" from "payment_schedules" where "user_id" = $1 and "id" = $2`
	rows, err := repo.Database().Query(selectStmt, user_id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var id int64
		var expiration_time time.Time
		var amount float64
		var payment_status int
		var credit_id int64

		err = rows.Scan(&id, &expiration_time, &amount, &payment_status, &credit_id)
		if err != nil {
			return nil, err
		}

		return &PaymentSchedule{
			ID:             id,
			ExpirationTime: expiration_time,
			Amount:         amount,
			PaymentStatus:  payment_status,
			CreditId:       credit_id,
			UserId:         user_id,
		}, nil
	}

	return nil, nil
}

// Список платежей пользователя
func (repo *PaymentRepository) GetPaymentsByUserId(user_id int64) (*[]PaymentSchedule, error) {
	selectStmt := `select "id", "expiration_time", "amount", "payment_status", "credit_id" from "payment_schedules" where "user_id" = $1`
	rows, err := repo.Database().Query(selectStmt, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []PaymentSchedule
	for rows.Next() {
		var id int64
		var expiration_time time.Time
		var amount float64
		var payment_status int
		var credit_id int64

		err = rows.Scan(&id, &expiration_time, &amount, &payment_status, &credit_id)
		if err != nil {
			return nil, err
		}

		payment := PaymentSchedule{
			ID:             id,
			ExpirationTime: expiration_time,
			Amount:         amount,
			PaymentStatus:  payment_status,
			CreditId:       credit_id,
			UserId:         user_id,
		}
		payments = append(payments, payment)
	}

	return &payments, nil
}

// Список платежей пользователя
func (repo *PaymentRepository) GetPaymentsByUserIdAndCreditId(user_id, credit_id int64) (*[]PaymentSchedule, error) {
	selectStmt := `select "id", "expiration_time", "amount", "payment_status" from "payment_schedules" where "user_id" = $1 and "credit_id" = $2`
	rows, err := repo.Database().Query(selectStmt, user_id, credit_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []PaymentSchedule
	for rows.Next() {
		var id int64
		var expiration_time time.Time
		var amount float64
		var payment_status int

		err = rows.Scan(&id, &expiration_time, &amount, &payment_status)
		if err != nil {
			return nil, err
		}

		payment := PaymentSchedule{
			ID:             id,
			ExpirationTime: expiration_time,
			Amount:         amount,
			PaymentStatus:  payment_status,
			CreditId:       credit_id,
			UserId:         user_id,
		}
		payments = append(payments, payment)
	}

	return &payments, nil
}

// Список платежей пользователя
func (repo *PaymentRepository) GetActivePayments() (*[]PaymentSchedule, error) {
	selectStmt := `select "id", "expiration_time", "amount", "payment_status", "credit_id", "user_id" from "payment_schedules" where "payment_status" = $1 and "expiration_time" <= CURRENT_DATE`
	rows, err := repo.Database().Query(selectStmt, 0)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []PaymentSchedule
	for rows.Next() {
		var id int64
		var expiration_time time.Time
		var amount float64
		var payment_status int
		var credit_id int64
		var user_id int64

		err = rows.Scan(&id, &expiration_time, &amount, &payment_status)
		if err != nil {
			return nil, err
		}

		payment := PaymentSchedule{
			ID:             id,
			ExpirationTime: expiration_time,
			Amount:         amount,
			PaymentStatus:  payment_status,
			CreditId:       credit_id,
			UserId:         user_id,
		}
		payments = append(payments, payment)
	}

	return &payments, nil
}
