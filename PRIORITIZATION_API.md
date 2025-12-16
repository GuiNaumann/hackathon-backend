# 📋 API de Priorização de Iniciativas

## 🎯 Visão Geral

Sistema de priorização anual de iniciativas por setor com controle de permissões.

---

## 🔐 Permissões

| Tipo de Usuário | Permissões |
|-----------------|------------|
| **User** | Vê e prioriza iniciativas do **seu setor**.  Ao salvar, bloqueia.  Para alterar novamente, precisa solicitar aprovação. |
| **Manager/Admin** | Vê e prioriza iniciativas de **todos os setores**.  Pode alterar a qualquer momento sem aprovação.  Pode aprovar/recusar solicitações. |

---

## 📡 Endpoints

### **1. Buscar Priorização do Meu Setor (User)**

```http
GET /api/private/prioritization? year=2025
Headers:

Code
Authorization: Bearer {token}
Response:

JSON
{
  "success": true,
  "data": {
    "id": 1,
    "sector_id": 4,
    "sector_name": "TI",
    "year": 2025,
    "is_locked": true,
    "initiatives": [
      {
        "id": 10,
        "title": "Automatizar processo de aprovação",
        "description": ".. .",
        "status": "Aprovada",
        "type": "Automação",
        "priority": "Alta",
        "sector":  "TI",
        "owner_name": "João Silva",
        "date":  "15 de jan, 2025"
      }
    ],
    "created_by_user_id": 5,
    "created_by_name": "Maria Santos",
    "created_at": "2025-01-10 10:00:00",
    "updated_at": "2025-01-10 10:00:00"
  }
}
2. Salvar Priorização (User)
HTTP
POST /api/private/prioritization
Content-Type: application/json

{
  "year": 2025,
  "priority_order": [10, 5, 12, 8, 3]
}
Regras:

✅ Se não existir priorização, cria e bloqueia
✅ Se existir e estiver desbloqueada, atualiza e bloqueia
❌ Se existir e estiver bloqueada, retorna erro (precisa solicitar aprovação)
✅ Admin/Manager pode salvar a qualquer momento
Response:

JSON
{
  "success": true,
  "message": "Priorização salva com sucesso",
  "data": { /* ... priorização completa ... */ }
}
3. Solicitar Mudança na Priorização (User)
Usado quando a priorização está bloqueada e o usuário quer alterar.

HTTP
POST /api/private/prioritization/request-change? year=2025
Content-Type:  application/json

{
  "new_priority_order": [12, 10, 8, 5, 3],
  "reason": "Mudança nas prioridades do trimestre devido a novas demandas do cliente"
}
Response:

JSON
{
  "success": true,
  "message": "Solicitação de mudança criada com sucesso",
  "data": {
    "id": 1,
    "prioritization_id": 1,
    "requested_by_user_id": 5,
    "requested_by_name": "João Silva",
    "new_priority_order": [12, 10, 8, 5, 3],
    "reason": "Mudança nas prioridades.. .",
    "status": "Pendente",
    "created_at": "2025-02-10T14:30:00Z"
  }
}
4. Buscar Priorização de Todos os Setores (Admin/Manager)
HTTP
GET /api/private/prioritization/all? year=2025
Response:

JSON
{
  "success": true,
  "data": {
    "year": 2025,
    "sectors": [
      {
        "id": 1,
        "sector_id": 1,
        "sector_name": "Comercial",
        "year": 2025,
        "is_locked": true,
        "initiatives": [ /* ... */ ],
        "created_by_user_id": 3,
        "created_by_name": "Admin User",
        "created_at":  "2025-01-05 09:00:00",
        "updated_at": "2025-01-05 09:00:00"
      },
      {
        "id": 2,
        "sector_id": 4,
        "sector_name":  "TI",
        "year": 2025,
        "is_locked": false,
        "initiatives": [ /* ... */ ]
      }
    ]
  }
}
5. Listar Solicitações Pendentes (Admin/Manager)
HTTP
GET /api/private/prioritization/change-requests
Response:

JSON
{
  "success": true,
  "count": 2,
  "data": [
    {
      "id": 1,
      "prioritization_id": 1,
      "requested_by_user_id": 5,
      "requested_by_name": "João Silva",
      "new_priority_order": [12, 10, 8, 5, 3],
      "reason":  "Mudança nas prioridades do trimestre",
      "status": "Pendente",
      "created_at": "2025-02-10T14:30:00Z"
    }
  ]
}
6. Aprovar/Recusar Solicitação de Mudança (Admin/Manager)
HTTP
POST /api/private/prioritization/change-requests/1/review
Content-Type: application/json

{
  "approved": true,
  "reason":  "Aprovado conforme justificativa apresentada"
}
Response:

JSON
{
  "success": true,
  "message": "Solicitação de mudança aprovada e a priorização foi desbloqueada para alteração"
}
Comportamento:

✅ Se aprovado: Desbloqueia a priorização para o usuário poder editar
❌ Se recusado: Mantém bloqueada
🎨 Exemplo de Fluxo (Frontend)
User:
Acessa /prioritization? year=2025
Se is_locked:  false → Pode arrastar e reorganizar
Clica em "Salvar Priorização" → Envia POST /prioritization
Sistema bloqueia (is_locked: true)
Se quiser alterar novamente → Clica em "Solicitar Mudança" → POST /prioritization/request-change
Aguarda aprovação do gestor
Admin/Manager:
Acessa /prioritization/all?year=2025
Vê todos os setores com suas priorizações
Pode arrastar e reorganizar qualquer setor a qualquer momento
Clica em "Salvar" → Atualiza sem bloqueio (pode alterar novamente depois)
Vê solicitações pendentes em /change-requests
Aprova/Recusa → POST /change-requests/{id}/review
🐛 Erros Comuns
Status	Erro	Solução
400	priorização já está bloqueada.  Solicite aprovação para alterá-la	Usuário normal tentou salvar priorização bloqueada. Deve usar /request-change
400	usuário não está vinculado a um setor	Usuário não tem sector_id. Precisa ser vinculado a um setor
400	já existe uma solicitação de mudança pendente	Aguardar aprovação da solicitação anterior
403	apenas administradores e gerentes podem... 	Endpoint restrito a admin/manager
✅ Checklist Frontend
 Tela de priorização para usuários (drag-and-drop)
 Botão "Salvar Priorização" (desabilitar se is_locked:  true)
 Modal "Solicitar Mudança" (quando bloqueado)
 Tela de priorização global para admin/manager (todos os setores)
 Listagem de solicitações pendentes
 Modal de aprovação/recusa de solicitações
 Indicador visual de priorização bloqueada/desbloqueada
 Filtro por ano


---

## ✅ **Resumo Final:**

✅ **Criado sistema completo de priorização**  
✅ **User**: Vê e prioriza seu setor, ao salvar bloqueia, precisa solicitar aprovação para alterar  
✅ **Admin/Manager**: Vê todos os setores, pode alterar a qualquer momento, aprova solicitações  
✅ **Tabelas**: `initiative_prioritization` e `prioritization_change_requests`  
✅ **Endpoints completos** com documentação  
✅ **Integrado no setup. go**  

Está tudo pronto!  🎉  
Quer que eu te ajude com algum ajuste ou com a implementação no frontend? 