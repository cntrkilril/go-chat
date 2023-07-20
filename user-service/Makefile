protoc:
	protoc -I. --go-grpc_out=. --go_out=. pkg/proto/$(file).proto

migrate-up:
	migrate -path db/migration -database "postgresql://root:pass@localhost:5432/users?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://root:pass@localhost:5432/users?sslmode=disable" -verbose down


