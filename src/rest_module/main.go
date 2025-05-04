package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/robfig/cron"

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
	// Создание объектов API операций
	var operRepository = InitOperationRepository(dbManager)
	var operManager = OperationManagerNewInstance(mailSender, userRepository, accountRepository, operRepository)
	var operController = OperationControllerNewInstance(operManager)
	// Создание объектов API кредитов
	var paymentRepository = InitPaymentRepository(dbManager)
	var creditRepository = InitCreditRepository(dbManager)
	var creditManager = CreditManagerNewInstance(mailSender, userRepository, accountRepository, creditRepository, paymentRepository)
	var creditController = CreditControllerNewInstance(creditManager)

	// Запускаем планировщик
	c := cron.New()
	// Установка задания списания платежей
	c.AddFunc("* 0 * * *", func() {
		log.Println("Старт задания списания платежей")
		err = creditManager.PaymentForCredit()
		if err != nil {
			log.Panicln(err.Error())
		}
	})
	// Старт планировщика
	c.Start()

	// Главный контроллер приложения
	api := ApiNewInstance(usersController, accountController, cardController, operController, creditController)
	// Запуск сетевой службы и HTTP-сервера
	// на всех локальных IP-адресах на порту 8080.
	err = http.ListenAndServe(":8080", api.Router())
	if err != nil {
		log.Fatal(err)
	}

	// Остановка планировщика
	c.Stop()
}
