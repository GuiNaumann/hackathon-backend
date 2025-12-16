-- Tabela de histórico de mudanças de status
CREATE TABLE IF NOT EXISTS initiative_history (
                                                  id BIGSERIAL PRIMARY KEY,
                                                  initiative_id BIGINT NOT NULL REFERENCES initiatives(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    old_status VARCHAR(50) NOT NULL,
    new_status VARCHAR(50) NOT NULL,
    reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

-- Índices
CREATE INDEX IF NOT EXISTS idx_history_initiative ON initiative_history(initiative_id);
CREATE INDEX IF NOT EXISTS idx_history_created_at ON initiative_history(created_at DESC);

-- Inserir histórico inicial para as iniciativas existentes
INSERT INTO initiative_history (initiative_id, user_id, old_status, new_status, reason, created_at)
SELECT
    i.id,
    i.owner_id,
    'Rascunho',
    i.status,
    'Iniciativa criada',
    i.created_at
FROM initiatives i
WHERE NOT EXISTS (
    SELECT 1 FROM initiative_history ih WHERE ih.initiative_id = i.id
);