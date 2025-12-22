# 🚀 Instruções para Criar Tenants Master

## ✅ Pré-requisitos

Antes de executar, certifique-se de que:

1. **Backend está rodando** na porta 3000
   ```bash
   cd backend
   go run cmd/api/main.go
   ```

2. **Firebase Auth está configurado** no backend
   - Credenciais em `backend/config/firebase-adminsdk.json`
   - Variáveis de ambiente configuradas

3. **Firestore está habilitado** no projeto Firebase

## 📝 Passo a Passo Completo

### PASSO 1: Iniciar Backend

```bash
cd backend
go run cmd/api/main.go
```

Aguarde até ver:
```
✓ Server running on :3000
✓ Firebase initialized
✓ Firestore connected
```

### PASSO 2: Executar Script de Criação

Abra um **novo terminal** e execute:

```powershell
powershell -ExecutionPolicy Bypass -File scripts/create-master-tenants.ps1
```

O script irá:
- ✅ Criar TENANT MASTER (ALTATECH Systems)
- ✅ Criar TENANT ZERO (ALTATECH Imóveis)
- ✅ Exibir os IDs dos tenants criados
- ⚠️ Mostrar instruções para próximos passos

**Saída esperada:**
```
=====================================
Criando Tenants Master da Plataforma
=====================================

1. Criando Tenant Master: ALTATECH Systems...
✅ Tenant Master criado com sucesso!
   Tenant ID: tenant-abc123xyz
   Broker ID: broker-xyz789
   Email: daniel.garcia@altatechsystems.com

2. Criando Tenant Zero: ALTATECH Imóveis...
✅ Tenant Zero criado com sucesso!
   Tenant ID: tenant-def456uvw
   Broker ID: broker-uvw123
   Email: administracao@altatechimoveis.com

=====================================
✅ RESUMO DA CRIAÇÃO
=====================================

🏢 TENANT MASTER - ALTATECH Systems
   Email: daniel.garcia@altatechsystems.com
   Senha: $%&AltatechSystems$%&
   Tenant ID: tenant-abc123xyz
   Login Admin: http://localhost:3002/login

🏢 TENANT ZERO - ALTATECH Imóveis
   Email: administracao@altatechimoveis.com
   Senha: $%&AltatechImoveis$%&
   Tenant ID: tenant-def456uvw
   Login Admin: http://localhost:3002/login
```

### PASSO 3: Anotar os IDs dos Tenants

**IMPORTANTE:** Copie os IDs exibidos no terminal:

```
TENANT_MASTER_ID=tenant-abc123xyz
TENANT_ZERO_ID=tenant-def456uvw
```

Você precisará deles para os próximos passos.

### PASSO 4: Configurar Metadados no Firestore

Acesse o Firebase Console:
https://console.firebase.google.com

#### A) Configurar Tenant Master

1. Vá em **Firestore Database**
2. Navegue até: `tenants` > `[TENANT_MASTER_ID]`
3. Clique em **"Add field"**
4. Adicionar os seguintes campos:

```
Campo: is_platform_admin
Tipo: boolean
Valor: true

Campo: cnpj
Tipo: string
Valor: 36.077.869/0001-81
```

5. Clicar em **Save**

#### B) Configurar Tenant Zero

1. Navegue até: `tenants` > `[TENANT_ZERO_ID]`
2. Adicionar os seguintes campos:

```
Campo: is_default_tenant
Tipo: boolean
Valor: true

Campo: cnpj
Tipo: string
Valor: 26.517.873/0001-60

Campo: creci
Tipo: string
Valor: 5733-J
```

3. Clicar em **Save**

### PASSO 5: Configurar Custom Claims (Platform Admin)

#### Opção A: Via Firebase Console (Mais Fácil)

1. No Firebase Console, vá em **Authentication**
2. Encontre o usuário: `daniel.garcia@altatechsystems.com`
3. Clique no usuário
4. Na aba **"Custom claims"**, adicione:

```json
{
  "is_platform_admin": true,
  "permissions": ["access_all_tenants", "manage_tenants", "debug_mode"]
}
```

5. Salvar

#### Opção B: Via Backend (Código)

Criar endpoint temporário no backend ou executar via script:

