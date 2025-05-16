```md
# asset-service

Simple Go microservice for uploading/downloading user assets with token-based auth.

## Setup
1. Install Go & Goose.
2. Configure `config.yaml`.
3. Run migrations:
   ```bash
goose -dir migrations postgres "$DATABASE_URL" up
```  
4. Build & run:
   ```bash
go build -o asset-service cmd/asset-service/main.go
./asset-service -config config.yaml
```  

## API
- POST `/api/auth` → `{token}`
- POST `/api/assets/{name}` Bearer → `{status:"ok"}`
- GET `/api/assets/{name}` Bearer → data
- GET `/api/assets` Bearer → `[]string`
- DELETE `/api/assets/{name}` Bearer → `{status:"ok"}`

## Tests
```bash
go test ./internal/service
```