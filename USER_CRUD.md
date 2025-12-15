# 👥 CRUD de Usuários

## Endpoints Disponíveis

### 1. Criar Usuário (Admin only)
```bash
POST /private/users
Content-Type: application/json
Cookie: auth_token=... 

{
  "email": "novo@exemplo.com",
  "name": "Novo Usuário",
  "password": "senha123",
  "type_ids": [1, 2]  // IDs dos tipos a atribuir
}
```

**Response:**
```json
{
  "success": true,
  "message": "Usuário criado com sucesso",
  "user": {
    "id": 2,
    "email": "novo@exemplo.com",
    "name": "Novo Usuário",
    "created_at": "2025-12-15T10:30:00Z"
  }
}
```

### 2. Listar Todos os Usuários (Admin only)
```bash
GET /private/users
Cookie: auth_token=... 
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "email": "admin@hackathon.com",
      "name": "Admin User",
      "user_types": [
        {
          "id": 1,
          "name": "admin",
          "description": "Administrador com acesso total"
        }
      ],
      "created_at": "2025-12-15 10:00:00"
    }
  ],
  "count": 1
}
```

### 3. Buscar Usuário por ID (Admin only)
```bash
GET /private/users/1
Cookie: auth_token=...
```

### 4. Atualizar Usuário (Admin only)
```bash
PUT /private/users/1
Content-Type: application/json
Cookie: auth_token=...

{
  "name": "Nome Atualizado",
  "email": "novo-email@exemplo.com",
  "type_ids": [2, 3]  // Opcional: atualizar tipos
}
```

### 5. Deletar Usuário (Admin only)
```bash
DELETE /private/users/2
Cookie: auth_token=... 
```

**Regra:** Não é possível deletar a própria conta.

### 6. Trocar Senha (Qualquer usuário autenticado)
```bash
POST /private/change-password
Content-Type: application/json
Cookie: auth_token=...

{
  "old_password": "senha_antiga",
  "new_password":  "nova_senha123"
}
```

## Validações

### Email
- Deve ser um email válido (regex)
- Deve ser único no sistema

### Nome
- Mínimo 3 caracteres

### Senha
- Mínimo 6 caracteres
- Armazenada com bcrypt

## Permissões

- **Admin**: Acesso total ao CRUD
- **Manager**: Apenas visualizar usuários
- **User**: Apenas trocar própria senha

## Exemplos de Uso

### Criar usuário manager
```bash
curl -X POST http://localhost:8080/private/users \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "email": "manager@hackathon.com",
    "name": "Manager User",
    "password": "manager123",
    "type_ids": [2]
  }'
```

### Atualizar apenas o nome
```bash
curl -X PUT http://localhost:8080/private/users/2 \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "name": "Novo Nome"
  }'
```

### Trocar senha
```bash
curl -X POST http://localhost:8080/private/change-password \
  -H "Content-Type:  application/json" \
  -b cookies.txt \
  -d '{
    "old_password": "senha_antiga",
    "new_password": "senha_nova123"
  }'
```