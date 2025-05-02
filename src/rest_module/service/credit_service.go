package service

import (
	"fmt"
	"rest_module/repository"
	"sync"
	"time"

	. "rest_module/model"
)

type CreditManager struct {
	m           sync.Mutex // мьютекс для синхронизации доступа
	mailSender  *MailSender
	userRepo    *repository.UserRepository    // репозиторий пользователей
	accountRepo *repository.AccountRepository // репозиторий счетов
	creditRepo  *repository.CreditRepository  // репозиторий кредитов
	paymentRepo *repository.PaymentRepository // репозиторий палтежей
}

// Конструктор сервиса
func CreditManagerNewInstance(mailSender *MailSender, userRepo *repository.UserRepository,
	accountRepo *repository.AccountRepository, creditRepo *repository.CreditRepository, paymentRepo *repository.PaymentRepository) *CreditManager {
	manager := CreditManager{}
	manager.mailSender = mailSender
	manager.userRepo = userRepo
	manager.accountRepo = accountRepo
	manager.creditRepo = creditRepo
	manager.paymentRepo = paymentRepo
	return &manager
}

// Создание карты
func (manager *CreditManager) AddCredit(credit Credit, user_id int64) (*Credit, error) {
	manager.m.Lock()
	defer manager.m.Unlock()
	var err error

	manager.creditRepo.Db.BeginTransaction()
	user, _ := manager.userRepo.GetUserByID(user_id)
	if user == nil {
		manager.creditRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Пользователь с таким логином не найден")
	}

	account, _ := manager.accountRepo.GetAccountByIDAndUserID(user_id, credit.AccountId)
	if account == nil {
		manager.creditRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Счет не найден")
	}

	account.Balance += credit.Amount
	err = manager.accountRepo.UpdateAccount(account)
	if err != nil {
		manager.creditRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Ошибка обновления счета %s", err.Error())
	}

	credit.UserId = user_id
	credit.StartDate = time.Now()
	credit.Rate, err = manager.getRate()
	if err != nil {
		manager.creditRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Ошибка получения ставки центробанка %s", err.Error())
	}
	credit.ID, err = manager.creditRepo.InsertCredit(&credit)
	if err != nil {
		manager.creditRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Ошибка добавления кредита %s", err.Error())
	}

	// Рассчет графика платежей
	pay := credit.Amount / float64(credit.MonthCount)
	for i := 1; i <= credit.MonthCount; i++ {
		payment := PaymentSchedule{}
		payment.ExpirationTime = time.Now().AddDate(0, i, 0)
		payment.Amount = pay
		payment.PaymentStatus = 0
		payment.CreditId = credit.ID
		payment.UserId = user_id
		payment.ID, err = manager.paymentRepo.InsertPayment(&payment)
		if err != nil {
			manager.creditRepo.Db.RollbackTransaction()
			return nil, fmt.Errorf("Ошибка создания графика платежей %s", err.Error())
		}
	}

	manager.creditRepo.Db.CommitTransaction()
	return &credit, nil
}

// Рассчет процентной ставки
func (manager *CreditManager) getRate() (float64, error) {
	rateService := CentralBankRateService{}
	rate, err := rateService.GetCentralBankRate()
	rate += 5
	return rate, err
}

// Поиск кредита по идентификатору
func (manager *CreditManager) FindCreditById(user_id, id int64) (*Credit, error) {
	manager.m.Lock()
	defer manager.m.Unlock()

	manager.creditRepo.Db.BeginTransaction()
	card, _ := manager.creditRepo.GetCreditByID(user_id, id)
	if card == nil {
		manager.creditRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Кредит с таким идентификатором не найден")
	}
	manager.creditRepo.Db.CommitTransaction()

	return card, nil
}

// Поиск кредитов пользователя
func (manager *CreditManager) FindCreditsByUserId(user_id int64) (*[]Credit, error) {
	manager.m.Lock()
	defer manager.m.Unlock()

	manager.creditRepo.Db.BeginTransaction()
	cards, _ := manager.creditRepo.GetCreditsByUserId(user_id)
	if cards == nil {
		manager.creditRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Кредиты пользователя не найдены")
	}
	manager.creditRepo.Db.CommitTransaction()

	return cards, nil
}

// График платежей по кредиту
func (manager *CreditManager) PaymentScheduleByCreditId(user_id, credit_id int64) (*[]PaymentSchedule, error) {
	manager.m.Lock()
	defer manager.m.Unlock()

	manager.creditRepo.Db.BeginTransaction()
	cards, _ := manager.paymentRepo.GetPaymentsByUserIdAndCreditId(user_id, credit_id)
	if cards == nil {
		manager.creditRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Кредиты пользователя не найдены")
	}
	manager.creditRepo.Db.CommitTransaction()

	return cards, nil
}
