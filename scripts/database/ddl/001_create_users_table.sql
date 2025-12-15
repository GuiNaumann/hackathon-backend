-- Tabela de usuários
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

-- Índice para busca por email
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Inserir usuário de teste (senha: admin123)
INSERT INTO users (email, name, password)
VALUES (
           'admin@hackathon.com',
           'Admin User',
           '$2a$10$LKprZ8AGbw.UMI6Flk6OB.ceOXgUeJtoLQzP.4mt4rGvBvGImtbIi'
       ) ON CONFLICT (email) DO NOTHING;