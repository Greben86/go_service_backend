package service

import (
	"fmt"
	"rest_module/repository"
	"sync"

	log "github.com/sirupsen/logrus"

	. "rest_module/model"
)

type OperationManager struct {
	m           sync.Mutex // мьютекс для синхронизации доступа
	mailSender  *MailSender
	userRepo    *repository.UserRepository      // репозиторий пользователей
	accountRepo *repository.AccountRepository   // репозиторий счетов
	operRepo    *repository.OperationRepository // репозиторий операций
}

// Конструктор сервиса
func OperationManagerNewInstance(mailSender *MailSender, userRepo *repository.UserRepository, accountRepo *repository.AccountRepository, operRepo *repository.OperationRepository) *OperationManager {
	manager := OperationManager{}
	manager.mailSender = mailSender
	manager.userRepo = userRepo
	manager.accountRepo = accountRepo
	manager.operRepo = operRepo
	return &manager
}

// Создание операции дебета
func (manager *OperationManager) AddOperationDebet(operation Operation, user_id int64) (*Operation, error) {
	log.Println("Создание операции дебета")
	manager.m.Lock()
	defer manager.m.Unlock()
	var err error

	manager.operRepo.Db.BeginTransaction()
	user, _ := manager.userRepo.GetUserByID(user_id)
	if user == nil {
		manager.operRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Пользователь с таким логином не найден")
	}

	account, _ := manager.accountRepo.GetAccountByIDAndUserID(user_id, operation.AccountId)
	if account == nil {
		manager.operRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Счет не найден")
	}

	account.Balance -= operation.SumValue
	err = manager.accountRepo.UpdateAccount(account)
	if err != nil {
		manager.operRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Ошибка обновления счета %s", err.Error())
	}

	operation.OperationType = "DEBET"
	operation.UserId = user_id
	operation.ID, err = manager.operRepo.InsertOperation(&operation)
	if err != nil {
		manager.operRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Ошибка добавления операции %s", err.Error())
	}
	manager.mailSender.SendEmailMessage(user.Email, operation.SumValue)
	manager.operRepo.Db.CommitTransaction()
	return &operation, nil
}

// Создание операции кредита
func (manager *OperationManager) AddOperationCredit(operation Operation, user_id int64) (*Operation, error) {
	log.Println("Создание операции кредита")
	manager.m.Lock()
	defer manager.m.Unlock()
	var err error

	manager.operRepo.Db.BeginTransaction()
	user, _ := manager.userRepo.GetUserByID(user_id)
	if user == nil {
		manager.operRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Пользователь с таким логином не найден")
	}

	account, _ := manager.accountRepo.GetAccountByIDAndUserID(user_id, operation.AccountId)
	if account == nil {
		manager.operRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Счет не найден")
	}

	account.Balance += operation.SumValue
	err = manager.accountRepo.UpdateAccount(account)
	if err != nil {
		manager.operRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Ошибка обновления счета %s", err.Error())
	}

	operation.OperationType = "CREDIT"
	operation.UserId = user_id
	operation.ID, err = manager.operRepo.InsertOperation(&operation)
	if err != nil {
		manager.operRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Ошибка добавления операции %s", err.Error())
	}
	manager.mailSender.SendEmailMessage(user.Email, operation.SumValue)
	manager.operRepo.Db.CommitTransaction()
	return &operation, nil
}

// Создание операции перевода
func (manager *OperationManager) AddOperationTransfer(sum_value float64, account_from, account_to, user_id int64) error {
	log.Println("Создание операции перевода")
	manager.m.Lock()
	defer manager.m.Unlock()
	var err error

	manager.operRepo.Db.BeginTransaction()
	user, _ := manager.userRepo.GetUserByID(user_id)
	if user == nil {
		manager.operRepo.Db.RollbackTransaction()
		return fmt.Errorf("Пользователь с таким логином не найден")
	}

	accountFrom, _ := manager.accountRepo.GetAccountByIDAndUserID(user_id, account_from)
	if accountFrom == nil {
		manager.operRepo.Db.RollbackTransaction()
		return fmt.Errorf("Счет не найден")
	}

	accountFrom.Balance -= sum_value
	err = manager.accountRepo.UpdateAccount(accountFrom)
	if err != nil {
		manager.operRepo.Db.RollbackTransaction()
		return fmt.Errorf("Ошибка обновления счета %s", err.Error())
	}

	operationFrom := Operation{}
	operationFrom.OperationType = "DEBET"
	operationFrom.SumValue = sum_value
	operationFrom.AccountId = account_from
	operationFrom.UserId = user_id
	operationFrom.ID, err = manager.operRepo.InsertOperation(&operationFrom)
	if err != nil {
		manager.operRepo.Db.RollbackTransaction()
		return fmt.Errorf("Ошибка добавления операции %s", err.Error())
	}

	accountTo, _ := manager.accountRepo.GetAccountByID(user_id)
	if accountTo == nil {
		manager.operRepo.Db.RollbackTransaction()
		return fmt.Errorf("Счет не найден")
	}

	accountTo.Balance += sum_value
	err = manager.accountRepo.UpdateAccount(accountTo)
	if err != nil {
		manager.operRepo.Db.RollbackTransaction()
		return fmt.Errorf("Ошибка обновления счета %s", err.Error())
	}

	operationTo := Operation{}
	operationTo.OperationType = "CREDIT"
	operationTo.SumValue = sum_value
	operationTo.AccountId = account_from
	operationTo.UserId = user_id
	operationTo.ID, err = manager.operRepo.InsertOperation(&operationTo)
	if err != nil {
		manager.operRepo.Db.RollbackTransaction()
		return fmt.Errorf("Ошибка добавления операции %s", err.Error())
	}

	manager.mailSender.SendEmailMessage(user.Email, sum_value)
	manager.operRepo.Db.CommitTransaction()
	return nil
}

// Поиск операций пользователя
func (manager *OperationManager) FindOperationsByUserId(user_id int64) (*[]Operation, error) {
	log.Println("Поиск операций пользователя")
	manager.m.Lock()
	defer manager.m.Unlock()

	manager.operRepo.Db.BeginTransaction()
	cards, _ := manager.operRepo.GetOperationsByUserId(user_id)
	if cards == nil {
		manager.operRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Операции пользователя не найдены")
	}
	manager.operRepo.Db.CommitTransaction()

	return cards, nil
}

// Поиск операций пользователя по счету
func (manager *OperationManager) FindOperationsByAccountId(user_id, account_id int64) (*[]Operation, error) {
	log.Println("Поиск операций пользователя по счету")
	manager.m.Lock()
	defer manager.m.Unlock()

	manager.operRepo.Db.BeginTransaction()
	cards, _ := manager.operRepo.GetOperationsByAccountId(user_id, account_id)
	if cards == nil {
		manager.operRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Операции по счету не найдены")
	}
	manager.operRepo.Db.CommitTransaction()

	return cards, nil
}
