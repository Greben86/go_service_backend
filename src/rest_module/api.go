package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

var JWT_SECRET = "53A73E5F1C4E0A2D3B5F2D784E6A1B423D6F247D1F6E5C3A596D635A75327855"

// API приложения.
type API struct {
	r *mux.Router // маршрутизатор запросов
	// db *db.DB     // база данных
}

// Конструктор API.
func ApiNewInstance() *API {
	api := API{}
	// api.db = db
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

func (api *API) exampleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("example")
	w.Write([]byte("example"))
}

func (api *API) registerHandler(w http.ResponseWriter, r *http.Request) {
	var token, _ = GenerateJWTToken("1")
	fmt.Println("register")
	w.Write([]byte("token: " + token))
}

func (api *API) loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("login")
	w.Write([]byte("login"))
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	// Предшествующий код
	// Public routes
	api.r.HandleFunc("/sign/up", api.registerHandler).Methods(http.MethodPost)
	api.r.HandleFunc("/sign/in", api.loginHandler).Methods(http.MethodPost)
	// Protected routes
	authRouter := api.r.PathPrefix("/").Subrouter()
	authRouter.Use(AuthMiddleware)

	authRouter.HandleFunc("/example", api.exampleHandler).Methods(http.MethodGet)
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
