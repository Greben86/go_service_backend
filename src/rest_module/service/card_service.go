package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"rest_module/repository"
	"sync"

	log "github.com/sirupsen/logrus"

	. "rest_module/model"

	"golang.org/x/crypto/bcrypt"
)

var SECRET_KEY = []byte("1234567890")

type CardManager struct {
	m          sync.Mutex // мьютекс для синхронизации доступа
	mailSender *MailSender
	userRepo   *repository.UserRepository // репозиторий пользователей
	cardRepo   *repository.CardRepository // репозиторий карт
}

// Конструктор сервиса
func CardManagerNewInstance(mailSender *MailSender, userRepo *repository.UserRepository, cardRepo *repository.CardRepository) *CardManager {
	manager := CardManager{}
	manager.mailSender = mailSender
	manager.userRepo = userRepo
	manager.cardRepo = cardRepo
	return &manager
}

// Создание карты
func (manager *CardManager) AddCard(card Card, user_id int64) (*Card, error) {
	log.Println("Создание карты")
	manager.m.Lock()
	defer manager.m.Unlock()

	manager.cardRepo.Db.BeginTransaction()
	user, _ := manager.userRepo.GetUserByID(user_id)
	if user == nil {
		manager.cardRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Пользователь с таким логином не найден")
	}

	exist, err := manager.cardRepo.GetCardByNumber(user_id, card.Number)
	if exist != nil || err != nil {
		manager.cardRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Карта с таким номер уже есть")
	}

	cvv, hashed := GenerateCVV()
	card.UserId = user_id
	card.CVV = hashed
	card.Number = EncodeHMAC(card.Number, SECRET_KEY)
	card.ID, err = manager.cardRepo.InsertCard(&card)
	if err != nil {
		manager.cardRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Ошибка добавления карты %s", err.Error())
	}
	manager.mailSender.SendCVVEmailMessage(user.Email, cvv)
	manager.cardRepo.Db.CommitTransaction()
	return &card, nil
}

// Создание CVV-кода по алгоритму Луна
func GenerateCVV() (string, string) {
	n, _ := rand.Int(rand.Reader, big.NewInt(900))
	cvv := fmt.Sprintf("%03d", n.Int64()+100)
	hashed, _ := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	return cvv, string(hashed)
}

// Генерация HMAC для номера карты
func EncodeHMAC(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// Поиск карты по идентификатору
func (manager *CardManager) FindCardById(user_id, id int64) (*Card, error) {
	log.Println("Поиск карты по идентификатору")
	manager.m.Lock()
	defer manager.m.Unlock()

	manager.cardRepo.Db.BeginTransaction()
	card, _ := manager.cardRepo.GetCardByID(user_id, id)
	if card == nil {
		manager.cardRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Карта с таким идентификатором не найдена")
	}
	manager.cardRepo.Db.CommitTransaction()

	return card, nil
}

// Поиск карты по номеру
func (manager *CardManager) FindCardByNumber(user_id int64, number string) (*Card, error) {
	log.Println("Поиск карты по номеру")
	manager.m.Lock()
	defer manager.m.Unlock()

	manager.cardRepo.Db.BeginTransaction()
	card, _ := manager.cardRepo.GetCardByNumber(user_id, number)
	if card == nil {
		manager.cardRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Карта с таким номером не найдена")
	}
	manager.cardRepo.Db.CommitTransaction()

	return card, nil
}

// Поиск карт пользователя
func (manager *CardManager) FindCardsByUserId(user_id int64) (*[]Card, error) {
	log.Println("Поиск карт пользователя")
	manager.m.Lock()
	defer manager.m.Unlock()

	manager.cardRepo.Db.BeginTransaction()
	cards, _ := manager.cardRepo.GetCardsByUserId(user_id)
	if cards == nil {
		manager.cardRepo.Db.RollbackTransaction()
		return nil, fmt.Errorf("Карты пользователя не найдены")
	}
	manager.cardRepo.Db.CommitTransaction()

	return cards, nil
}
