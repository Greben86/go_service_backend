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

type ResponseTransfer struct {
	AccountFrom int64   `json:"account_from"`
	AccountTo   int64   `json:"account_to"`
	SumValue    float64 `json:"sum_value"`
}

// API операций
type OperationController struct {
	operManager *OperationManager // сервис операцион
}

// Конструктор API карт
func OperationControllerNewInstance(operManager *OperationManager) *OperationController {
	api := OperationController{}
	api.operManager = operManager
	return &api
}

// Endpoint для добавления операции кредита
func (api *OperationController) AddOperationCreditHandler(w http.ResponseWriter, r *http.Request) {
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
	request := Operation{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if request.SumValue <= 0 {
		http.Error(w, "Сумма операции должна быть больше нуля", http.StatusBadRequest)
		return
	}

	user_id, _ := strconv.Atoi(r.Context().Value("id").(string))
	operation, err := api.operManager.AddOperationCredit(request, int64(user_id))
	// Проверяем наличие ошибок
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(&operation)
	w.Write(response)
}

// Endpoint для добавления операции дебета
func (api *OperationController) AddOperationDebetHandler(w http.ResponseWriter, r *http.Request) {
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
	request := Operation{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user_id, _ := strconv.Atoi(r.Context().Value("id").(string))
	operation, err := api.operManager.AddOperationDebet(request, int64(user_id))
	// Проверяем наличие ошибок
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(&operation)
	w.Write(response)
}

// Endpoint для добавления операции перевода
func (api *OperationController) AddOperationTransferHandler(w http.ResponseWriter, r *http.Request) {
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
	request := ResponseTransfer{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user_id, _ := strconv.Atoi(r.Context().Value("id").(string))
	err = api.operManager.AddOperationTransfer(request.SumValue, request.AccountFrom, request.AccountTo, int64(user_id))
	// Проверяем наличие ошибок
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(&request)
	w.Write(response)
}

// Endpoint списка операций пользователя
func (api *OperationController) OperationListHandler(w http.ResponseWriter, r *http.Request) {
	// Считывание параметра из контекста
	id, err := strconv.Atoi(r.Context().Value("id").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cards, err := api.operManager.FindOperationsByUserId(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(&cards)
	w.Write(response)
}

// Endpoint списка операций пользователя по счету
func (api *OperationController) AccountOperationListHandler(w http.ResponseWriter, r *http.Request) {
	// Считывание параметра из контекста
	user_id, err := strconv.Atoi(r.Context().Value("id").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Считывание параметра {id} из пути запроса.
	requestParam := mux.Vars(r)["id"]
	var account_id int
	account_id, err = strconv.Atoi(requestParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	card, err := api.operManager.FindOperationsByAccountId(int64(user_id), int64(account_id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(&card)
	w.Write(response)
}
