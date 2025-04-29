package main

import "sync"

type RequestSignUp struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseAuth struct {
	JwtToken string `json:"token"`
}

// Пользователь
type User struct {
	ID       int
	Username string
	Password string
}

// База данных заказов.
type DBMemory struct {
	m     sync.Mutex   //мьютекс для синхронизации доступа
	id    int          // текущее значение ID для нового заказа
	store map[int]User // БД пользователей
}

// Конструктор БД.
func NewDB() *DBMemory {
	db := DBMemory{
		id:    1, // первый номер пользователя
		store: map[int]User{},
	}
	return &db
}

// Users возвращает всех пользователей
func (db *DBMemory) Users() []User {
	db.m.Lock()
	defer db.m.Unlock()
	var data []User
	for _, v := range db.store {
		data = append(data, v)
	}
	return data
}

// NewUser создает нового пользователя
func (db *DBMemory) NewUser(u User) int {
	db.m.Lock()
	defer db.m.Unlock()
	u.ID = db.id
	db.store[u.ID] = u
	db.id++
	return u.ID
}

// GetUser обновляет данные пользователя по ID.
func (db *DBMemory) GetUser(id int) User {
	db.m.Lock()
	defer db.m.Unlock()
	return db.store[id]
}

// UpdateUser обновляет данные пользователя по ID.
func (db *DBMemory) UpdateUser(u User) {
	db.m.Lock()
	defer db.m.Unlock()
	db.store[u.ID] = u
}

// DeleteUser удаляет пользователя по ID.
func (db *DBMemory) DeleteUser(id int) {
	db.m.Lock()
	defer db.m.Unlock()
	delete(db.store, id)
}
