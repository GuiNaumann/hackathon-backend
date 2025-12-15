-- Tabela de iniciativas
CREATE TABLE IF NOT EXISTS initiatives
(
    id          BIGSERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL,
    benefits    TEXT         NOT NULL,
    status      VARCHAR(50)  NOT NULL DEFAULT 'Submetida',
    type        VARCHAR(50)  NOT NULL,
    priority    VARCHAR(20)  NOT NULL,
    sector      VARCHAR(100) NOT NULL,
    owner_id    BIGINT       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    deadline    DATE,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- Índices
CREATE INDEX IF NOT EXISTS idx_initiatives_owner ON initiatives (owner_id);
CREATE INDEX IF NOT EXISTS idx_initiatives_status ON initiatives (status);
CREATE INDEX IF NOT EXISTS idx_initiatives_type ON initiatives (type);
CREATE INDEX IF NOT EXISTS idx_initiatives_sector ON initiatives (sector);
CREATE INDEX IF NOT EXISTS idx_initiatives_created_at ON initiatives (created_at DESC);

-- Inserir iniciativas de exemplo
INSERT INTO initiatives (title, description, benefits, status, type, priority, sector, owner_id, deadline)
VALUES ('Automação do processo de admissão',
        'Implementar workflow automatizado para o processo de admissão de novos colaboradores, incluindo integração com sistemas de RH e documentação digital.',
        'Redução do tempo de admissão em 50%, eliminação de erros manuais e melhor experiência para novos colaboradores.',
        'Em Execução', 'Automação', 'Alta', 'Recursos Humanos', 1, '2024-03-15'),
       ('Dashboard de vendas em tempo real',
        'Criar painel executivo com indicadores de vendas atualizados em tempo real, incluindo metas, conversão e performance por vendedor.',
        'Visibilidade imediata dos resultados de vendas, tomada de decisão mais rápida e acompanhamento de metas em tempo real.',
        'Em Análise', 'Novo Projeto', 'Média', 'Comercial', 1, '2024-04-01'),
       ('Integração SPED Fiscal com ERP',
        'Desenvolver integração automática entre o sistema SPED Fiscal e o ERP corporativo para envio automatizado de obrigações fiscais.',
        'Eliminação de trabalho manual, redução de erros fiscais e conformidade automática com legislação. ',
        'Aprovada', 'Integração', 'Alta', 'Fiscal', 1, '2024-02-28'),
       ('Melhoria no módulo de relatórios',
        'Adicionar novos filtros e opções de exportação ao módulo de relatórios do sistema de gestão de produtos.',
        'Mais flexibilidade na geração de relatórios e melhor análise de dados de produtos.', 'Submetida', 'Melhoria',
        'Baixa', 'Produto', 1, NULL),
       ('Portal do colaborador mobile',
        'Desenvolver aplicativo mobile para acesso ao portal do colaborador com funcionalidades de ponto, férias e holerite.',
        'Maior comodidade para colaboradores, redução de demandas ao RH e modernização do acesso às informações.',
        'Devolvida', 'Novo Projeto', 'Média', 'Recursos Humanos', 1, '2024-06-01'),
       ('Integração com marketplace',
        'Conectar sistema de vendas com principais marketplaces (Mercado Livre, Amazon, B2W) para gestão centralizada. ',
        'Ampliação dos canais de venda, gestão unificada de estoque e pedidos, aumento do faturamento.', 'Reprovada',
        'Integração', 'Alta', 'Comercial', 1, NULL)
ON CONFLICT DO NOTHING;