-- Todos podem ver o histórico
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, '/api/private/initiatives/{id}/history', 'GET'
FROM user_type ut
WHERE ut.name IN ('admin', 'manager', 'user')
    ON CONFLICT DO NOTHING;