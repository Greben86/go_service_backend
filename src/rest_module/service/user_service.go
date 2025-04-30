package service

import (
	"fmt"
	"rest_module/repository"
	"sync"

	"golang.org/x/crypto/bcrypt"

	. "rest_module/domain/model"
)

type UserManager struct {
	m          sync.Mutex                 // мьютекс для синхронизации доступа
	repository *repository.UserRepository // репозиторий пользователей
}

// Конструктор сервиса
func UserManagerNewInstance(repository *repository.UserRepository) *UserManager {
	manager := UserManager{}
	manager.repository = repository
	return &manager
}

// Создание пользователя
func (manager *UserManager) AddUser(Username, Password string) (*User, error) {
	manager.m.Lock()
	defer manager.m.Unlock()

	exist := manager.repository.GetUserByName(Username)
	if exist != nil {
		return nil, fmt.Errorf("Пользователь с таким логином уже есть")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	user := User{Username: Username, Password: string(hashedPassword)}
	user.ID = manager.repository.InsertUser(&user)
	return &user, nil
}

// Поиск пользователя по идентификатору
func (manager *UserManager) FindUserById(id int) (*User, error) {
	manager.m.Lock()
	defer manager.m.Unlock()

	user := manager.repository.GetUserByID(id)
	if user == nil {
		return nil, fmt.Errorf("Пользователь с таким идентификатором не найден")
	}

	return user, nil
}

// Поиск пользователя по имени
func (manager *UserManager) FindUserByName(Username string) (*User, error) {
	manager.m.Lock()
	defer manager.m.Unlock()

	user := manager.repository.GetUserByName(Username)
	if user == nil {
		return nil, fmt.Errorf("Пользователь с таким логином не найден")
	}

	return user, nil
}
