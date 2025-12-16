-- Permissões para usar IA (todos os usuários autenticados)
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, '/api/private/ai/refine-text', 'POST'
FROM user_type ut
WHERE ut.name IN ('admin', 'manager', 'user')
ON CONFLICT DO NOTHING;