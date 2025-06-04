# Rate Limiter em Go com Redis

Este projeto implementa um middleware de **Rate Limiting** em Go, com suporte a limites por **IP** e **token de acesso**, usando **Redis como armazenamento**.

## Como executar o projeto

### 1. Clonar o repositório
```bash
git clone https://github.com/seu-usuario/rate-limiter-go.git
cd rate-limiter-go
```

### 2. Configurar variáveis de ambiente

Crie um arquivo .env na raiz com os seguintes campos:
```bash
RATE_BACKEND=redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
RATE_LIMIT_IP=5
RATE_LIMIT_TOKEN=10
BLOCK_DURATION_SECONDS=300
```

### 3. Subir o Servidor com Docker

Para subir o servidor completo (Redis + aplicação Go)
```bash
docker-compose up -d
```

### 4. Fazer uma requisição de teste
```bash
curl -H "API_KEY: abc123" http://localhost:8080
```
Você pode modificar ou remover o header API_KEY para testar os diferentes modos de limitação.

### 5. Rodar os testes
```bash
go test ./... -v
```

Mais informações

Para detalhes completos sobre:
  - Funcionamento do rate limiter
  - Explicação dos testes
  - Estratégias de armazenamento
  - Configurações detalhadas

Consulte a documentação técnica
