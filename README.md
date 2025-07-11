# Azugo-Training

> This repository is used solely to practice with azugo golang framework. 

## To run locally: 

1. Run GO API:
```bash
go run ./cmd/server
```

2. Run Docker with Keycloak: 
``` bash 
docker compose up -d
```

3. If any swagger changes are made:
``` bash
 swag init --generalInfo /cmd/server/main.go
```

## General info 

Keycloak will be available at http://localhost:8081, with both the username and password set to admin.
