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

// API карт
type CardController struct {
	cardManager *CardManager // сервис карт
}

// Конструктор API карт
func CardControllerNewInstance(cardManager *CardManager) *CardController {
	api := CardController{}
	api.cardManager = cardManager
	return &api
}

// Endpoint для добавления карты
func (api *CardController) AddCardHandler(w http.ResponseWriter, r *http.Request) {
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
	request := Card{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user_id, _ := strconv.Atoi(r.Context().Value("id").(string))
	card, err := api.cardManager.AddCard(request, int64(user_id))
	// Проверяем наличие ошибок
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(&card)
	w.Write(response)
}

// Endpoint информации карты
func (api *CardController) CardInfoHandler(w http.ResponseWriter, r *http.Request) {
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

	card, err := api.cardManager.FindCardById(int64(user_id), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(&card)
	w.Write(response)
}

// Endpoint списка карт пользователя
func (api *CardController) CardListHandler(w http.ResponseWriter, r *http.Request) {
	// Считывание параметра из контекста
	id, err := strconv.Atoi(r.Context().Value("id").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cards, err := api.cardManager.FindCardsByUserId(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(&cards)
	w.Write(response)
}
