FROM golang:1.23.3-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /transaction-system cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /transaction-system /app

EXPOSE 8080

CMD ["/app/transaction-system"]
