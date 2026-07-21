<div align="center">

<img src="docs/demo.gif" width="100%" alt="OAuth 2.0 + Auth0 SSO Demo">

# OAuth 2.0 + Auth0 SSO

**Single Sign-On приложение на Go с аутентификацией через Auth0 и Google**

<br/>

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Auth0](https://img.shields.io/badge/Auth0-EB5424?style=for-the-badge&logo=auth0&logoColor=white)](https://auth0.com/)
[![Gin](https://img.shields.io/badge/Gin-0089CF?style=for-the-badge&logo=gin&logoColor=white)](https://gin-gonic.com/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

<br/>

[Features](#features) • [Quick Start](#quick-start) • [How It Works](#how-it-works) • [Routes](#routes)

</div>

---

## Features

<div align="center">

| | |
|:---:|---|
| **Google SSO** | Вход через Google-аккаунт одним кликом |
| **OIDC-верификация** | ID-токен проверяется через Auth0 provider |
| **CSRF-защита** | Рандомный state-параметр при каждом входе |
| **Серверные сессии** | Claims хранятся в зашифрованной cookie |
| **Secure Logout** | Очистка сессии на стороне Auth0 |

</div>

---

## Tech Stack

<div align="center">

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-0089CF?style=flat-square&logo=gin&logoColor=white)
![Auth0](https://img.shields.io/badge/Auth0-EB5424?style=flat-square&logo=auth0&logoColor=white)
![OAuth2](https://img.shields.io/badge/OAuth2-FF6B35?style=flat-square)
![OIDC](https://img.shields.io/badge/OIDC-4A90D9?style=flat-square)
![TailwindCSS](https://img.shields.io/badge/Tailwind_CSS-06B6D4?style=flat-square&logo=tailwindcss&logoColor=white)

</div>

---

## Quick Start

### 1. Клонируйте

```bash
git clone https://github.com/qqwozz/OAuth_2.0.git
cd OAuth_2.0
```

### 2. Настройте Auth0

1. Зарегистрируйтесь на [auth0.com](https://auth0.com)
2. Создайте **Regular Web Application**
3. Скопируйте из настроек приложения:

| Параметр | Где найти |
|----------|-----------|
| `AUTH0_DOMAIN` | Applications → Your App → Domain |
| `AUTH0_CLIENT_ID` | Applications → Your App → Client ID |
| `AUTH0_CLIENT_SECRET` | Applications → Your App → Client Secret |

4. Добавьте в **Allowed Callback URLs**:
   ```
   http://localhost:8080/callback
   ```

5. Добавьте в **Allowed Logout URLs**:
   ```
   http://localhost:8080/
   ```

### 3. Конфигурация

```bash
cp .env.example .env
```

Заполните `.env`:

```env
AUTH0_DOMAIN=dev-xxxxxxxxxx.us.auth0.com
AUTH0_CLIENT_ID=your_client_id
AUTH0_CLIENT_SECRET=your_client_secret
AUTH0_REDIRECT_URL=http://localhost:8080/callback
SESSION_SECRET=generate_random_string_min_32_chars
```

### 4. Запуск

```bash
go run main.go
```

Откройте **http://localhost:8080** → нажмите **Sign In with Google**

---

## How It Works

<div align="center">

```
┌──────────┐            ┌──────────┐            ┌──────────┐
│          │  1. Login  │          │            │          │
│          │───────────▶│          │  2. Auth   │          │
│          │            │          │───────────▶│          │
│ Browser  │            │  Server  │            │  Auth0   │
│          │  4. User   │   (Go)   │  3. Tokens │          │
│          │◀───────────│          │◀───────────│          │
│          │            │          │            │          │
│          │  5. Profile│  6. Verify            │          │
│          │───────────▶│  ID Token│───────────▶│          │
│          │◀───────────│          │◀───────────│          │
└──────────┘            └──────────┘            └──────────┘
```

</div>

| Step | Description |
|:----:|-------------|
| 1 | User clicks "Sign In with Google" → redirect to Auth0 |
| 2 | Auth0 authenticates user via Google |
| 3 | Auth0 returns authorization code to callback URL |
| 4 | Server exchanges code for access token + ID token |
| 5 | Server verifies ID token via OIDC provider |
| 6 | User profile rendered with claims from verified token |

---

## Routes

| Method | Path | Auth | Description |
|:------:|------|:----:|-------------|
| `GET` | `/` | - | Главная страница |
| `GET` | `/login` | - | Редирект на Auth0 |
| `GET` | `/callback` | - | OAuth callback handler |
| `GET` | `/profile` | **Yes** | Профиль пользователя |
| `GET` | `/logout` | - | Выход из системы |
| `GET` | `/ping` | - | Health check |

---

## Project Structure

```
OAuth_2.0/
├── main.go                # Server, handlers, middleware
├── go.mod                 # Go module
├── go.sum                 # Dependencies checksum
├── .env.example           # Environment template
├── LICENSE                # MIT License
├── README.md
├── docs/
│   ├── demo.gif           # Demo GIF
│   └── demo.mp4           # Demo Video
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

---

## License

<div align="center">

[MIT](LICENSE) © [Dima Kiselev](https://github.com/qqwozz)

</div>
