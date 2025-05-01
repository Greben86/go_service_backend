package domain_model

// Пользователь
type User struct {
	ID       int64
	Username string
	Password string
	Email    string
}

// Счет пользователя
type Account struct {
	ID     int64
	Name   string
	Bank   string
	UserId int64
}
