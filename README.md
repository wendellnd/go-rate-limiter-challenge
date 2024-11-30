# go-rate-limiter-challenge

### Execução

Os arquivos `.env` e `.env.production` contém as variáveis de ambiente para configurar o rate limiter.

`.env` é utilizado para execução local
`.env.production` é utilizado no `docker compose up`

```
TOKEN_MAX_REQUESTS=100
TOKEN_EXPIRATION_IN_SECONDS=60
IP_MAX_REQUESTS=10
IP_EXPIRATION_IN_SECONDS=60
```

TOKEN_MAX_REQUESTS - Configura o máximo de requisições por token no rate limiter. O token pode ser definido no header `API_KEY`.
TOKEN_EXPIRATION_IN_SECONDS - Configura o tempo de expiração do rate limiter para os tokens.

IP_MAX_REQUESTS - Configura o máximo de requisições por ip no rate limiter.
IP_EXPIRATION_IN_SECONDS - Configura o tempo de expiração do rate limiter para os ips.

O teste do rate limiter foi feito no arquivo [./middleware/rate_limit_test.go](./middleware/rate_limit_test.go).
