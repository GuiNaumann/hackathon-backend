-- Permissões para comentários (todos os usuários autenticados podem comentar)
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, endpoint, method
FROM user_type ut
         CROSS JOIN (
    VALUES
        ('/api/private/initiatives/{initiativeId}/comments', 'GET'),
        ('/api/private/initiatives/{initiativeId}/comments', 'POST'),
        ('/api/private/comments/{id}', 'PUT'),
        ('/api/private/comments/{id}', 'DELETE')
) AS perms(endpoint, method)
WHERE ut.name IN ('admin', 'manager', 'user')
ON CONFLICT DO NOTHING;