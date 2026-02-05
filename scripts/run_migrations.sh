#!/bin/bash

# Script de migration para o projeto hackathon-backend
# Usa as credenciais definidas no docker-compose.yaml

DB_USER="hackathon_tic"
DB_PASSWORD="hackathon2137689712"
DB_NAME="hackathon"
DB_HOST="localhost"
DB_PORT="5432"

# Diretório das migrations
MIGRATIONS_DIR="scripts/database/ddl"

# Função para executar uma migration
run_migration() {
    local file=$1
    echo "Executando: $file"
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$file"
    if [ $? -eq 0 ]; then
        echo "✓ Sucesso: $file"
    else
        echo "✗ Falha: $file"
        exit 1
    fi
}

# Verificar se o diretório existe
if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo "Erro: Diretório de migrations não encontrado: $MIGRATIONS_DIR"
    exit 1
fi

# Verificar se existem arquivos SQL
if [ ! "$(ls -A $MIGRATIONS_DIR/*.sql 2>/dev/null)" ]; then
    echo "Erro: Nenhum arquivo SQL encontrado em $MIGRATIONS_DIR"
    exit 1
fi

echo "Iniciando migrations..."
echo "Database: $DB_NAME"
echo "Host: $DB_HOST:$DB_PORT"
echo ""

# Executar migrations em ordem numérica
for file in $(ls -1 "$MIGRATIONS_DIR"/*.sql | sort -V); do
    run_migration "$file"
done

echo ""
echo "Migrations concluídas com sucesso!"
