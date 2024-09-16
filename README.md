# Book Go


## Introduction
This project is a online book store Go-based application, order. It uses various tools like Docker, `migrate`, `mockgen`, and `golangci-lint` to ensure smooth development, testing, and deployment workflows.

## Prerequisites to contribute

If you want to contribute to this project, ensure you have the following dependencies installed:

- [Go](https://golang.org/dl/) (v1.XX or higher)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Migrate](https://github.com/golang-migrate/migrate) (for database migrations)
- [Mockgen](https://github.com/golang/mock) (for generating mocks)
- [Wire](https://github.com/google/wire) (dependency injection)
- [Golangci-lint](https://golangci-lint.run/usage/install/) (for linting)

Install these dependencies using the appropriate installation guides.

## Commands
1. Run the app
```bash
make start
```

2. Create a new database migration
```bash
make migrateInit name="your migration name"
```

2. Create a new database migration
```bash
make migrateInit
```

3. Run the tests with race conditions and code coverage
```bash
make test
```

4. Generate repository mock
```bash
make mockRepo
```

5. Run the linter
```bash
make checkLint
```

## Postman documentation
https://www.postman.com/gadhittana/development/collection/0arthyu/book-go