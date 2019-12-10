curl -X POST -H "Content-Type: application/json" localhost:8082/v1/internal/games -d@games.json
curl -X POST -H "Content-Type: application/json" localhost:8082/v1/internal/modules -d@mod_breaker.json
curl -X POST -H "Content-Type: application/json" localhost:8082/v1/internal/modules -d@mod_games.json
