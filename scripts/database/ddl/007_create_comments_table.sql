-- Tabela de comentários
CREATE TABLE IF NOT EXISTS initiative_comments
(
    id            BIGSERIAL PRIMARY KEY,
    initiative_id BIGINT    NOT NULL REFERENCES initiatives (id) ON DELETE CASCADE,
    user_id       BIGINT    NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    content       TEXT      NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Índices
CREATE INDEX IF NOT EXISTS idx_comments_initiative ON initiative_comments (initiative_id);
CREATE INDEX IF NOT EXISTS idx_comments_user ON initiative_comments (user_id);
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON initiative_comments (created_at DESC);

-- Inserir comentários de exemplo
INSERT INTO initiative_comments (initiative_id, user_id, content)
VALUES (1, 1, 'Excelente iniciativa! Vamos priorizar esta automação.'),
       (1, 1, 'Precisamos definir o cronograma de implementação.'),
       (2, 1, 'Dashboard aprovado, aguardando início do desenvolvimento.'),
       (3, 1, 'Integração crítica para o fechamento fiscal.'),
       (4, 1, 'Melhoria aprovada para o próximo sprint.')
ON CONFLICT DO NOTHING;