start:
	docker-compose up

q2c:
	sqlc generate

test:
	go test -v -covermode=atomic -race -coverpkg=./... ./... \
	-coverprofile coverage.out.tmp && cat coverage.out.tmp | grep -v "_mock.go" | grep -v "injector.go" | grep -v "_gen.go" > coverage.out && rm coverage.out.tmp && \
	go tool cover -func coverage.out
	
generateInjector:
	wire ./...

migrateInit:
	migrate create -ext sql -dir db/migration -seq $(name)

mockRepo:
	mockgen -package mockrepo -destination db/repository/mock/repository_mock.go -source=db/repository/repository.go -aux_files github.com/gadhittana-01/book-go/db/repository=db/repository/querier.go

mockAuthMiddleware:
	mockgen -package mockutl -source=./utils/auth_middleware.go -destination=./utils/mock/auth_middleware_mock.go

mockToken:
	mockgen -package mockutl -source=./utils/token.go -destination=./utils/mock/token_mock.go

mockUserSvc:
	mockgen -package mocksvc -source=./service/user_service.go -destination=./service/mock/user_service_mock.go

mockBookSvc:
	mockgen -package mocksvc -source=./service/book_service.go -destination=./service/mock/book_service_mock.go

mockOrderSvc:
	mockgen -package mocksvc -source=./service/order_service.go -destination=./service/mock/order_service_mock.go

checkLint:
	golangci-lint run ./... -v

fixLint:
	golangci-lint run --fix

swaggerLocal:
	swag init && rm -rf docs/swagger.json