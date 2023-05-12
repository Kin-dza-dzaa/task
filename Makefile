run:
	go run ./cmd/app/ 1 > logs.log

swagger:
	swag fmt
	swag init --parseDependency -g internal/transport/http/v1/rest/currency.go

migrateup:
	migrate -database postgresql://testuser:12345@localhost:5432/test?sslmode=disable -path ./internal/repository/postgresql/migrations/ up 

migratedown:
	migrate -database postgresql://testuser:12345@localhost:5432/test?sslmode=disable -path ./internal/repository/postgresql/migrations/ down 

test:
	go test ./...

# For tests you will need docker API on port 2375 
# with disabled tls
cover:
	go test -coverprofile cover.out ./...
	go tool cover -html cover.out