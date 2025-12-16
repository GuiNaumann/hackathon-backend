-- Tabela de priorização por setor (salva a ordem das iniciativas)
CREATE TABLE IF NOT EXISTS initiative_prioritization (
                                                         id BIGSERIAL PRIMARY KEY,
                                                         sector_id BIGINT NOT NULL REFERENCES sectors(id) ON DELETE CASCADE,
    year INT NOT NULL, -- Ano da priorização (ex: 2025)
    priority_order JSONB NOT NULL, -- Array com IDs das iniciativas na ordem de prioridade
    is_locked BOOLEAN NOT NULL DEFAULT false, -- Se está bloqueada para edição
    created_by_user_id BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(sector_id, year) -- Apenas uma priorização por setor por ano
    );

-- Tabela de solicitações de mudança de priorização
CREATE TABLE IF NOT EXISTS prioritization_change_requests (
                                                              id BIGSERIAL PRIMARY KEY,
                                                              prioritization_id BIGINT NOT NULL REFERENCES initiative_prioritization(id) ON DELETE CASCADE,
    requested_by_user_id BIGINT NOT NULL REFERENCES users(id),
    new_priority_order JSONB NOT NULL, -- Nova ordem proposta
    reason TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'Pendente', -- Pendente, Aprovada, Reprovada
    reviewed_by_user_id BIGINT REFERENCES users(id),
    review_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    reviewed_at TIMESTAMP
    );

-- Índices
CREATE INDEX IF NOT EXISTS idx_prioritization_sector_year ON initiative_prioritization(sector_id, year);
CREATE INDEX IF NOT EXISTS idx_prioritization_locked ON initiative_prioritization(is_locked);
CREATE INDEX IF NOT EXISTS idx_change_requests_status ON prioritization_change_requests(status);
CREATE INDEX IF NOT EXISTS idx_change_requests_prioritization ON prioritization_change_requests(prioritization_id);

COMMENT ON TABLE initiative_prioritization IS 'Priorização anual de iniciativas por setor';
COMMENT ON TABLE prioritization_change_requests IS 'Solicitações de mudança na priorização';

-- Permissões para priorização
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT ut.id, endpoint, method
FROM user_type ut
         CROSS JOIN (
    VALUES
        -- Todos podem ver e salvar priorização do seu setor
        ('/api/private/prioritization', 'GET'),
        ('/api/private/prioritization', 'POST'),
        ('/api/private/prioritization/request-change', 'POST'),

        -- Admin e Manager podem ver tudo e aprovar mudanças
        ('/api/private/prioritization/all', 'GET'),
        ('/api/private/prioritization/change-requests', 'GET'),
        ('/api/private/prioritization/change-requests/{id}/review', 'POST')
) AS perms(endpoint, method)
WHERE ut.name IN ('admin', 'manager', 'user')
    ON CONFLICT DO NOTHING;