📋 Documentação de Mudanças - Sistema de Priorização de Iniciativas
🆕 NOVA FUNCIONALIDADE: Priorização Anual de Iniciativas
🎯 O que é?
Agora o sistema tem uma tela de priorização onde os usuários podem organizar a ordem de importância das iniciativas do seu setor para o ano.

É como uma lista de top 10 iniciativas que o setor vai focar no ano.

👥 **Como funciona para cada tipo de usuário? **
📌 Usuário Normal (User):
Vê apenas as iniciativas do SEU setor (aprovadas)
Pode arrastar e soltar para organizar a ordem de prioridade
Ao clicar em "Salvar Priorização":
✅ Sistema salva a ordem
🔒 BLOQUEIA a priorização (não pode mais alterar)
Se quiser alterar depois de salvar:
❌ Não pode mais mover as iniciativas
✅ Precisa clicar em "Solicitar Mudança"
✉️ Abre um modal para justificar por que quer mudar
⏳ Aguarda Gestor/Admin aprovar
Se aprovarem:
✅ Pode editar novamente
Se recusarem:
❌ Continua bloqueado
👔 Gestor/Admin (Manager/Admin):
Vê TODAS as iniciativas de TODOS os setores agrupadas por setor
Exemplo:
Code
📦 Setor Comercial (5 iniciativas)
📦 Setor TI (8 iniciativas)
📦 Setor RH (3 iniciativas)
Pode arrastar e soltar para organizar a priorização entre setores
Pode alterar a qualquer momento (não bloqueia para admin/manager)
Vê uma lista de solicitações de mudança pendentes dos usuários
Pode aprovar ou recusar cada solicitação
🎨 Telas que vocês precisam criar:
1. Tela de Priorização (User)
   Rota sugerida: /prioritization

Layout:

Code
┌─────────────────────────────────────────┐
│  📊 Priorização 2025 - Setor TI        │
│                                         │
│  Status: 🔒 Bloqueada                  │
│  [Solicitar Mudança]                    │
│                                         │
│  Arraste para organizar:                 │
│                                         │
│  1. 🔹 Automatizar processo de RH       │
│  2. 🔹 Integração com ERP               │
│  3. 🔹 Dashboard de vendas              │
│  4. 🔹 App mobile para clientes         │
│  5. 🔹 Migração para cloud              │
│                                         │
└─────────────────────────────────────────┘
Funcionalidades:

✅ Drag-and-drop para reordenar
✅ Botão "Salvar Priorização" (desabilitar se is_locked:  true)
✅ Botão "Solicitar Mudança" (só aparecer se is_locked: true)
✅ Filtro de ano (2024, 2025, 2026...)
✅ Indicador visual se está bloqueado ou não
2. Modal "Solicitar Mudança" (User)
   Abre quando o usuário clica em "Solicitar Mudança":

Code
┌──────────────────────────────────────────┐
│  Solicitar Mudança na Priorização        │
│                                          │
│  Nova ordem proposta:                    │
│  1. App mobile para clientes             │
│  2. Automatizar processo de RH           │
│  3. Integração com ERP                   │
│  ...                                      │
│                                          │
│  Justificativa:  *                        │
│  ┌────────────────────────────────────┐ │
│  │ Mudança nas prioridades devido a   │ │
│  │ novas demandas do cliente X...      │ │
│  └────────────────────────────────────┘ │
│                                          │
│  [Cancelar]  [Enviar Solicitação]       │
└──────────────────────────────────────────┘
3. Tela de Priorização Global (Admin/Manager)
   Rota sugerida: /prioritization/all

Layout:

Code
┌────────────────────────────────────────────┐
│  📊 Priorização Global 2025                │
│                                            │
│  🔹 Setor Comercial (5 iniciativas)       │
│     1. CRM para vendedores                 │
│     2. Portal do cliente                   │
│     3. Automação de propostas              │
│     ...                                     │
│                                            │
│  🔹 Setor TI (8 iniciativas)              │
│     1. Migração para cloud                 │
│     2. Dashboard BI                        │
│     ...                                     │
│                                            │
│  🔹 Setor RH (3 iniciativas)              │
│     1. Portal do colaborador               │
│     ...                                     │
│                                            │
│  [Salvar Priorização Global]               │
└────────────────────────────────────────────┘
Funcionalidades:

✅ Vê todos os setores ao mesmo tempo
✅ Pode arrastar entre setores para priorizar globalmente
✅ Pode alterar a qualquer momento (sem bloqueio)
✅ Accordion/collapse para cada setor
4. Tela de Solicitações Pendentes (Admin/Manager)
   Rota sugerida: /prioritization/requests

Layout:

Code
┌─────────────────────────────────────────────┐
│  📬 Solicitações de Mudança de Priorização │
│                                             │
│  🔔 2 solicitações pendentes                │
│                                             │
│  ┌───────────────────────────────────────┐ │
│  │ 👤 João Silva (Setor TI)              │ │
│  │ 📅 15 de fev, 2025                    │ │
│  │                                       │ │
│  │ Motivo:                                │ │
│  │ "Mudança nas prioridades devido a..." │ │
│  │                                       │ │
│  │ Nova ordem proposta:                  │ │
│  │ 1. App mobile                         │ │
│  │ 2. Automatizar RH                     │ │
│  │ ...                                    │ │
│  │                                       │ │
│  │ [Recusar] [Aprovar]                   │ │
│  └───────────────────────────────────────┘ │
└─────────────────────────────────────────────┘
5. Modal de Aprovação/Recusa (Admin/Manager)
   Code
   ┌──────────────────────────────────────────┐
   │  ✅ Aprovar Solicitação                  │
   │                                          │
   │  Solicitante: João Silva (TI)            │
   │  Data: 15 de fev, 2025                   │
   │                                          │
   │  Justificativa para aprovação:  *         │
   │  ┌────────────────────────────────────┐ │
   │  │ Aprovado conforme alinhamento com  │ │
   │  │ as novas diretrizes estratégicas   │ │
   │  └────────────────────────────────────┘ │
   │                                          │
   │  [Cancelar]  [Aprovar]                   │
   └──────────────────────────────────────────┘
   📡 Novos Endpoints
