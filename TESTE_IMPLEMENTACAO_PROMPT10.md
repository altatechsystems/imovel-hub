# Como Testar a Implementação do PROMPT 10

## Data: 06 de Janeiro de 2026

Este guia contém todos os comandos e URLs para testar a implementação completa do sistema de separação entre corretores e usuários administrativos.

## 1. Verificar Backend

### Compilar o Backend
```bash
cd backend
go build -o bin/server.exe cmd/server/main.go
```

**Resultado esperado:** Compilação sem erros

### Rodar os Testes Unitários
```bash
cd backend
go test ./internal/models -v
```

**Resultado esperado:**
```
=== RUN   TestUser_HasPermission
--- PASS: TestUser_HasPermission (0.00s)
=== RUN   TestUser_AddPermission
--- PASS: TestUser_AddPermission (0.00s)
=== RUN   TestUser_RemovePermission
--- PASS: TestUser_RemovePermission (0.00s)
=== RUN   TestIsValidUserRole
--- PASS: TestIsValidUserRole (0.00s)
=== RUN   TestValidUserRoles
--- PASS: TestValidUserRoles (0.00s)
PASS
ok      github.com/altatech/ecosistema-imob/backend/internal/models
```

### Iniciar o Servidor Backend
```bash
cd backend
go run cmd/server/main.go
```

**Resultado esperado:**
```
Connecting to Firestore database: imob-dev
✅ Connected to Firestore database: imob-dev
Firebase initialized successfully
Repositories initialized
Services initialized
Handlers initialized
[GIN-debug] Listening on :8080
```

### Verificar Rotas Registradas
Ao iniciar o servidor, procurar pelas rotas de usuários:
```
[GIN-debug] POST   /api/v1/admin/:tenant_id/users
[GIN-debug] GET    /api/v1/admin/:tenant_id/users/:userId
[GIN-debug] PUT    /api/v1/admin/:tenant_id/users/:userId
[GIN-debug] DELETE /api/v1/admin/:tenant_id/users/:userId
[GIN-debug] GET    /api/v1/admin/:tenant_id/users
[GIN-debug] POST   /api/v1/admin/:tenant_id/users/:userId/permissions
[GIN-debug] DELETE /api/v1/admin/:tenant_id/users/:userId/permissions/:permission
```

## 2. Testar API com cURL/Postman

### Listar Usuários Administrativos
```bash
curl -X GET "http://localhost:8080/api/v1/admin/{TENANT_ID}/users" \
  -H "Authorization: Bearer {TOKEN}"
```

**Resultado esperado:** JSON com lista de 4 usuários (os que foram migrados)

### Listar Apenas Usuários Ativos
```bash
curl -X GET "http://localhost:8080/api/v1/admin/{TENANT_ID}/users?active=true" \
  -H "Authorization: Bearer {TOKEN}"
```

### Buscar Usuário Específico
```bash
curl -X GET "http://localhost:8080/api/v1/admin/{TENANT_ID}/users/{USER_ID}" \
  -H "Authorization: Bearer {TOKEN}"
```

### Criar Novo Usuário Administrativo
```bash
curl -X POST "http://localhost:8080/api/v1/admin/{TENANT_ID}/users" \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "firebase_uid": "test-uid-123",
    "name": "Teste Usuário",
    "email": "teste@example.com",
    "role": "manager",
    "is_active": true,
    "permissions": ["property.view", "property.create"]
  }'
```

**Nota:** Este endpoint NÃO deve aceitar usuários sem especificar role ou com CRECI.

### Listar Corretores (para verificar que não há usuários admin)
```bash
curl -X GET "http://localhost:8080/api/v1/admin/{TENANT_ID}/brokers" \
  -H "Authorization: Bearer {TOKEN}"
```

**Resultado esperado:** JSON com lista de apenas 4 corretores (os que têm CRECI válido)

## 3. Verificar Collections no Firestore

### Acessar Console do Firestore
```
https://console.firebase.google.com/project/ecosistema-imob-dev/firestore
```

### Verificar Collection `/tenants/{TENANT_ID}/users`
**Deve conter 4 documentos:**
1. Daniel Garcia (admin - ALTATECH Systems)
2. Administrao (admin - ALTATECH Imóveis)
3. Pablo Silva (sem CRECI válido)
4. Franco Barroso (sem CRECI válido)

**Cada documento deve ter:**
- `role`: "admin" ou "manager"
- `is_active`: true/false
- `permissions`: array (pode estar vazio para admins)
- **NÃO deve ter** campo `creci`

### Verificar Collection `/tenants/{TENANT_ID}/brokers`
**Deve conter 4 documentos:**
1. Alex Reis (CRECI: 48.346)
2. Daniel Garcia (CRECI: 34.134)
3. Suzana Costa (CRECI: 52.648)
4. Fernanda Reis (CRECI: 34.166)

**Cada documento deve ter:**
- `creci`: string com número válido (não vazio, não "PENDENTE")
- `role`: "broker" ou "broker_admin"
- `is_active`: true/false

## 4. Testar Frontend Admin

### Iniciar Frontend
```bash
cd frontend-admin
npm run dev
```

**URL:** http://localhost:3000

### Testar Página de Equipe
**URL:** http://localhost:3000/dashboard/equipe

**Verificações:**
- [ ] Página carrega sem erros
- [ ] Mostra card informativo azul explicando diferença entre usuários e corretores
- [ ] Link para página de Corretores funciona
- [ ] Lista os 4 usuários administrativos
- [ ] Filtro "Mostrar apenas ativos" funciona
- [ ] Contador mostra "Total: 4 usuário(s)"
- [ ] Cada usuário mostra:
  - Avatar ou inicial
  - Nome e documento
  - Email e telefone
  - Badge de perfil (Admin = roxo, Manager = azul)
  - Permissões ("Acesso total" para admin, "X permissões" para manager)
  - Status (Ativo/Inativo)
  - Botões Editar e Excluir
