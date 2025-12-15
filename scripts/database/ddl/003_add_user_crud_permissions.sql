-- Adicionar permissões de CRUD de usuários para admin
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, endpoint, method
FROM user_type ut
         CROSS JOIN (
    VALUES
        ('/api/private/users', 'GET'),
        ('/api/private/users', 'POST'),
        ('/api/private/users/{id}', 'GET'),
        ('/api/private/users/{id}', 'PUT'),
        ('/api/private/users/{id}', 'DELETE'),
        ('/api/private/change-password', 'POST'),
        ('/api/private/me', 'GET'),
        ('/api/private/personal-information', 'GET'),
        ('/api/private/admin/users/{userId}/types/{typeId}', 'POST'),
        ('/api/private/admin/users/{userId}/types/{typeId}', 'DELETE')
) AS perms(endpoint, method)
WHERE ut.name = 'admin'
ON CONFLICT DO NOTHING;

-- Adicionar permissão de trocar senha para todos os tipos
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, endpoint, method
FROM user_type ut
         CROSS JOIN (
    VALUES
        ('/api/private/change-password', 'POST'),
        ('/api/private/me', 'GET'),
        ('/api/private/personal-information', 'GET')
) AS perms(endpoint, method)
WHERE ut.name IN ('manager', 'user')
ON CONFLICT DO NOTHING;