# Script de Migração: Separar Brokers de Users

## 📋 Objetivo

Este script migra usuários administrativos (sem CRECI) da collection `/tenants/{id}/brokers/` para a nova collection `/tenants/{id}/users/`.

## 🎯 Lógica de Decisão

O script analisa cada registro em `/brokers` e decide se deve migrar baseado em:

### ✅ **Migrar para /users** (Admin Users)
- **Sem CRECI válido** E role = "admin" ou "manager"
- **Sem CRECI válido** E sem role de broker
- CRECI vazio, "-", "PENDENTE", "n/a", ou muito curto (< 3 chars)

### ⏸️ **Manter em /brokers** (Real Brokers)
- **Tem CRECI válido** (formato: XXXXX-F/UF ou XXXXX-J/UF)
- Independente do role, se tem CRECI = é broker real

### 🔄 **Ações Especiais**
- **Broker com CRECI + role "admin"** → Mantém em /brokers mas muda role para "broker_admin"
- **Sem CRECI + sem role definido** → Migra para /users como "admin"

## 🚀 Como Usar

### 1. **Dry-Run (Recomendado) - Apenas Relatório**

```bash
cd backend
go run cmd/migrate-users/main.go --dry-run=true --csv=migration-report.csv
```

**O que faz:**
- ✅ Analisa todos os brokers
- ✅ Gera relatório detalhado em CSV
- ✅ Mostra estatísticas no console
- ❌ **NÃO faz alterações** no Firestore

### 2. **Execução Real - Aplica Mudanças**

⚠️ **ATENÇÃO:** Faça backup do Firestore antes!

```bash
cd backend
go run cmd/migrate-users/main.go --dry-run=false --csv=migration-execution.csv
```

**O que faz:**
- ✅ Cria registros em `/tenants/{id}/users/`
- ✅ Remove registros de `/tenants/{id}/brokers/`
- ✅ Mantém mesmo ID para traceability
- ✅ Atualiza Firebase Custom Claims (se necessário)
- ✅ Rollback automático em caso de erro

## 📊 Relatório CSV

O script gera um arquivo CSV com as seguintes colunas:

| Coluna | Descrição |
|--------|-----------|
| Tenant ID | ID do tenant |
| Tenant Name | Nome do tenant |
| Broker ID | ID do registro |
| Name | Nome da pessoa |
| Email | Email |
| CRECI | CRECI registrado (ou vazio) |
| Role | Role atual (admin, broker, manager) |
| Type | Classificação (Real Broker, Admin User, etc.) |
| Action | Ação tomada (Migrate, Keep, etc.) |
| Status | Status da operação (OK, ERROR, MIGRATED) |
| Notes | Notas adicionais ou mensagens de erro |

### Exemplo de Relatório:

```csv
Tenant ID,Tenant Name,Broker ID,Name,Email,CRECI,Role,Type,Action,Status,Notes
tenant-1,ALTATECH Imóveis,user-1,Administrao,admin@example.com,-,admin,Admin User (No CRECI),Migrate to /users,MIGRATED,
tenant-1,ALTATECH Imóveis,broker-1,João Silva,joao@example.com,12345-F/SP,broker,Real Broker,Keep in /brokers,OK,
tenant-1,ALTATECH Imóveis,broker-2,Maria Admin,maria@example.com,67890-J/RJ,admin,Broker Admin (Has CRECI),Keep in /brokers,OK,Update role to broker_admin
```

## 📈 Estatísticas

Ao final, o script mostra:

```
===========================================
MIGRATION STATISTICS
===========================================
Total Brokers Found:    7
Real Brokers (CRECI):   6
Admin Users (No CRECI): 1
Migrated to /users:     1
Kept in /brokers:       6
Errors:                 0
===========================================
```

## 🔍 Exemplos de Casos

### Caso 1: Admin Puro (SEM CRECI)
```
Nome: Administrao
Email: admin@altatech.com
CRECI: -
Role: admin
→ AÇÃO: Migrar para /users com role="admin"
```

### Caso 2: Corretor Real
```
Nome: João Silva
Email: joao@example.com
CRECI: 12345-F/SP
Role: broker
→ AÇÃO: Manter em /brokers
```

### Caso 3: Corretor que é Admin
```
Nome: Maria Santos
Email: maria@example.com
CRECI: 67890-J/RJ
Role: admin
→ AÇÃO: Manter em /brokers, atualizar role para "broker_admin"
```

### Caso 4: Registro sem CRECI e sem Role
```
Nome: Pedro Gerente
Email: pedro@example.com
CRECI: (vazio)
Role: (vazio)
→ AÇÃO: Migrar para /users com role="admin"
```

## ⚠️ Avisos Importantes

1. **Backup First!**
   - Faça backup completo do Firestore antes de executar
   - Use o Firebase Console → Firestore → Export

2. **Dry-Run Always First**
   - SEMPRE execute com `--dry-run=true` primeiro
   - Revise o CSV antes de executar de verdade

3. **Firebase Custom Claims**
   - Após migração, pode ser necessário atualizar custom claims
   - Usuários migrados para /users podem precisar re-login

4. **Testing**
   - Teste primeiro em ambiente de staging
   - Valide que a aplicação funciona após migração

5. **Rollback**
   - Mantenha o CSV de execução para possível rollback manual
   - Rollback automático só funciona durante a transação

## 🧪 Teste em Staging

```bash
# 1. Copie serviceAccountKey de staging
cp serviceAccountKey.staging.json backend/serviceAccountKey.json

# 2. Execute dry-run
cd backend
go run cmd/migrate-users/main.go --dry-run=true

# 3. Revise o CSV
cat migration-report.csv

# 4. Se OK, execute de verdade
go run cmd/migrate-users/main.go --dry-run=false

# 5. Valide manualmente no Firebase Console
```

## 🔧 Troubleshooting

### Erro: "serviceAccountKey.json not found"
- Baixe a service account key do Firebase Console
- Coloque na raiz de `backend/`

### Erro: "Permission denied"
- Verifique que a service account tem permissão de admin no Firestore

### Erro: "failed to create user"
- Verifique se já existe um user com mesmo ID
- Revise logs para detalhes específicos

## 📝 Pós-Migração

Após executar com sucesso:

1. ✅ Verificar no Firebase Console:
   - `/tenants/{id}/brokers/` deve ter apenas brokers com CRECI
   - `/tenants/{id}/users/` deve ter os admins migrados

2. ✅ Testar login:
   - Admins migrados podem precisar fazer logout/login

3. ✅ Verificar frontend:
   - Página "Corretores" não deve mostrar admins
   - Nova página "Equipe" deve mostrar admins

4. ✅ Logs de atividade:
   - Verificar que não há erros relacionados

## 📞 Suporte

Em caso de problemas:
1. Pare a migração (Ctrl+C)
2. Salve o CSV gerado
3. Verifique os logs de erro
4. Reverta mudanças se necessário (use backup)
