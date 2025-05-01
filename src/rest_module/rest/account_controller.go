package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	. "rest_module/domain/dto"
	. "rest_module/service"

	"github.com/gorilla/mux"
)

// API счетов
type AccountController struct {
	accountManager *AccountManager // сервис счетов
}

// Конструктор API счетов
func AccountControllerNewInstance(accountManager *AccountManager) *AccountController {
	api := AccountController{}
	api.accountManager = accountManager
	return &api
}

// Endpoint для регистрации
func (api *AccountController) AddAccountHandler(w http.ResponseWriter, r *http.Request) {
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
	request := RequestAccount{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user_id, _ := strconv.Atoi(r.Context().Value("id").(string))
	account, err := api.accountManager.AddAccount(request.Name, request.Bank, int64(user_id))
	// Проверяем наличие ошибок
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseDTO := ResponseAccount{ID: account.ID, Name: account.Name, Bank: account.Bank}
	response, _ := json.Marshal(&responseDTO)
	w.Write(response)
}

// Endpoint информации о счете
func (api *AccountController) AccountInfoHandler(w http.ResponseWriter, r *http.Request) {
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

	account, err := api.accountManager.FindAccountById(int64(user_id), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseDTO := ResponseAccount{ID: account.ID, Name: account.Name, Bank: account.Bank}
	response, _ := json.Marshal(&responseDTO)
	w.Write(response)
}

// Endpoint списка счетов пользователя
func (api *AccountController) AccountListHandler(w http.ResponseWriter, r *http.Request) {
	// Считывание параметра из контекста
	id, err := strconv.Atoi(r.Context().Value("id").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accounts, err := api.accountManager.FindAccountsByUserId(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var responseDTOs []ResponseAccount
	for _, account := range *accounts {
		dto := ResponseAccount{ID: account.ID, Name: account.Name, Bank: account.Bank}
		responseDTOs = append(responseDTOs, dto)
	}
	response, _ := json.Marshal(&responseDTOs)
	w.Write(response)
}
