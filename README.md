<div align="center">

<img src="docs/demo.gif" width="100%" alt="Demo">

# OAuth 2.0 + Auth0 SSO

Single Sign-On на Go с аутентификацией через Auth0 и Google.

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

```bash
git clone https://github.com/qqwozz/OAuth_2.0.git
cd OAuth_2.0
cp .env.example .env
# Заполните .env своими данными Auth0
go run main.go
```

Откройте **http://localhost:8080** → нажмите **Sign In with Google**

---

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

---

## Configuration

| Variable | Required | Description |
|----------|:--------:|-------------|
| `AUTH0_DOMAIN` | Yes | Auth0 domain |
| `AUTH0_CLIENT_ID` | Yes | Auth0 client ID |
| `AUTH0_CLIENT_SECRET` | Yes | Auth0 client secret |
| `AUTH0_REDIRECT_URL` | Yes | `http://localhost:8080/callback` |
| `PORT` | No | `8080` |

---

## Routes

| Method | Path | Auth | Description |
|:------:|------|:----:|-------------|
| `GET` | `/` | - | Home |
| `GET` | `/login` | - | Redirect to Auth0 |
| `GET` | `/callback` | - | OAuth callback |
| `GET` | `/profile` | Yes | User profile |
| `GET` | `/avatar` | Yes | Avatar proxy |
| `GET` | `/logout` | - | Logout |

---

## License

[MIT](LICENSE) © [Dima Kiselev](https://github.com/qqwozz)
