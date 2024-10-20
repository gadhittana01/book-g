# Build stage
FROM golang:1.23-alpine3.19 AS builder

WORKDIR /app
COPY ./ ./

ARG git_token
ENV git_token=$git_token

RUN apk update && apk --no-cache add git

RUN git config \
    --global \
    url."https://gadhittana01:${git_token}@github.com".insteadOf \
    "https://github.com"
RUN go env -w GOPRIVATE=github.com/gadhittana01

RUN go mod tidy
RUN go build -o main .

# Run stage
FROM alpine:3.19

WORKDIR /app
COPY --from=builder ./app/main ./
COPY ./config/app.env ./config/app.env
COPY ./db/migration ./db/migration

CMD ["/app/main"]