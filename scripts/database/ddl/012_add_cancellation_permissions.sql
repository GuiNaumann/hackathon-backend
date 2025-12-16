-- Todos podem solicitar cancelamento
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, '/api/private/initiatives/{id}/request-cancellation', 'POST'
FROM user_type ut
WHERE ut.name IN ('admin', 'manager', 'user')
    ON CONFLICT DO NOTHING;

-- Apenas admin e manager podem gerenciar cancelamentos
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, endpoint, method
FROM user_type ut
         CROSS JOIN (
    VALUES
        ('/api/private/cancellation-requests', 'GET'),
        ('/api/private/cancellation-requests/{id}/review', 'POST')
) AS perms(endpoint, method)
WHERE ut.name IN ('admin', 'manager')
    ON CONFLICT DO NOTHING;