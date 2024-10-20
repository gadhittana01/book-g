# Build stage
FROM golang:1.23-alpine3.19 AS builder

WORKDIR /app
COPY ./ ./

# Install Git
RUN apt-get update && apt-get install -y git
RUN go mod tidy
RUN go build -o main .

# Run stage
FROM alpine:3.19

WORKDIR /app
COPY --from=builder ./app/main ./
COPY ./config/app.env ./config/app.env
COPY ./db/migration ./db/migration

CMD ["/app/main"]