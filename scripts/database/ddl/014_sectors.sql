-- Tabela de setores
CREATE TABLE IF NOT EXISTS sectors (
                                       id BIGSERIAL PRIMARY KEY,
                                       name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

-- Índices
CREATE INDEX IF NOT EXISTS idx_sectors_name ON sectors(name);
CREATE INDEX IF NOT EXISTS idx_sectors_active ON sectors(active);

-- Inserir setores padrão
INSERT INTO sectors (name, description) VALUES
                                            ('Recursos Humanos', 'Departamento de gestão de pessoas e talentos'),
                                            ('Tecnologia da Informação', 'Departamento de TI e infraestrutura'),
                                            ('Financeiro', 'Departamento financeiro e contábil'),
                                            ('Comercial', 'Departamento de vendas e relacionamento com clientes'),
                                            ('Marketing', 'Departamento de marketing e comunicação'),
                                            ('Operações', 'Departamento de operações e logística'),
                                            ('Jurídico', 'Departamento jurídico e compliance'),
                                            ('Produto', 'Departamento de gestão de produtos')
    ON CONFLICT (name) DO NOTHING;

-- Adicionar coluna sector_id na tabela users
ALTER TABLE users ADD COLUMN IF NOT EXISTS sector_id BIGINT REFERENCES sectors(id) ON DELETE SET NULL;

-- Criar índice
CREATE INDEX IF NOT EXISTS idx_users_sector ON users(sector_id);

COMMENT ON TABLE sectors IS 'Setores/Departamentos da empresa';
COMMENT ON COLUMN users.sector_id IS 'Setor ao qual o usuário pertence';

-- Permissões para CRUD de setores (Admin e Manager)
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, endpoint, method
FROM user_type ut
         CROSS JOIN (
    VALUES
        ('/api/private/sectors', 'GET'),
        ('/api/private/sectors', 'POST'),
        ('/api/private/sectors/{id}', 'GET'),
        ('/api/private/sectors/{id}', 'PUT'),
        ('/api/private/sectors/{id}', 'DELETE')
) AS perms(endpoint, method)
WHERE ut.name IN ('admin', 'manager')
    ON CONFLICT DO NOTHING;

-- Permissão para usuários comuns apenas visualizarem
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, endpoint, method
FROM user_type ut
         CROSS JOIN (
    VALUES
        ('/api/private/sectors', 'GET'),
        ('/api/private/sectors/{id}', 'GET')
) AS perms(endpoint, method)
WHERE ut.name = 'user'
    ON CONFLICT DO NOTHING;