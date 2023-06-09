FROM golang:1.20-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o discord-proxy ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/discord-proxy .

EXPOSE 8000

CMD ["./discord-proxy"]