<div align="center">

# OAuth 2.0 + Auth0 SSO

**Single Sign-On приложение на Go с аутентификацией через Auth0 и Google**

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Auth0](https://img.shields.io/badge/Auth0-EB5424?style=for-the-badge&logo=auth0&logoColor=white)](https://auth0.com/)
[![Gin](https://img.shields.io/badge/Gin-0089CF?style=for-the-badge&logo=gin&logoColor=white)](https://gin-gonic.com/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

<br/>

[Features](#features) | [Quick Start](#quick-start) | [How It Works](#how-it-works) | [Routes](#routes)

</div>

---

## Demo

<div align="center">
  <img src="demogif.gif" width="750" alt="OAuth 2.0 + Auth0 Demo" style="border-radius: 12px; box-shadow: 0 4px 24px rgba(0,0,0,0.15);">
</div>

<br/>

---

## Features

- **Google SSO** — вход через Google-аккаунт одним кликом
- **OIDC-верификация** — ID-токен проверяется через Auth0 provider
- **CSRF-защита** — рандомный state-параметр при каждом входе
- **Серверные сессии** — claims хранятся в зашифрованной cookie
- **Secure logout** — очистка сессии на стороне Auth0

## Tech Stack

<table>
<tr>
<td><b>Backend</b></td>
<td>Go 1.26 + Gin</td>
</tr>
<tr>
<td><b>Auth Provider</b></td>
<td>Auth0 (Google, Email/Password, Social)</td>
</tr>
<tr>
<td><b>Libraries</b></td>
<td>go-oidc, oauth2, gin-sessions, godotenv</td>
</tr>
<tr>
<td><b>Frontend</b></td>
<td>HTML Templates + Tailwind CSS</td>
</tr>
</table>

---

## Quick Start

### 1. Clone

```bash
git clone https://github.com/qqwozz/OAuth_2.0.git
cd OAuth_2.0
```

### 2. Setup Auth0

1. Зарегистрируйтесь на [auth0.com](https://auth0.com)
2. Создайте **Regular Web Application**
3. Скопируйте из настроек приложения:
   - **Domain** → `AUTH0_DOMAIN`
   - **Client ID** → `AUTH0_CLIENT_ID`
   - **Client Secret** → `AUTH0_CLIENT_SECRET`
4. Добавьте в **Allowed Callback URLs**:
   ```
   http://localhost:8080/callback
   ```
5. Добавьте в **Allowed Logout URLs**:
   ```
   http://localhost:8080/
   ```

### 3. Configure

```bash
cp .env.example .env
```

Заполните `.env`:

```env
AUTH0_DOMAIN=dev-xxxxxxxxxx.us.auth0.com
AUTH0_CLIENT_ID=your_client_id
AUTH0_CLIENT_SECRET=your_client_secret
AUTH0_REDIRECT_URL=http://localhost:8080/callback
SESSION_SECRET=generate_a_random_string_min_32_chars
```

### 4. Run

```bash
go run main.go
```

Откройте **http://localhost:8080** и нажмите "Sign In with Google"

---

## How It Works

```
 ┌──────────┐              ┌──────────┐              ┌──────────┐
 │  Browser  │────────────▶│  Server  │────────────▶│  Auth0   │
 │           │◀────────────│  (Go)    │◀────────────│          │
 └──────────┘              └──────────┘              └──────────┘
      │                          │                          │
      │   1. GET /login          │                          │
      │─────────────────────────▶│                          │
      │                          │                          │
      │   2. 307 → Auth0         │                          │
      │◀─────────────────────────│  /authorize              │
      │                          │─────────────────────────▶│
      │                          │                          │
      │   3. User logs in via    │                          │
      │      Google              │                          │
      │───────────────────────────────────────────────────▶│
      │                          │                          │
      │   4. GET /callback       │                          │
      │      ?code=xxx&state=yyy │                          │
      │─────────────────────────▶│                          │
      │                          │   5. Exchange code       │
      │                          │      for tokens          │
      │                          │─────────────────────────▶│
      │                          │◀─────────────────────────│
      │                          │                          │
      │                          │   6. Verify ID token     │
      │                          │      via OIDC provider   │
      │                          │   (cached well-known)    │
      │                          │                          │
      │   7. 307 → /profile      │                          │
      │◀─────────────────────────│                          │
      │                          │                          │
      │   8. GET /profile        │                          │
      │─────────────────────────▶│                          │
      │◀──── 200 + HTML ─────────│                          │
```

---

## Routes

| Method | Path | Description |
|:------:|------|-------------|
| `GET` | `/` | Главная страница с кнопкой входа |
| `GET` | `/login` | Редирект на Auth0 |
| `GET` | `/callback` | Обработка OAuth callback |
| `GET` | `/profile` | Профиль пользователя (🔒) |
| `GET` | `/logout` | Выход из системы |
| `GET` | `/ping` | Health check |

---

## Project Structure

```
.
├── main.go               # Server, handlers, middleware
├── go.mod / go.sum       # Dependencies
├── .env.example          # Environment template
├── LICENSE               # MIT
├── demogif.gif           # Demo
└── web/
    ├── static/img/       # Assets
    └── template/         # HTML templates
```

---

## License

[MIT](LICENSE) &copy; [Dima Kiselev](https://github.com/qqwozz)

<div align="center">
  <sub>Made with Go + Auth0</sub>
</div>
