package domain_model

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
