package main

import (
	"log"
	"net/http"
)

func main() {
	// Инициализация БД в памяти.
	db := NewDB()
	// Создание объекта API, использующего БД в памяти.
	api := ApiNewInstance(db)
	// Запуск сетевой службы и HTTP-сервера
	// на всех локальных IP-адресах на порту 8081.
	err := http.ListenAndServe(":8081", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}
