package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob" // Регистрация типов для сериализации в сессиях
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// UserInfo - структура для хранения информации о пользователе из Auth0
type UserInfo struct {
	Sub           string    `json:"sub"`
	GivenName     string    `json:"given_name"`
	Nickname      string    `json:"nickname"`
	Name          string    `json:"name"`
	Picture       string    `json:"picture"`
	UpdatedAt     int64     `json:"updated_at"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
}

// server - основная структура сервера
type server struct {
	router       *gin.Engine
	oauth2Config *oauth2.Config
}

// NewServer - создает новый экземпляр сервера
func NewServer() (*server, error) {
	router := gin.Default() // используем Default для логирования и recovery

	oauth2Config, err := NewOauth2Config()
	if err != nil {
		return nil, fmt.Errorf("could not create new oauth2 config: %v", err)
	}

	return &server{
		router:       router,
		oauth2Config: oauth2Config,
	}, nil
}

// NewOauth2Config - создает конфигурацию OAuth2 для Auth0
func NewOauth2Config() (*oauth2.Config, error) {
	// Получаем переменные окружения
	domain := os.Getenv("AUTH0_DOMAIN")
	if domain == "" {
		return nil, fmt.Errorf("AUTH0_DOMAIN environment variable is not set")
	}

	clientID := os.Getenv("AUTH0_CLIENT_ID")
	if clientID == "" {
		return nil, fmt.Errorf("AUTH0_CLIENT_ID environment variable is not set")
	}

	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")
	if clientSecret == "" {
		return nil, fmt.Errorf("AUTH0_CLIENT_SECRET environment variable is not set")
	}

	redirectURL := os.Getenv("AUTH0_REDIRECT_URL")
	if redirectURL == "" {
		return nil, fmt.Errorf("AUTH0_REDIRECT_URL environment variable is not set")
	}

	// Создаем OIDC провайдера
	providerURL := fmt.Sprintf("https://%s/", domain)
	log.Printf("Creating OIDC provider with URL: %s", providerURL)

	// Таймаут для соединения с Auth0
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		return nil, fmt.Errorf("could not create new provider: %v", err)
	}

	log.Println("OIDC provider created successfully")

	// Возвращаем конфигурацию OAuth2
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "picture"},
		Endpoint:     provider.Endpoint(),
	}, nil
}

// loginHandler - обрабатывает вход пользователя
func (s *server) loginHandler(ctx *gin.Context) {
	log.Println("Login handler called")

	// Генерируем случайную строку state для защиты от CSRF
	state, err := generateRandomString()
	if err != nil {
		log.Printf("Error generating random string: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate state"})
		return
	}

	// Сохраняем state в сессии
	session := sessions.Default(ctx)
	session.Set("state", state)

	if err := session.Save(); err != nil {
		log.Printf("Error saving session: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not save session"})
		return
	}

	// Перенаправляем на Auth0 для аутентификации
	authURL := s.oauth2Config.AuthCodeURL(state)
	log.Printf("Redirecting to Auth0: %s", authURL)
	ctx.Redirect(http.StatusTemporaryRedirect, authURL)
}

// logoutHandler - обрабатывает выход пользователя
func (s *server) logoutHandler(ctx *gin.Context) {
	// Удаляем все cookies и значения сессий
	// устанавливаем отрицательное время для cookie
	ctx.SetCookie("at", "", -1, "/", "", false, true)
	ctx.SetCookie("u", "", -1, "/", "", false, true)
	ctx.SetCookie("auth-sessions", "", -1, "/", "", false, true)

	// Вызываем endpoint выхода Auth0 для очистки сессии и токена на стороне Auth0
	logoutURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not logout"})
		return
	}

	// Проверяем, был ли запрос выполнен через http или https
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	// Перенаправляем пользователя обратно на главную страницу
	redirectUrl, err := url.Parse(scheme + "://" + ctx.Request.Host)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not parse URL"})
		return
	}

	// Добавляем параметры URL
	params := url.Values{}
	params.Add("returnTo", redirectUrl.String())
	params.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutURL.RawQuery = params.Encode()

	ctx.Redirect(http.StatusTemporaryRedirect, logoutURL.String())
}

// callbackHandler - обрабатывает callback от Auth0
func (s *server) callbackHandler(ctx *gin.Context) {
	log.Println("Callback handler called")

	// Проверяем, совпадает ли state в сессии и в запросе
	session := sessions.Default(ctx)
	state := session.Get("state")
	queryState := ctx.Query("state")

	log.Printf("Session state: %v, Query state: %v", state, queryState)

	if state != queryState {
		log.Printf("State mismatch: session=%v, query=%v", state, queryState)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid state param"})
		return
	}

	// Удаляем state из сессии после проверки
	session.Delete("state")
	if err := session.Save(); err != nil {
		log.Printf("Error saving session after deleting state: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not save session"})
		return
	}

	// Получаем code из запроса
	code := ctx.Query("code")
	if code == "" {
		log.Println("Code parameter is missing")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing code parameter"})
		return
	}

	// Обмениваем code на токен
	log.Println("Exchanging code for token...")
	token, err := s.oauth2Config.Exchange(ctx, code)
	if err != nil {
		log.Printf("Error exchanging code: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not exchange code"})
		return
	}

	// Проверяем валидность токена
	if !token.Valid() {
		log.Println("Token is invalid or expired")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		return
	}

	// Получаем информацию о пользователе для отображения в профиле
	client := s.oauth2Config.Client(ctx, token)
	resp, err := client.Get("https://" + os.Getenv("AUTH0_DOMAIN") + "/userinfo")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch user information"})
		return
	}
	defer resp.Body.Close()

	// Парсим тело ответа
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not parse response body"})
		return
	}

	// Сохраняем access token и тело ответа в cookie
	// (Примечание: cookie должен быть зашифрован перед сохранением)
	// u -> userInfo
	ctx.SetCookie("u", string(b), int(time.Now().Add(1*time.Hour).Unix()), "/", "", false, true)
	// at -> access token
	ctx.SetCookie("at", token.AccessToken, int(time.Now().Add(1*time.Hour).Unix()), "/", "", false, true)

	// Извлекаем ID токен из ответа
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Println("No id_token in response")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no id_token in response"})
		return
	}

	// Создаем провайдера для верификации токена
	providerURL := fmt.Sprintf("https://%s/", os.Getenv("AUTH0_DOMAIN"))
	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		log.Printf("Error creating provider for verification: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create provider"})
		return
	}

	// Верифицируем ID токен
	verifier := provider.Verifier(&oidc.Config{ClientID: os.Getenv("AUTH0_CLIENT_ID")})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Printf("Error verifying id token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not verify id token"})
		return
	}

	// Извлекаем данные пользователя из токена
	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		log.Printf("Error parsing claims: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not parse claims"})
		return
	}

	// Сохраняем пользователя в сессии
	session.Set("user", claims)

	if err := session.Save(); err != nil {
		log.Printf("Error saving user session: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not save user session"})
		return
	}

	log.Printf("User authenticated successfully: %v", claims["email"])
	ctx.Redirect(http.StatusTemporaryRedirect, "/profile")
}

func main() {
	// Регистрируем тип map для gob, чтобы можно было сохранять в сессии
	gob.Register(map[string]interface{}{})

	// Проверяем наличие всех необходимых переменных окружения
	requiredVars := []string{"AUTH0_DOMAIN", "AUTH0_CLIENT_ID", "AUTH0_CLIENT_SECRET", "AUTH0_REDIRECT_URL"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("Environment variable %s is not set", v)
		}
	}

	// Создаем сервер
	server, err := NewServer()
	if err != nil {
		log.Fatalf("could not create new server: %v", err)
	}

	// Настройка сессий
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "superSecretValue"
		log.Println("WARNING: Using default session secret. Set SESSION_SECRET environment variable in production!")
	}

	// Убеждаемся, что секрет достаточно длинный
	for len(secret) < 32 {
		secret += "0"
	}

	// Создаем хранилище сессий
	store := cookie.NewStore([]byte(secret))
	server.router.Use(sessions.Sessions("auth-sessions", store))

	// Статические файлы и шаблоны
	server.router.Static("/public", "web/static")
	server.router.LoadHTMLGlob("web/template/*")

	// Определяем маршруты
	server.router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Дебаг-эндпоинт для проверки переменных окружения
	server.router.GET("/debug/env", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"AUTH0_DOMAIN":        os.Getenv("AUTH0_DOMAIN") != "",
			"AUTH0_CLIENT_ID":     os.Getenv("AUTH0_CLIENT_ID") != "",
			"AUTH0_CLIENT_SECRET": os.Getenv("AUTH0_CLIENT_SECRET") != "",
			"AUTH0_REDIRECT_URL":  os.Getenv("AUTH0_REDIRECT_URL"),
		})
	})

	// Главная страница
	server.router.GET("/", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		user := session.Get("user")
		ctx.HTML(http.StatusOK, "home.html", gin.H{
			"loggedIn": user != nil,
			"user":     user,
		})
	})

	// Страница профиля (требует аутентификации)
	server.router.GET("/profile", IsAuthenticate(), func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		user := session.Get("user")
		if user == nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		// Получаем информацию о пользователе из cookie
		userInfo, err := ctx.Cookie("u")
		if err != nil {
			// Если cookie с информацией о пользователе не существует, перенаправляем на главную
			ctx.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		// Десериализуем значение в структуру UserInfo
		var u UserInfo
		if err := json.Unmarshal([]byte(userInfo), &u); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something wrong. Please try logging in again"})
			return
		}

		// Отображаем страницу профиля с данными пользователя
		ctx.HTML(http.StatusOK, "profile.html", gin.H{
			"Profile": u,
		})
	})

	// Выход из системы
	server.router.GET("/logout", server.logoutHandler)

	// Маршруты аутентификации
	server.router.GET("/login", server.loginHandler)
	server.router.GET("/callback", server.callbackHandler)

	// Запускаем сервер
	log.Println("Server starting on :8080")
	if err := server.router.Run(":8080"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

// IsAuthenticate - middleware для проверки аутентификации пользователя
// Проверяет наличие валидного access token в cookie
func IsAuthenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := ctx.Cookie("at")
		if err != nil || accessToken == "" {
			// Cookie не существует или пустой, перенаправляем на главную
			ctx.Redirect(http.StatusTemporaryRedirect, "/")
			ctx.Abort()
			return
		}
		// Если все хорошо, передаем запрос дальше
		ctx.Next()
	}
}

// generateRandomString - генерирует случайную строку для state
func generateRandomString() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("crypto/rand failed: %v", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}