# asset-service

Simple Go microservice for uploading and downloading user assets with token-based authentication.

---

## ⚙️ Setup

1. **Install Go & Goose:**
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

2. **Configure `config.yaml`:**  
   (Пример конфигурации)
   ```yaml
   database_url: "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
   server:
    address: ":8443"
    tls_cert_file: "server.crt"
    tls_key_file: "server.key"
   ```

3. **Run migrations:**
   ```bash
   goose -dir migrations postgres "$DATABASE_URL" up
   ```

4. **Build & run:**
   ```bash
   go build -o asset-service cmd/asset-service/main.go
   ./asset-service -config config.yaml
   ```

---

## 📦 API

### 🔐 Auth

- **POST** `/api/auth`  
  **Body:**
  ```json
  {
    "login": "alice",
    "password": "secret"
  }
  ```
  **Response:**
  ```json
  {
    "token": "abc123"
  }
  ```

---

### 📁 Asset Operations

- **POST** `/api/assets/{name}`  
  Upload asset  
  **Header**: `Authorization: Bearer {token}`  
  **Body**: binary  
  **Response:**
  ```json
  { "status": "ok" }
  ```

- **GET** `/api/assets/{name}`  
  Download asset  
  **Header**: `Authorization: Bearer {token}`  
  **Response**: binary

- **GET** `/api/assets`  
  List user assets  
  **Header**: `Authorization: Bearer {token}`  
  **Response:**
  ```json
  [ "file1.txt", "report.pdf" ]
  ```

- **DELETE** `/api/assets/{name}`  
  Delete asset  
  **Header**: `Authorization: Bearer {token}`  
  **Response:**
  ```json
  { "status": "ok" }
  ```

---

## ✅ Tests

Run unit tests:

```bash
go test ./internal/...
```

---

