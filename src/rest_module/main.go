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
	err := dbManager.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	var mailSender, _ = InitMailSender()

	// Создание объектов API пользователя
	var userRepository = InitUserRepository(dbManager)
	var userManager = UserManagerNewInstance(userRepository)
	var usersController = UsersControllerNewInstance(userManager)
	// Создание объектов API счета
	var accountRepository = InitAccountRepository(dbManager)
	var accountManager = AccountManagerNewInstance(userRepository, accountRepository)
	var accountController = AccountControllerNewInstance(accountManager)
	// Создание объектов API карт
	var cardRepository = InitCardRepository(dbManager)
	var cardManager = CardManagerNewInstance(mailSender, userRepository, cardRepository)
	var cardController = CardControllerNewInstance(cardManager)

	// Главный контроллер приложения
	api := ApiNewInstance(usersController, accountController, cardController)
	// Запуск сетевой службы и HTTP-сервера
	// на всех локальных IP-адресах на порту 8080.
	err = http.ListenAndServe(":8080", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}
