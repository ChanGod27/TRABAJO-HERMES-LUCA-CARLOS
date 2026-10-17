#!/bin/bash
set -e

# Función de limpieza mejorada: apaga la DB y mata el servidor Go
cleanup() {
    echo "Limpiando servidor y contenedores..."
    if [ -n "$SERVER_PID" ]; then
        kill "$SERVER_PID" 2>/dev/null || true
    fi
    docker compose down -v
}

echo "1. Limpiando contenedores y volúmenes previos..."
docker compose down -v --remove-orphans

echo "2. Generando código con sqlc..."
sqlc generate

echo "3. Levantando base de datos..."
docker compose up -d

# Registramos la nueva función de limpieza
trap cleanup EXIT

echo "4. Esperando a que la base de datos esté lista (Healthcheck)..."
until [ "$(docker inspect -f '{{.State.Health.Status}}' $(docker compose ps -q db) 2>/dev/null)" = "healthy" ]; do
    sleep 1
done

echo "5. Ejecutando tests unitarios/integración de Go..."
go test -v ./...

echo "6. Levantando el servidor HTTP en segundo plano..."
go run . &
SERVER_PID=$! # Guardamos el ID del proceso para matarlo en el trap

echo "7. Esperando a que el servidor HTTP responda en :8080..."
until curl -s http://localhost:8080/ > /dev/null 2>&1 \
    || curl -s http://localhost:8080/usuarios > /dev/null 2>&1; do
    sleep 1
done

echo "8. Ejecutando tests de API con Hurl..."
hurl --test request.hurl

echo "Pruebas finalizadas con éxito."