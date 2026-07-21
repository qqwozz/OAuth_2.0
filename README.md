<h1 align="center">OAuth 2.0 + Auth0</h1>

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/Auth0-EB5424?style=for-the-badge&logo=auth0&logoColor=white" alt="Auth0">
  <img src="https://img.shields.io/badge/Gin-0089CF?style=for-the-badge&logo=gin&logoColor=white" alt="Gin">
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
</p>

<p align="center">
  Single Sign-On (SSO) приложение на Go с аутентификацией через <b>Auth0</b> и Google.  
  Реализует полный флоу OAuth 2.0 / OIDC: login → callback → session → profile.
</p>

---

## Demo

<p align="center">
  <img src="demogif.gif" width="700" alt="Demo">
</p>

## Stack

| Технология | Назначение |
|---|---|
| **Go + Gin** | HTTP-сервер и роутинг |
| **Auth0** | Identity Provider (Google SSO) |
| **go-oidc** | OIDC-провайдер и верификация ID-токена |
| **oauth2** | Обмен authorization code на токены |
| **gin-sessions** | Серверные сессии (cookie store) |
| **godotenv** | Загрузка `.env` |

## Структура проекта

```
.
├── main.go              # Сервер, хендлеры, middleware
├── go.mod
├── go.sum
├── .env.example         # Шаблон переменных окружения
├── LICENSE
├── demogif.gif          # Demo
└── web/
    ├── static/
    │   └── img/
    │       ├── google.png
    │       └── password.png
    └── template/
        ├── header.html
        ├── footer.html
        ├── home.html
        └── profile.html
```

## Как запустить

### 1. Клонируйте репозиторий

```bash
git clone https://github.com/<your-username>/OAuth_2.0.git
cd OAuth_2.0
```

### 2. Настройте Auth0

1. Создайте аккаунт на [auth0.com](https://auth0.com)
2. Создайте **Regular Web Application**
3. В настройках приложения скопируйте:
   - **Domain** → `AUTH0_DOMAIN`
   - **Client ID** → `AUTH0_CLIENT_ID`
   - **Client Secret** → `AUTH0_CLIENT_SECRET`
4. В **Allowed Callback URLs** добавьте: `http://localhost:8080/callback`
5. В **Allowed Logout URLs** добавьте: `http://localhost:8080/`

### 3. Создайте `.env`

```bash
cp .env.example .env
```

Заполните:

```env
AUTH0_DOMAIN=dev-xxxxxxxxxx.us.auth0.com
AUTH0_CLIENT_ID=your_client_id
AUTH0_CLIENT_SECRET=your_client_secret
AUTH0_REDIRECT_URL=http://localhost:8080/callback
SESSION_SECRET=your_random_secret_at_least_32_chars
```

### 4. Запустите

```bash
go run main.go
```

Откройте [http://localhost:8080](http://localhost:8080)

## Флоу аутентификации

```
┌──────────┐         ┌──────────┐         ┌──────────┐
│  Browser  │────────▶│  Server  │────────▶│  Auth0   │
│           │◀────────│  (Go)    │◀────────│          │
└──────────┘         └──────────┘         └──────────┘
     │                     │                     │
     │  1. GET /login      │                     │
     │────────────────────▶│                     │
     │                     │  2. redirect to     │
     │◀─── 307 ────────────│  Auth0 /authorize   │
     │                     │────────────────────▶│
     │                     │                     │
     │  3. User logs in    │                     │
     │─────────────────────────────────────────▶│
     │                     │                     │
     │  4. GET /callback?code=...                │
     │────────────────────▶│                     │
     │                     │  5. Exchange code   │
     │                     │  for tokens         │
     │                     │────────────────────▶│
     │                     │◀────────────────────│
     │                     │                     │
     │  6. 307 → /profile  │                     │
     │◀─── 307 ────────────│                     │
     │                     │                     │
     │  7. GET /profile    │                     │
     │────────────────────▶│                     │
     │◀──── 200 + HTML ────│                     │
```

## Маршруты

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/` | Главная страница (кнопка входа) |
| GET | `/login` | Редирект на Auth0 для входа |
| GET | `/callback` | Обработка ответа от Auth0 |
| GET | `/profile` | Профиль пользователя (требует авторизации) |
| GET | `/logout` | Выход (очистка сессии + Auth0 logout) |
| GET | `/ping` | Health check |

## Лицензия

[MIT](LICENSE) &copy; Dima Kiselev
