
# Rate Limiter em Go

## Configuração e Execução

### Passo 1: Configurar Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:

```
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
DEFAULT_IP_LIMIT=10
DEFAULT_TOKEN_LIMIT=10
BLOCK_DURATION=5s
```

### Passo 2: Docker Compose

Certifique-se de ter o Docker e Docker Compose instalados. Crie um arquivo `docker-compose.yml` na raiz do projeto para configurar o Redis:

```yaml
version: '3.8'
services:
  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  redis_data:
```

Para iniciar o Redis, execute:

```sh
docker-compose up -d
```

### Passo 3: Executar o Projeto

Para rodar o projeto, use os seguintes comandos:

```sh
go mod tidy
go run cmd/main.go
```

### Estrutura do Projeto

A estrutura do projeto deve ser semelhante a esta:

```
├── cmd
│   └── main.go
├── docker-compose.yml
├── .env
├── go.mod
├── go.sum
└── pkg
    ├── middleware.go
    ├── rate_limiter.go
    └── rate_test.go
```
