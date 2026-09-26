#!/bin/bash

# Obtener el directorio donde reside este script de forma robusta
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
ENV_FILE="$ROOT_DIR/.env"

# Cargar archivo .env si existe, exportando las variables automáticamente
if [ -f "$ENV_FILE" ]; then
  set -a
  source "$ENV_FILE"
  set +a
else
  echo "Advertencia: Archivo .env no encontrado en $ROOT_DIR"
fi

# Variables de configuración (ahora tomarán los valores del .env)
CONTAINER="lvbp-postgres"
USER="${POSTGRES_USER//$'\r'/}"
DB="${POSTGRES_DB//$'\r'/}"

case "$1" in
  start)
    echo "Levantando infraestructura de base de datos..."
    docker compose -f "$ROOT_DIR/docker-compose.yml" --env-file "$ENV_FILE" up -d
    ;;
  stop)
    echo "Apagando la base de datos..."
    docker compose -f "$ROOT_DIR/docker-compose.yml" --env-file "$ENV_FILE" down
    ;;
  shell)
    echo "Abriendo la consola de PostgreSQL..."
    docker exec -it "$CONTAINER" psql -U "$USER" -d "$DB"
    ;;
  *)
    echo "Uso correcto: $0 {start|stop|shell}"
    exit 1
    ;;
esac