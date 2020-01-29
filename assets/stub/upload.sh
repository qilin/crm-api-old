curl -X POST -H "Content-Type: application/json" localhost:8082/internal/v1/games -d@games.json
curl -X POST -H "Content-Type: application/json" localhost:8082/internal/v1/modules -d@mod_breaker.json
curl -X POST -H "Content-Type: application/json" localhost:8082/internal/v1/modules -d@mod_games.json
