# Etapa 1: build
FROM golang:1.23 AS builder

WORKDIR /app

# Copiamos los archivos go.mod y go.sum primero (para aprovechar cache)
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el código
COPY . .

# Compilamos el binario
RUN go build -o main .

# Etapa 2: runtime
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
