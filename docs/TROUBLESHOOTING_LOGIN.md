# Troubleshooting: Erro "Broker not found" no Login

**Data:** 22/12/2025
**Status:** ✅ Resolvido

---

## Problema

Ao tentar fazer login com credenciais válidas, o sistema retornava erro:
```
{"error": "Broker not found"}
```

## Causa Raiz

O backend utiliza **CollectionGroup query** para buscar brokers em todas as subcoleções:

```go
brokersQuery := h.firestoreDB.CollectionGroup("brokers").
    Where("firebase_uid", "==", userRecord.UID).
    Limit(1)
```

**CollectionGroup queries requerem índices especiais** no Firestore com scope `COLLECTION_GROUP`, não apenas `COLLECTION`.

Sem o índice apropriado, a query falhava silenciosamente retornando 0 resultados, mesmo com o broker existindo no banco.

## Diagnóstico

### 1. Verificar se o broker existe

Script criado em `scripts/check-broker.js`:

```javascript
const userRecord = await admin.auth().getUserByEmail(email);
const brokerRef = db.collection('tenants').doc(tenantId)
  .collection('brokers').doc(brokerId);
const brokerDoc = await brokerRef.get();

// Tentativa de CollectionGroup query
const brokersQuery = await db.collectionGroup('brokers')
  .where('firebase_uid', '==', userRecord.uid)
  .get();
```

**Resultado:**
- ✅ Broker existe em `/tenants/{id}/brokers/{broker_id}`
- ✅ Campo `firebase_uid` está preenchido corretamente
- ❌ CollectionGroup query falha com erro `9 FAILED_PRECONDITION`

### 2. Erro FAILED_PRECONDITION

```
Error: 9 FAILED_PRECONDITION
```

Este erro indica que **falta um índice no Firestore** para suportar a query.

## Solução

### 1. Adicionar índice ao firestore.indexes.json

```json
{
  "collectionGroup": "brokers",
  "queryScope": "COLLECTION_GROUP",
  "fields": [
    { "fieldPath": "firebase_uid", "order": "ASCENDING" }
  ]
}
```

### 2. Criar índice no Firebase Console

**Caminho:** Firebase Console → Firestore → Índices → Adicionar isenção

**Configuração:**
- **Código do conjunto:** `brokers`
- **Escopo da coleção:** **Grupo de coleções** (Collection group)
- **Campos:**
  - Campo: `firebase_uid`
  - Modo: Crescente (Ascending)

### 3. Aguardar criação do índice

O índice leva ~2-5 minutos para ser criado. Status visível em:
- Firebase Console → Firestore → Índices

### 4. Deploy via Firebase CLI (Alternativa)

```bash
firebase deploy --only firestore:indexes --project ecosistema-imob-dev
```

**Nota:** Requer autenticação via `firebase login`

## Verificação

Após o índice ser criado, testar o login:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "daniel.garcia@altatechsystems.com",
    "password": "$%&AltatechSystems$%&"
  }'
```

**Resposta esperada:**
```json
{
  "token": "eyJhbGc...",
  "tenant_id": "391b12f8-ebe4-426a-8c99-ec5a10b1f361",
  "broker": {
    "id": "73f624cc-2db1-4a2f-9a95-8b21abffc8d7",
    "email": "daniel.garcia@altatechsystems.com",
    "role": "admin"
  }
}
```

## Logs Úteis

Logs adicionados em `auth_handler.go:205`:

```go
log.Printf("Looking for broker with firebase_uid: %s", userRecord.UID)
// ...
log.Printf("No broker found for firebase_uid: %s (email: %s)",
    userRecord.UID, userRecord.Email)
```

Verificar logs do backend para debug:
```bash
# Windows
type backend-*.log | findstr "broker"

# Linux/Mac
grep "broker" backend-*.log
```

## Índices Firestore Necessários

### Para Autenticação
```json
{
  "collectionGroup": "brokers",
  "queryScope": "COLLECTION_GROUP",
  "fields": [
    { "fieldPath": "firebase_uid", "order": "ASCENDING" }
  ]
},
{
  "collectionGroup": "brokers",
  "queryScope": "COLLECTION_GROUP",
  "fields": [
    { "fieldPath": "email", "order": "ASCENDING" }
  ]
}
```

## Arquivos Modificados

1. **backend/internal/handlers/auth_handler.go** - Adicionados logs de debug
2. **firestore.indexes.json** - Adicionados índices COLLECTION_GROUP
3. **scripts/check-broker.js** - Script de diagnóstico criado

## Lições Aprendidas

1. **CollectionGroup queries sempre requerem índices explícitos**
   - Mesmo queries simples como `Where("field", "==", value)` precisam de índice
   - Scope deve ser `COLLECTION_GROUP`, não `COLLECTION`

2. **Erro FAILED_PRECONDITION = falta índice**
   - Tanto no SDK Node.js quanto no Go
   - Sempre verificar índices antes de usar CollectionGroup

3. **Índices levam tempo para serem criados**
   - Não são imediatos (2-5 minutos em média)
   - Console mostra status "Building" durante criação

4. **Ferramentas de diagnóstico são essenciais**
   - Scripts Node.js para testar queries diretamente
   - Logs detalhados no backend para troubleshooting

## Prevenção

Para evitar este problema em novos handlers:

1. **Sempre criar índices ANTES de usar CollectionGroup**
2. **Adicionar ao firestore.indexes.json durante desenvolvimento**
3. **Testar queries complexas em ambiente de dev primeiro**
4. **Adicionar logs detalhados em queries críticas**

## Referências

- [Firestore Collection Group Queries](https://firebase.google.com/docs/firestore/query-data/queries#collection-group-query)
- [Firestore Index Management](https://firebase.google.com/docs/firestore/query-data/indexing)
- [Error: 9 FAILED_PRECONDITION](https://firebase.google.com/docs/firestore/query-data/indexing#index_exemptions)

---

**Implementado por:** Claude Code
**Data de Resolução:** 22 de dezembro de 2025
**Tempo de Diagnóstico:** ~30 minutos
**Status Final:** ✅ Resolvido
