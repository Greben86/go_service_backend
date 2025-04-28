package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

var JWT_SECRET = "53A73E5F1C4E0A2D3B5F2D784E6A1B423D6F247D1F6E5C3A596D635A75327855"

// API приложения.
type API struct {
	r  *mux.Router // маршрутизатор запросов
	db *DB         // база данных
}

// Конструктор API.
func ApiNewInstance(db *DB) *API {
	api := API{}
	api.db = db
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

func (api *API) exampleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Authorization: " + r.Header.Values("Authorization")[0])
	fmt.Println("ID: " + fmt.Sprintf("%#v", r.Context().Value("id")))
	id, _ := strconv.Atoi(fmt.Sprintf("%#v", r.Context().Value("id")))
	user := api.db.GetUser(id)
	fmt.Println(id)

	w.Write([]byte(user.Username))
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

	user := User{Username: request.Username, Password: request.Password}
	user.ID = api.db.NewUser(user)

	var token, _ = GenerateJWTToken(fmt.Sprint(user.ID))

	fmt.Println(request)
	var responseDTO = ResponseAuth{JwtToken: token}

	response, _ := json.Marshal(&responseDTO)
	w.Write(response)

	fmt.Println("register")
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

	var users []User = api.db.Users()
	var user User
	for i := 0; i < len(users); i++ {
		element := users[i]
		if element.Username == request.Username && element.Password == request.Password {
			user = element
			break
		}
	}

	var token, _ = GenerateJWTToken(string(user.ID))
	var responseDTO = ResponseAuth{JwtToken: token}

	response, _ := json.Marshal(&responseDTO)
	w.Write(response)

	fmt.Println("login")
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
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(JWT_SECRET), nil
			})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "id", claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GenerateJWTToken(id string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWT_SECRET))
}
