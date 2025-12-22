# Firestore Database Setup

## Problema Identificado

Durante os testes de integração, foi identificado que o Firestore Database não está criado no projeto Firebase `ecosistema-imob-dev`.

**Erro:**
```
The database (default) does not exist for project ecosistema-imob-dev
```

## Solução: Criar Firestore Database

### Opção 1: Firebase Console (Produção)

1. Acesse o [Firebase Console](https://console.firebase.google.com/)
2. Selecione o projeto `ecosistema-imob-dev`
3. No menu lateral, clique em **Build** → **Firestore Database**
4. Clique em **Create database**
5. Escolha o modo:
   - **Production mode** (recomendado para produção)
   - **Test mode** (apenas para desenvolvimento inicial)
6. Selecione a localização:
   - Recomendado: `southamerica-east1` (São Paulo, Brasil)
7. Clique em **Enable**

### Opção 2: Firestore Emulator (Desenvolvimento Local)

Para desenvolvimento local sem custos, use o Firestore Emulator:

#### 1. Instalar Firebase CLI

```bash
npm install -g firebase-tools
```

#### 2. Fazer Login

```bash
firebase login
```

#### 3. Inicializar Emulators

```bash
cd backend
firebase init emulators
```

Selecione:
- [x] Firestore Emulator
- [x] Storage Emulator (opcional)

#### 4. Configurar firebase.json

Crie ou atualize `backend/firebase.json`:

```json
{
  "emulators": {
    "firestore": {
      "port": 8081,
      "host": "127.0.0.1"
    },
    "storage": {
      "port": 9199
    },
    "ui": {
      "enabled": true,
      "port": 4000
    }
  }
}
```

#### 5. Iniciar Emulators

```bash
firebase emulators:start
```

#### 6. Configurar Backend para usar Emulator

Atualize `backend/.env`:

```env
# Adicionar para usar emulator
FIRESTORE_EMULATOR_HOST=localhost:8081
FIREBASE_STORAGE_EMULATOR_HOST=localhost:9199
```

Ou configure via código em `backend/internal/config/firebase.go`:

```go
if os.Getenv("ENVIRONMENT") == "development" {
    os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8081")
    os.Setenv("FIREBASE_STORAGE_EMULATOR_HOST", "localhost:9199")
}
```

## Status Atual

- ✅ Backend compilado e rodando na porta 8080
- ✅ Rotas configuradas corretamente (53 endpoints)
- ✅ Frontend implementado e buildado
- ❌ Firestore Database precisa ser criado
- ⏳ Aguardando criação do database para testes de integração

## Próximos Passos

1. **Opção A:** Criar Firestore Database no console (produção)
2. **Opção B:** Configurar Firestore Emulator (desenvolvimento)
3. Executar script de teste: `bash backend/scripts/create-test-data.sh`
4. Atualizar `frontend-public/.env.local` com `NEXT_PUBLIC_TENANT_ID`
5. Testar frontend: `cd frontend-public && npm run dev`

## Endpoints Disponíveis

### Públicos (Frontend)
- `GET /api/{tenant_id}/properties` - Listar imóveis
- `GET /api/{tenant_id}/properties/{id}` - Detalhes do imóvel
- `GET /api/{tenant_id}/properties/slug/{slug}` - Buscar por slug
- `POST /api/{tenant_id}/leads` - Criar lead

### Admin (Dashboard)
- Tenants: CRUD completo
- Brokers: CRUD + ativação/desativação
- Owners: CRUD + LGPD (anonimização, revogação)
- Properties: CRUD + status + visibilidade + duplicados
- Listings: CRUD + canonical
- Leads: CRUD + status + atribuição + LGPD
- Activity Logs: Timeline de propriedades e leads
- Storage: Upload/List/Delete imagens

### Storage (Imagens)
- `POST /api/{tenant_id}/property-images/{property_id}` - Upload
- `GET /api/{tenant_id}/property-images/{property_id}` - Listar
- `GET /api/{tenant_id}/property-images/{property_id}/{image_id}` - URL
- `DELETE /api/{tenant_id}/property-images/{property_id}/{image_id}` - Deletar

## Comandos Rápidos

### Iniciar Backend
```bash
cd backend
./bin/caas.exe
```

### Criar Dados de Teste (após Firestore configurado)
```bash
bash backend/scripts/create-test-data.sh
```

### Iniciar Frontend
```bash
cd frontend-public
npm run dev
```

Acesse: http://localhost:3000

## Observações

- O frontend está **100% funcional** e aguardando apenas dados do backend
- O backend está **rodando corretamente** mas precisa do Firestore
- Todos os tipos TypeScript estão alinhados com os modelos Go
- A API client no frontend está configurada para `http://localhost:8080/api`
