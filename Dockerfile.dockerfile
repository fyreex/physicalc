# ── Estágio 1: Build ──────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder
 
WORKDIR /app
 
# Copia os arquivos de dependências primeiro (cache eficiente)
COPY go.mod go.sum ./
RUN go mod download
 
# Copia o restante do código
COPY . .
 
# Compila o binário estático (sem dependências externas)
RUN CGO_ENABLED=0 GOOS=linux go build -o physicalc ./cmd/server
 
# ── Estágio 2: Imagem final mínima ────────────────────────────────────────────
FROM alpine:3.19
 
WORKDIR /app
 
# Certificados SSL (necessários para HTTPS)
RUN apk --no-cache add ca-certificates
 
# Copia apenas o binário compilado
COPY --from=builder /app/physicalc .
 
EXPOSE 8080
 
CMD ["./physicalc"]