- [ ] Botão "Novo Usuário" redireciona para `/dashboard/equipe/novo`
- [ ] Loading state aparece durante carregamento
- [ ] Estado vazio aparece quando não há usuários

### Testar Página de Corretores
**URL:** http://localhost:3000/dashboard/corretores

**Verificações:**
- [ ] Página carrega sem erros
- [ ] Subtítulo mostra "(CRECI obrigatório)"
- [ ] Card informativo azul explicando que a página é só para corretores
- [ ] Link para página de Equipe funciona
- [ ] Lista apenas os 4 corretores (com CRECI válido)
- [ ] **NÃO mostra** "Administrao" nem outros usuários sem CRECI
- [ ] Estatísticas mostram apenas corretores reais

### Testar Exclusão de Usuário
1. Ir para `/dashboard/equipe`
2. Clicar em "Excluir" em algum usuário
3. Confirmar exclusão
4. Verificar que:
   - [ ] Modal de confirmação aparece
   - [ ] Após confirmar, usuário é removido da lista
   - [ ] Contador é atualizado
   - [ ] Request DELETE é feito para API

## 5. Verificar Validações

### Backend - Tentar Criar Corretor sem CRECI
```bash
curl -X POST "http://localhost:8080/api/v1/admin/{TENANT_ID}/brokers" \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "firebase_uid": "test-uid",
    "name": "Teste Corretor",
    "email": "teste@example.com",
    "role": "broker"
  }'
```

**Resultado esperado:** Erro indicando que CRECI é obrigatório

### Backend - Tentar Criar Corretor com Role Admin
```bash
curl -X POST "http://localhost:8080/api/v1/admin/{TENANT_ID}/brokers" \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "firebase_uid": "test-uid",
    "name": "Teste Admin",
    "email": "teste@example.com",
    "creci": "12345",
    "role": "admin"
  }'
```

**Resultado esperado:** Erro dizendo "invalid role for broker: must be 'broker' or 'broker_admin'. Administrative users should be created in /users collection"

### Backend - Tentar Criar Usuário com Role Broker
```bash
curl -X POST "http://localhost:8080/api/v1/admin/{TENANT_ID}/users" \
  -H "Authorization: Bearer {TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "firebase_uid": "test-uid",
    "name": "Teste User",
    "email": "teste@example.com",
    "role": "broker"
  }'
```

**Resultado esperado:** Erro dizendo "invalid role: must be 'admin' or 'manager'"

## 6. Verificar Logs de Atividades

### Acessar ActivityLog Collection
```
Firestore > tenants > {TENANT_ID} > activity_logs
```

**Deve conter logs de:**
- Criação de usuários
- Atualização de usuários
- Exclusão de usuários
- Concessão de permissões
- Revogação de permissões

Cada log deve ter:
- `event_type`: "user.created", "user.updated", etc.
- `actor_type`: "user" ou "broker"
- `actor_id`: ID do usuário que executou a ação
- `metadata`: Detalhes da operação
- `created_at`: Timestamp

## 7. Arquivos de Migração

### Verificar Relatórios de Migração
```bash
cd backend
cat migration-report.csv
cat migration-execution-report.csv
```

**Deve mostrar:**
- 8 linhas de dados (além do header)
- 4 registros com Action = "Migrate to /users" e Status = "MIGRATED"
- 4 registros com Action = "Keep in /brokers" e Status = "OK"

## 8. Checklist Final

### Backend
- [ ] ✅ Servidor compila sem erros
- [ ] ✅ Testes unitários passam (6/6)
- [ ] ✅ Rotas de usuários registradas
- [ ] ✅ API retorna apenas usuários admin em /users
- [ ] ✅ API retorna apenas corretores com CRECI em /brokers
- [ ] ✅ Validações impedem criação incorreta
- [ ] ✅ Activity logs são criados

### Frontend
- [ ] ✅ Página Equipe carrega e funciona
- [ ] ✅ Página Corretores atualizada com avisos
- [ ] ✅ TypeScript types funcionam
- [ ] ✅ API client faz requests corretas
- [ ] ✅ Filtros funcionam
- [ ] ✅ Ações (editar, excluir) funcionam

### Firestore
- [ ] ✅ Collection /users contém 4 documentos
- [ ] ✅ Collection /brokers contém 4 documentos
- [ ] ✅ Nenhum usuário admin em /brokers
- [ ] ✅ Nenhum corretor em /users

### UX
- [ ] ✅ Cards informativos claros
- [ ] ✅ Links entre páginas funcionam
- [ ] ✅ Diferença entre corretores e usuários é clara
- [ ] ✅ Usuários entendem onde gerenciar cada tipo

---

## Problemas Conhecidos

### Frontend - Páginas de Criação/Edição
As páginas `/dashboard/equipe/novo` e `/dashboard/equipe/[id]` ainda não foram implementadas. Ao clicar em "Novo Usuário" ou "Editar", o usuário será redirecionado mas verá erro 404.

**Solução:** Implementar essas páginas ou desabilitar os botões temporariamente.

### Backend - Signup Flow
O fluxo de signup ainda não diferencia automaticamente entre broker e admin user. Usuários devem ser criados manualmente via API.

**Solução:** Implementar lógica no signup para detectar se tem CRECI e criar no endpoint correto.

---

**Data do Teste:** _______________________

**Testado por:** _______________________

**Resultado:** [ ] ✅ Todos os testes passaram  [ ] ❌ Alguns testes falharam

**Observações:**
_____________________________________________
_____________________________________________
_____________________________________________
