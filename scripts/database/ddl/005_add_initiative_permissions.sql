-- Permissões para initiatives
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, endpoint, method
FROM user_type ut
         CROSS JOIN (
    VALUES
        -- Admin: acesso total
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

-- Manager e User:  podem criar e ver suas próprias
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, endpoint, method
FROM user_type ut
         CROSS JOIN (
    VALUES
        ('/api/private/initiatives', 'GET'),
        ('/api/private/initiatives', 'POST'),
        ('/api/private/initiatives/{id}', 'GET'),
        ('/api/private/my-initiatives', 'GET')
) AS perms(endpoint, method)
WHERE ut.name IN ('manager', 'user')
    ON CONFLICT DO NOTHING;