# Этап сборки
FROM golang:1.24.1-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /transaction-system ./cmd/main.go

# Этап запуска
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache postgresql-client
COPY --from=build /transaction-system /app
COPY wait-for-postgres.sh .
RUN chmod +x wait-for-postgres.sh
EXPOSE 8080
CMD ["sh", "-c", "./wait-for-postgres.sh db 5432 -- /app/transaction-system"]
