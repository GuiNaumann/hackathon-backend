-- Tabela de solicitações de cancelamento
CREATE TABLE IF NOT EXISTS initiative_cancellation_requests (
                                                                id BIGSERIAL PRIMARY KEY,
                                                                initiative_id BIGINT NOT NULL REFERENCES initiatives(id) ON DELETE CASCADE,
    requested_by_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'Pendente', -- Pendente, Aprovada, Reprovada
    reviewed_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    review_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    reviewed_at TIMESTAMP,
    CONSTRAINT unique_pending_cancellation UNIQUE (initiative_id, status)
    );

-- Índices
CREATE INDEX IF NOT EXISTS idx_cancellation_initiative ON initiative_cancellation_requests(initiative_id);
CREATE INDEX IF NOT EXISTS idx_cancellation_status ON initiative_cancellation_requests(status);
CREATE INDEX IF NOT EXISTS idx_cancellation_requested_by ON initiative_cancellation_requests(requested_by_user_id);

COMMENT ON TABLE initiative_cancellation_requests IS 'Solicitações de cancelamento de iniciativas';
COMMENT ON COLUMN initiative_cancellation_requests.status IS 'Status da solicitação:   Pendente, Aprovada, Reprovada';