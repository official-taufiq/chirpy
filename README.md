
> ⚠️ **Note:** This project was built by following a tutorial from [Boot.dev](https://www.boot.dev/).  
> It was created for learning purposes — to understand how HTTP servers, routing, databases, and authentication work in Go.

---

# 🐦 Chirpy — A Lightweight Social Network API (Twitter Clone)

Chirpy is a backend web application built in **Go (Golang)** that mimics the core functionality of a social media platform like Twitter — where users can post short messages (“chirps”), manage their accounts, and interact with others securely via JWT-based authentication.

This project is designed as a **learning and showcase project** to demonstrate REST API design, database interaction, authentication, and web server management in Go.

---

## 🚀 Features

- 🐤 **Chirps (Posts)** — Create, fetch, update, and delete chirps  
- 👤 **User Authentication** — Register, login, and update user profiles  
- 🔐 **JWT Authentication** — Secure endpoints using access and refresh tokens  
- 🔁 **Token Refresh & Revocation** — Supports refresh token rotation and revocation for session security  
- 🧠 **Polka Webhooks Integration** — Handles asynchronous external service events  
- 📊 **Admin Metrics** — Tracks file server hits and server statistics  
- 🧹 **Admin Reset Endpoint** — Clear metrics or reset app data for testing/development  
- 💾 **PostgreSQL Database** — Persistent storage with migrations managed by Goose  
- ⚙️ **Environment-based Configuration** — Uses `.env` for environment variables  
- 🧱 **Modular Code Structure** — Clean separation of handlers, middleware, and database logic  

---

## 🛠️ Tech Stack

| Component | Description |
|------------|-------------|
| **Language** | Go (Golang) |
| **Database** | PostgreSQL |
| **Environment Management** | `godotenv` |
| **Database Driver** | `lib/pq` |
| **ORM / Query Builder** | Custom queries using `sqlc` (via `internal/database`) |
| **Authentication** | JWT Tokens (Access & Refresh) |
| **HTTP Router** | Standard `net/http` multiplexer |
| **Package Management** | Go Modules (`go.mod`) |

---

## 📁 Project Structure

```
chirpy/
├── main.go
├── go.mod
├── go.sum
├── .env
├── internal/
│   ├── auth/          # JWT handling, token utilities
│   ├── database/      # Generated SQL queries and database access layer
│   ├── handlers/      # API route handlers (user, chirps, metrics, etc.)
│   ├── middleware/    # Reusable middlewares (metrics, auth, etc.)
│   └── utils/         # Helper utilities
└── migrations/        # Goose migrations for database schema
```

---

## ⚙️ Setup & Installation

### 1️⃣ Clone the repository
```bash
git clone https://github.com/official-taufiq/chirpy.git
cd chirpy
```

### 2️⃣ Install dependencies
```bash
go mod tidy
```

### 3️⃣ Setup your environment variables
Create a `.env` file in the root directory:
```bash
DB_URL=postgres://username:password@localhost:5432/chirpy?sslmode=disable
JWT_SECRET=your_jwt_secret_key
POLKA_KEY=your_polka_api_key
PLATFORM=dev
```

### 4️⃣ Setup the database
Make sure PostgreSQL is running, then run migrations (using Goose or your preferred migration tool):
```bash
goose up
```

### 5️⃣ Run the application
```bash
go run main.go
```

The server will start on **http://localhost:8080**

---

## 🌐 API Endpoints Overview

| Method | Endpoint | Description |
|---------|-----------|-------------|
| `GET` | `/api/healthz` | Health check |
| `POST` | `/api/users` | Create new user |
| `POST` | `/api/login` | Login and receive JWT |
| `PUT` | `/api/users` | Update user info |
| `POST` | `/api/refresh` | Refresh access token |
| `POST` | `/api/revoke` | Revoke refresh token |
| `POST` | `/api/chirps` | Create a chirp |
| `GET` | `/api/chirps` | Get all chirps |
| `POST` | `/api/chirps/{chirpID}` | Get one chirp |
| `DELETE` | `/api/chirps/{chirpID}` | Delete a chirp |
| `POST` | `/api/polka/webhooks` | Handle Polka webhooks |
| `GET` | `/admin/metrics` | View metrics |
| `POST` | `/admin/reset` | Reset app metrics |

---

## 🧪 Example Request

### Create a New Chirp
```bash
curl -X POST http://localhost:8080/api/chirps   -H "Content-Type: application/json"   -H "Authorization: Bearer <access_token>"   -d '{"body": "Hello world! This is my first chirp."}'
```

---

## 🔒 Authentication Flow

1. **User Login** → `/api/login`  
   Returns an **access token** (short-lived) and **refresh token** (long-lived).  
2. **Access Protected Routes** → Include the access token in the `Authorization` header.  
3. **Token Expiry** → Use `/api/refresh` with the refresh token to get a new access token.  
4. **Logout / Revoke** → Call `/api/revoke` to invalidate a refresh token.

---

## 🧰 Development Notes

- This project follows Go’s standard project layout and idiomatic conventions.  
- The handlers and database access are structured around dependency injection via `apiConfig`.  
- Built using **Go 1.21+** (recommended).  
- Code formatted with `go fmt`.  

---

## Motivation
## quickstart
## usage
## contributing
