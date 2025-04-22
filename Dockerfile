FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copiar o código fonte
COPY . .

# Instalar o Swagger e gerar documentação
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go -o ./docs

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o cep-geo-api ./cmd/main.go

FROM alpine:3.19

WORKDIR /app

# Copiar o binário compilado
COPY --from=builder /app/cep-geo-api .

# Expor a porta
EXPOSE 8070

# Definir variável de ambiente
ENV PORT=8070

# Comando para iniciar a aplicação
CMD ["./cep-geo-api"]