```go
// backend/cmd/set-platform-admin/main.go
package main

import (
    "context"
    "log"

    firebase "firebase.google.com/go"
    "firebase.google.com/go/auth"
    "google.golang.org/api/option"
)

func main() {
    ctx := context.Background()

    // Inicializar Firebase
    sa := option.WithCredentialsFile("backend/config/firebase-adminsdk.json")
    app, err := firebase.NewApp(ctx, nil, sa)
    if err != nil {
        log.Fatalf("error initializing app: %v\n", err)
    }

    client, err := app.Auth(ctx)
    if err != nil {
        log.Fatalf("error getting Auth client: %v\n", err)
    }

    // Obter usuário por email
    user, err := client.GetUserByEmail(ctx, "daniel.garcia@altatechsystems.com")
    if err != nil {
        log.Fatalf("error getting user: %v\n", err)
    }

    // Setar custom claims
    claims := map[string]interface{}{
        "tenant_id": "SEU_TENANT_MASTER_ID_AQUI",
        "role": "admin",
        "is_platform_admin": true,
        "permissions": []string{"access_all_tenants", "manage_tenants"},
    }

    err = client.SetCustomUserClaims(ctx, user.UID, claims)
    if err != nil {
        log.Fatalf("error setting custom claims: %v\n", err)
    }

    log.Println("✅ Custom claims set successfully for platform admin!")
}
```

Executar:
```bash
cd backend
go run cmd/set-platform-admin/main.go
```

### PASSO 6: Migrar Imóveis para Tenant Zero

#### A) Preparar Ambiente

```bash
# Instalar firebase-admin (se ainda não tiver)
npm install firebase-admin
```

#### B) Executar Migração

```bash
node scripts/migrate-properties.js <TENANT_ZERO_ID>
```

Substitua `<TENANT_ZERO_ID>` pelo ID real obtido no PASSO 3.

**Exemplo:**
```bash
node scripts/migrate-properties.js tenant-def456uvw
```

#### C) Verificar Migração

1. Acessar Firestore Console
2. Navegar até: `tenants` > `[TENANT_ZERO_ID]` > `properties`
3. Verificar se os imóveis foram migrados
4. Conferir se todos têm o campo `tenant_id` preenchido

### PASSO 7: Testar Login dos Tenants

#### A) Testar Tenant Master

1. Acessar: http://localhost:3002/login
2. Login:
   ```
   Email: daniel.garcia@altatechsystems.com
   Senha: $%&AltatechSystems$%&
   ```
3. Verificar acesso ao dashboard
4. Verificar se custom claims estão presentes (F12 > Console):
   ```javascript
   firebase.auth().currentUser.getIdTokenResult()
     .then(token => console.log(token.claims))
   ```

#### B) Testar Tenant Zero

1. Fazer logout
2. Login:
   ```
   Email: administracao@altatechimoveis.com
   Senha: $%&AltatechImoveis$%&
   ```
3. Verificar acesso ao dashboard
4. Verificar se os imóveis aparecem na listagem

### PASSO 8: Fazer Backup

```bash
# Exportar Firestore
firebase firestore:export backup_$(date +%Y%m%d)

# Ou via gcloud
gcloud firestore export gs://[BUCKET_NAME]/backup_$(date +%Y%m%d)
```

## 🔍 Troubleshooting

### Erro: "Email já cadastrado"

**Problema:** Tentando criar tenant que já existe

**Solução:**
1. O script automaticamente tentará fazer login
2. Se falhar, delete o usuário no Firebase Console > Authentication
3. Execute o script novamente

### Erro: "Connection refused" ao criar tenant

**Problema:** Backend não está rodando

**Solução:**
```bash
cd backend
go run cmd/api/main.go
```

### Erro: "Firebase not initialized"

**Problema:** Credenciais do Firebase não configuradas

**Solução:**
1. Verificar se existe: `backend/config/firebase-adminsdk.json`
2. Baixar do Firebase Console se necessário
3. Verificar variáveis de ambiente

### Custom Claims não aparecem

**Problema:** Token não foi renovado

**Solução:**
1. Fazer logout do Firebase Auth
2. Fazer login novamente
3. Token será renovado com novos claims

### Imóveis não aparecem após migração

**Problema:** Tenant ID incorreto ou migração falhou

**Solução:**
1. Verificar logs da migração
2. Conferir Firestore manualmente
3. Executar migração novamente (é idempotente)

## 📋 Checklist Final

Antes de considerar concluído, verificar:

- [ ] TENANT MASTER criado no Firestore
- [ ] TENANT ZERO criado no Firestore
- [ ] Metadados `is_platform_admin` adicionados ao TENANT MASTER
- [ ] Metadados `is_default_tenant` e `creci` adicionados ao TENANT ZERO
- [ ] Custom claims configurados para daniel.garcia@altatechsystems.com
- [ ] Imóveis migrados para `/tenants/[TENANT_ZERO_ID]/properties`
- [ ] Login testado para ambos os tenants
- [ ] Dashboard acessível para ambos os tenants
- [ ] Backup do Firestore realizado
- [ ] IDs dos tenants documentados em local seguro

## 📞 Suporte

Se encontrar problemas:

1. Verificar logs do backend
2. Verificar console do navegador (F12)
3. Conferir Firestore no Firebase Console
4. Revisar documentação em [TENANTS_MASTER.md](TENANTS_MASTER.md)

---

**Status**: 📝 Aguardando Execução
**Última atualização**: 2025-12-22
