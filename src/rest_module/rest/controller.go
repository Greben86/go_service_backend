package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	. "rest_module/domain/dto"
	. "rest_module/service"
)

// API приложения.
type API struct {
	r           *mux.Router  // маршрутизатор запросов
	userManager *UserManager // сервис пользователей
}

// Конструктор API.
func ApiNewInstance(userManager *UserManager) *API {
	api := API{}
	api.userManager = userManager
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

func (api *API) exampleHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.Context().Value("id").(string))
	user, err := api.userManager.FindUserById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ResponseUser{ID: user.ID, Username: user.Username}
	json, _ := json.Marshal(&response)
	w.Write(json)
}

func (api *API) registerHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса с помощью io.ReadAll
	body, err := io.ReadAll(r.Body)

	// Закрываем тело запроса
	defer r.Body.Close()

	// Проверяем наличие ошибок
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Выводим тело запроса в ответ
	request := RequestSignUp{}
	err = json.Unmarshal(body, &request)

	// Проверяем наличие ошибок
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := api.userManager.AddUser(request.Username, request.Password)
	// Проверяем наличие ошибок
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, _ := GenerateJWTToken(fmt.Sprint(user.ID))
	responseDTO := ResponseAuth{JwtToken: token}
	response, _ := json.Marshal(&responseDTO)
	w.Write(response)
}

func (api *API) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса с помощью io.ReadAll
	body, err := io.ReadAll(r.Body)

	// Закрываем тело запроса
	defer r.Body.Close()

	// Проверяем наличие ошибок
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Выводим тело запроса в ответ
	request := RequestSignUp{}
	err = json.Unmarshal([]byte(body), &request)

	// Проверяем наличие ошибок
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := api.userManager.FindUserByName(request.Username)
	// Проверяем наличие ошибок
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = CheckPasswordForUser(user, request.Password)
	// Проверяем наличие ошибок
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, _ := GenerateJWTToken(fmt.Sprint(user.ID))
	responseDTO := ResponseAuth{JwtToken: token}
	response, _ := json.Marshal(&responseDTO)
	w.Write(response)
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	// Предшествующий код
	// Public routes
	api.Router().HandleFunc("/register", api.registerHandler).Methods(http.MethodPost) // регистрация
	api.Router().HandleFunc("/login", api.loginHandler).Methods(http.MethodPost)       // аутентификация
	// Protected routes
	authRouter := api.Router().PathPrefix("/").Subrouter()
	authRouter.Use(AuthMiddleware)

	authRouter.HandleFunc("/accounts", api.exampleHandler).Methods(http.MethodPost)                    // создать счет
	authRouter.HandleFunc("/cards", api.exampleHandler).Methods(http.MethodPost)                       // выпустить карту
	authRouter.HandleFunc("/transfer", api.exampleHandler).Methods(http.MethodPost)                    // перевод средств
	authRouter.HandleFunc("/analytics", api.exampleHandler).Methods(http.MethodGet)                    // получить аналитику
	authRouter.HandleFunc("/credits/{creditId}/schedule", api.exampleHandler).Methods(http.MethodGet)  // график платежей по кредиту
	authRouter.HandleFunc("/accounts/{accountId}/predict", api.exampleHandler).Methods(http.MethodGet) // прогноз баланса
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		id, err := CheckTokenAndGetId(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
