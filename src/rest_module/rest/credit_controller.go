package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	. "rest_module/model"
	. "rest_module/service"

	"github.com/gorilla/mux"
)

// API кредитов
type CreditController struct {
	creditManager *CreditManager // сервис карт
}

// Конструктор API кредитов
func CreditControllerNewInstance(creditManager *CreditManager) *CreditController {
	api := CreditController{}
	api.creditManager = creditManager
	return &api
}

// Endpoint для добавления кредита
func (api *CreditController) AddCreditHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса с помощью io.ReadAll
	body, err := io.ReadAll(r.Body)

	// Закрываем тело запроса
	defer r.Body.Close()

	// Проверяем наличие ошибок
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Выводим тело запроса в ответ
	request := Credit{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if request.Amount <= 0 {
		http.Error(w, "Сумма кридита должна быть больше нуля", http.StatusBadRequest)
		return
	}

	user_id, _ := strconv.Atoi(r.Context().Value("id").(string))
	credit, err := api.creditManager.AddCredit(request, int64(user_id))
	// Проверяем наличие ошибок
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(&credit)
	w.Write(response)
}

// Endpoint информации о кредит
func (api *CreditController) CreditInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Считывание параметра из контекста
	user_id, err := strconv.Atoi(r.Context().Value("id").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Считывание параметра {id} из пути запроса.
	requestParam := mux.Vars(r)["id"]
	var id int
	id, err = strconv.Atoi(requestParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	credit, err := api.creditManager.FindCreditById(int64(user_id), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(&credit)
	w.Write(response)
}

// Endpoint списка кредитов пользователя
func (api *CreditController) CreditListHandler(w http.ResponseWriter, r *http.Request) {
	// Считывание параметра из контекста
	user_id, err := strconv.Atoi(r.Context().Value("id").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cards, err := api.creditManager.FindCreditsByUserId(int64(user_id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(&cards)
	w.Write(response)
}

// Endpoint графика платежей по кредиту
func (api *CreditController) PaymentScheduleHandler(w http.ResponseWriter, r *http.Request) {
	// Считывание параметра из контекста
	user_id, err := strconv.Atoi(r.Context().Value("id").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Считывание параметра {id} из пути запроса.
	requestParam := mux.Vars(r)["id"]
	var credit_id int
	credit_id, err = strconv.Atoi(requestParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cards, err := api.creditManager.PaymentScheduleByCreditId(int64(user_id), int64(credit_id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(&cards)
	w.Write(response)
}

// Endpoint прогноза баланса
func (api *CreditController) AccountPredictHandler(w http.ResponseWriter, r *http.Request) {
	// Считывание параметра из контекста
	user_id, err := strconv.Atoi(r.Context().Value("id").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Считывание параметра {id} из пути запроса.
	requestParam := mux.Vars(r)["id"]
	var credit_id int
	credit_id, err = strconv.Atoi(requestParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cards, err := api.creditManager.AccountPredictByCreditId(int64(user_id), int64(credit_id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(&cards)
	w.Write(response)
}
