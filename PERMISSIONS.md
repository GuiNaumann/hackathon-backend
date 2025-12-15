# 🔐 Sistema de Permissões

## Como Funciona

1. **user_type**: Define os tipos de usuário (admin, manager, user)
2. **type_user**: Relaciona usuários com seus tipos
3. **user_type_permissions**: Define quais endpoints cada tipo pode acessar

## Middlewares

### 1. AuthMiddleware
Valida o token e injeta o usuário no contexto.

### 2. PermissionMiddleware
Verifica se o usuário tem permissão para acessar o endpoint específico.

## Endpoints

### GET /private/personal-information
Retorna informações do usuário e seus tipos.

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "admin@hackathon.com",
      "name": "Admin User"
    },
    "user_types": [
      {
        "id": 1,
        "name": "admin",
        "description": "Administrador com acesso total"
      }
    ]
  }
}
```

### POST /private/admin/users/{userId}/types/{typeId}
Atribui um tipo a um usuário (apenas admin).

### DELETE /private/admin/users/{userId}/types/{typeId}
Remove um tipo de um usuário (apenas admin).

## Exemplo de Uso

### 1. Fazer Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@hackathon.com","password":"admin123"}' \
  -c cookies.txt
```

### 2. Ver Informações Pessoais
```bash
curl -X GET http://localhost:8080/private/personal-information \
  -b cookies.txt
```

### 3. Tentar Acessar Endpoint Sem Permissão
Se você não for admin, receberá:
```json
{
  "success": false,
  "error": "Você não tem permissão para acessar este recurso",
  "code": 403
}
```

## Como Adicionar Novas Permissões

```sql
-- 1. Criar novo tipo (opcional)
INSERT INTO user_type (name, description) 
VALUES ('moderator', 'Moderador com permissões específicas');

-- 2. Adicionar permissão ao tipo
INSERT INTO user_type_permissions (user_type_id, endpoint, method)
SELECT id, '/private/novo-endpoint', 'GET'
FROM user_type WHERE name = 'moderator';

-- 3. Atribuir tipo ao usuário
INSERT INTO type_user (user_id, user_type_id)
VALUES (1, (SELECT id FROM user_type WHERE name = 'moderator'));
```

## Middleware Customizado por Tipo

Se você quiser proteger uma rota específica por tipo:

```go
// No setup. go ou no módulo
adminRouter := privateRouter.PathPrefix("/admin").Subrouter()
adminRouter.Use(RequireUserType(permUseCase, "admin"))
```