1. Buscar minha priorização
   TypeScript
   GET /api/private/prioritization? year=2025
2. Salvar priorização
   TypeScript
   POST /api/private/prioritization
   Body: {
   year: 2025,
   priority_order: [10, 5, 12, 8, 3] // IDs na ordem
   }
3. Solicitar mudança
   TypeScript
   POST /api/private/prioritization/request-change? year=2025
   Body: {
   new_priority_order: [12, 10, 8, 5, 3],
   reason: "Motivo da mudança..."
   }
4. Buscar priorização global (Admin/Manager)
   TypeScript
   GET /api/private/prioritization/all?year=2025
5. Listar solicitações pendentes (Admin/Manager)
   TypeScript
   GET /api/private/prioritization/change-requests
6. Aprovar/Recusar solicitação (Admin/Manager)
   TypeScript
   POST /api/private/prioritization/change-requests/{id}/review
   Body: {
   approved: true,
   reason: "Justificativa..."
   }
   🎨 Biblioteca Recomendada para Drag-and-Drop
   bash
   npm install @dnd-kit/core @dnd-kit/sortable
   ou

bash
npm install react-beautiful-dnd
🔄 Fluxo Completo (Exemplo)
Cenário: Usuário Normal
Janeiro/2025: João (TI) acessa /prioritization
Arrasta as 8 iniciativas do setor TI na ordem que acha mais importante
Clica em "Salvar Priorização"
Sistema salva e bloqueia ✅
Fevereiro/2025: João percebe que precisa mudar a ordem
Tenta arrastar mas está bloqueado 🔒
Clica em "Solicitar Mudança"
Reorganiza a ordem e escreve a justificativa
Envia → Status "Pendente" ⏳
Gestor recebe notificação
Gestor abre /prioritization/requests
Clica em "Aprovar" ✅
João recebe notificação: "Sua solicitação foi aprovada!"
João volta em /prioritization
Agora está desbloqueado e pode editar 🎉
Cenário: Admin/Manager
Janeiro/2025: Admin acessa /prioritization/all
Vê 3 setores:
Comercial (5 iniciativas)
TI (8 iniciativas)
RH (3 iniciativas)
Arrasta e organiza a priorização entre setores
Clica em "Salvar"
Qualquer momento: Pode voltar e alterar de novo (sem bloqueio)
Acessa /prioritization/requests para ver solicitações dos usuários
Aprova ou recusa conforme necessário
⚠️ Regras Importantes
Situação	User	Admin/Manager
Primeira vez salvando	Bloqueia ✅	Não bloqueia ❌
Quer alterar depois de salvar	Precisa solicitar aprovação ✉️	Pode alterar livremente ✅
Vê iniciativas de outros setores	Não ❌	Sim ✅
Pode aprovar solicitações	Não ❌	Sim ✅
🐛 Mensagens de Erro para Tratar
TypeScript
// Quando usuário tenta salvar priorização bloqueada
{
"success": false,
"error":  "priorização já está bloqueada.  Solicite aprovação para alterá-la"
}

// Quando já existe solicitação pendente
{
"success":  false,
"error": "já existe uma solicitação de mudança pendente para esta priorização"
}

// Quando usuário não tem setor
{
"success":  false,
"error": "usuário não está vinculado a um setor"
}
✅ Checklist de Implementação
User:
Tela /prioritization com drag-and-drop
Botão "Salvar Priorização"
Botão "Solicitar Mudança" (só se bloqueado)
Modal de solicitação de mudança
Indicador visual de bloqueado/desbloqueado
Filtro de ano
Notificação quando solicitação for aprovada/recusada
Admin/Manager:
Tela /prioritization/all com todos os setores
Drag-and-drop entre setores
Botão "Salvar Priorização Global"
Tela /prioritization/requests com solicitações pendentes
Modal de aprovação/recusa
Badge de notificação com quantidade de solicitações pendentes
📚 Exemplo de Código React
TSX
// Exemplo simples de drag-and-drop
import { DndContext, closestCenter } from '@dnd-kit/core';
import { SortableContext, arrayMove } from '@dnd-kit/sortable';

function Prioritization() {
const [initiatives, setInitiatives] = useState([]);
const [isLocked, setIsLocked] = useState(false);

const handleDragEnd = (event) => {
const { active, over } = event;
if (active.id !== over.id) {
setInitiatives((items) => {
const oldIndex = items. findIndex((i) => i.id === active. id);
const newIndex = items.findIndex((i) => i.id === over.id);
return arrayMove(items, oldIndex, newIndex);
});
}
};

const handleSave = async () => {
await api.post('/prioritization', {
year: 2025,
priority_order: initiatives.map(i => i.id)
});
};

return (
<div>
<h1>Priorização 2025</h1>
<p>Status: {isLocked ? '🔒 Bloqueada' : '🔓 Desbloqueada'}</p>

      <DndContext collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
        <SortableContext items={initiatives}>
          {initiatives.map((initiative, index) => (
            <DraggableItem key={initiative. id} initiative={initiative} index={index} />
          ))}
        </SortableContext>
      </DndContext>

      {! isLocked && <button onClick={handleSave}>Salvar Priorização</button>}
      {isLocked && <button onClick={openChangeRequestModal}>Solicitar Mudança</button>}
    </div>
);
}
