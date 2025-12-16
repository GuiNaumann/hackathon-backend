-- Permissões para revisar iniciativas (Admin e Manager)
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, endpoint, method
FROM user_type ut
         CROSS JOIN (
    VALUES
        ('/api/private/initiatives/submitted', 'GET'),
        ('/api/private/initiatives/{id}/review', 'POST')
) AS perms(endpoint, method)
WHERE ut.name IN ('admin', 'manager')
ON CONFLICT DO NOTHING;