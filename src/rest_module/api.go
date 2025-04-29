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
	"golang.org/x/crypto/bcrypt"
)

var JWT_SECRET = "53A73E5F1C4E0A2D3B5F2D784E6A1B423D6F247D1F6E5C3A596D635A75327855"

// API приложения.
type API struct {
	r  *mux.Router // маршрутизатор запросов
	db *DBManager  // база данных
}

// Конструктор API.
func ApiNewInstance(db *DBManager) *API {
	api := API{}
	api.db = db
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

func (api *API) exampleHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.Context().Value("id").(string))
	user := api.db.GetUserByID(id)
	if user == nil {
		w.Write([]byte("Not found"))
		return
	}

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

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	user := User{Username: request.Username, Password: string(hashedPassword)}
	user.ID = api.db.InsertUser(&user)

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

	user := api.db.GetUserByName(request.Username)
	if user == nil {
		w.Write([]byte("Not found"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Password is wrong"))
		return
	}

	var token, _ = GenerateJWTToken(fmt.Sprint(user.ID))
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
		ctx := context.WithValue(r.Context(), "id", claims.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GenerateJWTToken(id string) (string, error) {
	claims := jwt.RegisteredClaims{
		ID:        id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWT_SECRET))
}
