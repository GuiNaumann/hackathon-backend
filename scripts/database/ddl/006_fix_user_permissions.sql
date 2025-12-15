-- Limpar permissões antigas que estavam erradas
DELETE FROM user_type_permissions WHERE endpoint LIKE '/api/private/admin/users%';

-- Permissões corretas para CRUD de usuários (ADMIN)
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, endpoint, method
FROM user_type ut
         CROSS JOIN (
    VALUES
        -- CRUD de usuários
        ('/api/private/users', 'GET'),
        ('/api/private/users', 'POST'),
        ('/api/private/users/{id}', 'GET'),
        ('/api/private/users/{id}', 'PUT'),
        ('/api/private/users/{id}', 'DELETE'),

        -- Gerenciar tipos de usuários
        ('/api/private/admin/users/{userId}/types/{typeId}', 'POST'),
        ('/api/private/admin/users/{userId}/types/{typeId}', 'DELETE'),

        -- Listar tipos disponíveis
        ('/api/private/user-types', 'GET'),

        -- Informações pessoais
        ('/api/private/me', 'GET'),
        ('/api/private/personal-information', 'GET'),
        ('/api/private/change-password', 'POST'),

        -- Iniciativas (acesso total)
        ('/api/private/initiatives', 'GET'),
        ('/api/private/initiatives', 'POST'),
        ('/api/private/initiatives/{id}', 'GET'),
        ('/api/private/initiatives/{id}', 'PUT'),
        ('/api/private/initiatives/{id}', 'DELETE'),
        ('/api/private/initiatives/{id}/status', 'PATCH'),
        ('/api/private/my-initiatives', 'GET')
) AS perms(endpoint, method)
WHERE ut.name = 'admin'
    ON CONFLICT DO NOTHING;

-- Permissões para MANAGER
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, endpoint, method
FROM user_type ut
         CROSS JOIN (
    VALUES
        -- Apenas visualizar usuários
        ('/api/private/users', 'GET'),
        ('/api/private/users/{id}', 'GET'),

        -- Listar tipos disponíveis
        ('/api/private/user-types', 'GET'),

        -- Informações pessoais
        ('/api/private/me', 'GET'),
        ('/api/private/personal-information', 'GET'),
        ('/api/private/change-password', 'POST'),

        -- Iniciativas (pode criar e ver)
        ('/api/private/initiatives', 'GET'),
        ('/api/private/initiatives', 'POST'),
        ('/api/private/initiatives/{id}', 'GET'),
        ('/api/private/my-initiatives', 'GET')
) AS perms(endpoint, method)
WHERE ut.name = 'manager'
    ON CONFLICT DO NOTHING;

-- Permissões para USER
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, endpoint, method
FROM user_type ut
         CROSS JOIN (
    VALUES
        -- Informações pessoais
        ('/api/private/me', 'GET'),
        ('/api/private/personal-information', 'GET'),
        ('/api/private/change-password', 'POST'),

        -- Iniciativas (pode criar e ver próprias)
        ('/api/private/initiatives', 'GET'),
        ('/api/private/initiatives', 'POST'),
        ('/api/private/initiatives/{id}', 'GET'),
        ('/api/private/my-initiatives', 'GET')
) AS perms(endpoint, method)
WHERE ut.name = 'user'
    ON CONFLICT DO NOTHING;