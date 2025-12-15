-- Tabela de tipos de usuário
CREATE TABLE IF NOT EXISTS user_type (
                                         id BIGSERIAL PRIMARY KEY,
                                         name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

-- Tabela de relacionamento usuário <-> tipo
CREATE TABLE IF NOT EXISTS type_user (
                                         id BIGSERIAL PRIMARY KEY,
                                         user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_type_id BIGINT NOT NULL REFERENCES user_type(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, user_type_id)
    );

-- Tabela de permissões por endpoint
CREATE TABLE IF NOT EXISTS user_type_permissions (
                                                     id BIGSERIAL PRIMARY KEY,
                                                     user_type_id BIGINT NOT NULL REFERENCES user_type(id) ON DELETE CASCADE,
    endpoint VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_type_id, endpoint, method)
    );

-- Índices
CREATE INDEX IF NOT EXISTS idx_type_user_user_id ON type_user(user_id);
CREATE INDEX IF NOT EXISTS idx_type_user_type_id ON type_user(user_type_id);
CREATE INDEX IF NOT EXISTS idx_permissions_type_endpoint ON user_type_permissions(user_type_id, endpoint, method);

-- Inserir tipos de usuário padrão
INSERT INTO user_type (name, description) VALUES
                                              ('admin', 'Administrador com acesso total'),
                                              ('manager', 'Gerente com acesso intermediário'),
                                              ('user', 'Usuário comum com acesso básico')
    ON CONFLICT (name) DO NOTHING;

-- Atribuir tipo admin ao usuário de teste
INSERT INTO type_user (user_id, user_type_id)
SELECT u.id, ut.id
FROM users u, user_type ut
WHERE u.email = 'admin@hackathon.com' AND ut. name = 'admin'
    ON CONFLICT DO NOTHING;

-- Configurar permissões para admin (acesso total aos exemplos)
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, endpoint, method
FROM user_type ut
         CROSS JOIN (
    VALUES
        ('/api/private/me', 'GET'),
        ('/api/private/personal-information', 'GET'),
        ('/api/private/admin/users', 'GET'),
        ('/api/private/admin/users', 'POST'),
        ('/api/private/admin/users', 'DELETE')
) AS perms(endpoint, method)
WHERE ut.name = 'admin'
    ON CONFLICT DO NOTHING;

-- Configurar permissões para manager (acesso intermediário)
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, endpoint, method
FROM user_type ut
         CROSS JOIN (
    VALUES
        ('/api/private/me', 'GET'),
        ('/api/private/personal-information', 'GET'),
        ('/api/private/admin/users', 'GET')
) AS perms(endpoint, method)
WHERE ut.name = 'manager'
    ON CONFLICT DO NOTHING;

-- Configurar permissões para user (acesso básico)
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, endpoint, method
FROM user_type ut
         CROSS JOIN (
    VALUES
        ('/api/private/me', 'GET'),
        ('/api/private/personal-information', 'GET')
) AS perms(endpoint, method)
WHERE ut.name = 'user'
    ON CONFLICT DO NOTHING;