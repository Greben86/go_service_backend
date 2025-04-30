package main

import (
	"log"
	"net/http"
	. "rest_module/repository"
	. "rest_module/rest"
	. "rest_module/service"
)

func main() {
	// Инициализация соединения с БД
	var dbManager = NewDBManager()
	defer dbManager.CloseConnection()
	var userRepository = InitUserRepository(dbManager)
	var userManager = UserManagerNewInstance(userRepository)
	// Создание объекта API, использующего БД в памяти.
	api := ApiNewInstance(userManager)
	// Запуск сетевой службы и HTTP-сервера
	// на всех локальных IP-адресах на порту 8081.
	err := http.ListenAndServe(":8081", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}
