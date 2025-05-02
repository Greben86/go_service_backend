package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	. "rest_module/domain/dto"
	. "rest_module/service"
)

// API приложения.
type API struct {
	r                 *mux.Router        // маршрутизатор запросов
	usersController   *UsersController   // контроллер пользователей
	accountController *AccountController // контроллер счетов
	cardController    *CardController    // контроллер карт
}

// Конструктор API.
func ApiNewInstance(usersController *UsersController, accountController *AccountController, cardController *CardController) *API {
	api := API{}
	api.usersController = usersController
	api.accountController = accountController
	api.cardController = cardController
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// Endpoint для проверки работы сервиса
func (api *API) healthHandler(w http.ResponseWriter, r *http.Request) {
	response := ResponseHealth{Status: "UP"}
	json, _ := json.Marshal(&response)
	w.Write(json)
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	// Public routes
	api.Router().HandleFunc("/health", api.healthHandler).Methods(http.MethodGet)                      // проверка
	api.Router().HandleFunc("/register", api.usersController.RegisterHandler).Methods(http.MethodPost) // регистрация
	api.Router().HandleFunc("/login", api.usersController.LoginHandler).Methods(http.MethodPost)       // аутентификация
	// Protected routes
	authRouter := api.Router().PathPrefix("/").Subrouter()
	authRouter.Use(AuthMiddleware)

	// Счета
	authRouter.HandleFunc("/accounts/add", api.accountController.AddAccountHandler).Methods(http.MethodPost)      // создать счет
	authRouter.HandleFunc("/accounts/{id}/get", api.accountController.AccountInfoHandler).Methods(http.MethodGet) // получить счет
	authRouter.HandleFunc("/accounts/all", api.accountController.AccountListHandler).Methods(http.MethodGet)      // получить список счетов
	// Карты
	authRouter.HandleFunc("/cards/add", api.cardController.AddCardHandler).Methods(http.MethodPost)      // выпустить карту
	authRouter.HandleFunc("/cards/{id}/get", api.cardController.CardInfoHandler).Methods(http.MethodGet) // выпустить карту
	authRouter.HandleFunc("/cards/all", api.cardController.CardListHandler).Methods(http.MethodGet)      // выпустить карту

	authRouter.HandleFunc("/transfer", api.usersController.UserInfoHandler).Methods(http.MethodPost)                    // перевод средств
	authRouter.HandleFunc("/analytics", api.usersController.UserInfoHandler).Methods(http.MethodGet)                    // получить аналитику
	authRouter.HandleFunc("/credits/{creditId}/schedule", api.usersController.UserInfoHandler).Methods(http.MethodGet)  // график платежей по кредиту
	authRouter.HandleFunc("/accounts/{accountId}/predict", api.usersController.UserInfoHandler).Methods(http.MethodGet) // прогноз баланса
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

// Проверка токена и добавление идентификатора пользователя в контекст
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
