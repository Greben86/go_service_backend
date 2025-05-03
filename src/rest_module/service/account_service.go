package service

import (
	"fmt"
	"rest_module/repository"
	"sync"

	log "github.com/sirupsen/logrus"

	. "rest_module/model"
)

type AccountManager struct {
	m           sync.Mutex                    // мьютекс для синхронизации доступа
	userRepo    *repository.UserRepository    // репозиторий пользователей
	accountRepo *repository.AccountRepository // репозиторий счетов
}

// Конструктор сервиса
func AccountManagerNewInstance(userRepo *repository.UserRepository, accountRepo *repository.AccountRepository) *AccountManager {
	manager := AccountManager{}
	manager.userRepo = userRepo
	manager.accountRepo = accountRepo
	return &manager
}

// Создание счета
func (manager *AccountManager) AddAccount(account Account, user_id int64) (*Account, error) {
	log.Println("Создание счета")
	manager.m.Lock()
	defer manager.m.Unlock()

	manager.accountRepo.Db.BeginTransaction()
	user, _ := manager.userRepo.GetUserByID(user_id)
	if user == nil {
		manager.accountRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Пользователь с таким логином не найден")
	}

	exist, _ := manager.accountRepo.GetAccountByName(user_id, account.Name)
	if exist != nil {
		manager.accountRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Счет с таким названием уже есть")
	}

	var err error
	account.Balance = 0.
	account.UserId = user_id
	account.ID, err = manager.accountRepo.InsertAccount(&account)
	if err != nil {
		manager.accountRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Ошибка добавления счета %s", err.Error())
	}
	manager.accountRepo.Db.CommitTransaction()
	return &account, nil
}

// Поиск счета по идентификатору
func (manager *AccountManager) FindAccountById(user_id, id int64) (*Account, error) {
	log.Println("Поиск счета по идентификатору")
	manager.m.Lock()
	defer manager.m.Unlock()

	manager.accountRepo.Db.BeginTransaction()
	account, _ := manager.accountRepo.GetAccountByIDAndUserID(user_id, id)
	if account == nil {
		manager.accountRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Счет с таким идентификатором не найден")
	}
	manager.accountRepo.Db.CommitTransaction()

	return account, nil
}

// Поиск счета по названию
func (manager *AccountManager) FindAccountByName(user_id int64, name string) (*Account, error) {
	log.Println("Поиск счета по названию")
	manager.m.Lock()
	defer manager.m.Unlock()

	manager.accountRepo.Db.BeginTransaction()
	user, _ := manager.accountRepo.GetAccountByName(user_id, name)
	if user == nil {
		manager.accountRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Счет с таким названием не найден")
	}
	manager.accountRepo.Db.CommitTransaction()

	return user, nil
}

// Поиск счетов пользователя
func (manager *AccountManager) FindAccountsByUserId(user_id int64) (*[]Account, error) {
	log.Println("Поиск счетов пользователя")
	manager.m.Lock()
	defer manager.m.Unlock()

	manager.accountRepo.Db.BeginTransaction()
	accounts, _ := manager.accountRepo.GetAccountsByUserId(user_id)
	if accounts == nil {
		manager.accountRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Счета пользователя не найдены")
	}
	manager.accountRepo.Db.CommitTransaction()

	return accounts, nil
}

// Аналитика счетов пользователя
func (manager *AccountManager) GetFinancialSummaryByUserId(user_id int64) (map[string]any, error) {
	log.Println("Аналитика счетов пользователя")
	manager.m.Lock()
	defer manager.m.Unlock()

	manager.accountRepo.Db.BeginTransaction()
	accounts, _ := manager.accountRepo.GetAccountsByUserId(user_id)
	if accounts == nil {
		manager.accountRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Счета пользователя не найдены")
	}

	var totalBalance float64
	for _, acc := range *accounts {
		totalBalance += acc.Balance
	}

	summary := map[string]interface{}{
		"user_id":               user_id,
		"total_account_balance": totalBalance,
		"number_of_accounts":    len(*accounts),
	}

	log.Printf("Generated financial summary for user %d", user_id)
	manager.accountRepo.Db.CommitTransaction()

	return summary, nil
}
