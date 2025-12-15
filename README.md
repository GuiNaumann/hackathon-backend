# Hackathon Backend - Clean Architecture

Backend completo em Go seguindo Clean Architecture / Hexagonal, pronto para produção.

## 🚀 Quick Start

### 1. Configurar Banco de Dados

```bash
# Criar banco PostgreSQL
createdb hackathon_db

# Executar migrations
psql -U postgres -d hackathon_db -f scripts/database/ddl/001_create_users_table.sql