# Build stage
FROM golang:1.22-alpine3.19 AS builder

WORKDIR /app
COPY ./ ./

RUN go mod tidy
RUN go build -o main .

# Run stage
FROM alpine:3.19

WORKDIR /app
COPY --from=builder ./app/main ./
COPY ./config/app.env ./config/app.env
COPY ./db/migration ./db/migration

CMD ["/app/main"]