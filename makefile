start:
	docker-compose up

migrateInit:
	migrate create -ext sql -dir db/migration -seq $(name)

mockRepo:
	mockgen -package mockrepo -destination db/repository/mock/repository_mock.go -source=db/repository/repository.go -aux_files github.com/gadhittana-01/book-go/db/repository=db/repository/querier.go

checkLint:
	golangci-lint run ./... -v