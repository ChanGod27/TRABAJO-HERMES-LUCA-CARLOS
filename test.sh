#!/bin/bash
set -e

echo "1. Limpiando contenedores y volúmenes previos..."
docker compose down -v --remove-orphans

echo "2. Generando código con sqlc..."
sqlc generate

echo "3. Levantando base de datos..."
docker compose up -d

# Registrado aquí: si falla la espera o los tests, limpia igual
trap 'docker compose down -v' EXIT

echo "4. Esperando a que la base de datos esté lista (Healthcheck)..."
until [ "$(docker inspect -f '{{.State.Health.Status}}' $(docker compose ps -q db) 2>/dev/null)" = "healthy" ]; do
    sleep 1
done

echo "5. Ejecutando tests..."
go test -v ./...

echo "Pruebas finalizadas con éxito."