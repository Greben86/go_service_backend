package domain_model

import "time"

// Пользователь
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"email"`
}

// Счет пользователя
type Account struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Bank    string  `json:"bank"`
	Balance float64 `json:"balance"`
	UserId  int64   `json:"-"`
}

// Карта пользователя
type Card struct {
	ID              int64  `json:"id"`
	Number          string `json:"number"`
	ExpirationMonth int    `json:"expiration_month"`
	ExpirationYear  int    `json:"expiration_year"`
	CVV             string `json:"-"`
	AccountId       int64  `json:"account_id"`
	UserId          int64  `json:"-"`
}

// Операция
type Operation struct {
	ID            int64   `json:"id"`
	SumValue      float64 `json:"sum_value"`
	OperationType string  `json:"operation_type"`
	AccountId     int64   `json:"account_id"`
	UserId        int64   `json:"-"`
}

// Кредит
type Credit struct {
	ID         int64     `json:"id"`
	Amount     float64   `json:"amount"`
	Rate       float64   `json:"rate"`
	MonthCount int       `json:"month_count"`
	StartDate  time.Time `json:"start_date"`
	AccountId  int64     `json:"account_id"`
	UserId     int64     `json:"-"`
}

// Платеж по кредиту
type PaymentSchedule struct {
	ID             int64     `json:"id"`
	ExpirationTime time.Time `json:"expiration_time"`
	Amount         float64   `json:"amount"`
	PaymentStatus  int       `json:"payment_status"`
	CreditId       int64     `json:"credit_id"`
	UserId         int64     `json:"-"`
}
