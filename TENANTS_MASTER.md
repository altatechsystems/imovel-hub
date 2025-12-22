# Tenants Master da Plataforma

## 📋 Visão Geral

Este documento descreve a estrutura de **tenants especiais** da plataforma, que possuem privilégios diferenciados dos tenants comuns (imobiliárias clientes).

## 🏢 Estrutura de Tenants

```
┌─────────────────────────────────────────────────────────────┐
│                    PLATAFORMA                                │
│  (Multi-tenant com hierarquia de privilégios)               │
└─────────────────────────────────────────────────────────────┘
                            │
                ┌───────────┴───────────┐
                │                       │
┌───────────────▼─────────┐  ┌─────────▼──────────────┐
│   TENANT MASTER          │  │   TENANT ZERO          │
│   ALTATECH Systems       │  │   ALTATECH Imóveis     │
│                          │  │                        │
│   - Platform Admin       │  │   - Tenant Padrão      │
│   - Acesso Full          │  │   - Proprietária dos   │
│   - Desenvolvimento      │  │     imóveis importados │
│   - Suporte              │  │   - CRECI 5733-J       │
│                          │  │                        │
└──────────────────────────┘  └────────────────────────┘
                            │
                ┌───────────┴───────────┐
                │                       │
┌───────────────▼─────────┐  ┌─────────▼──────────────┐
│   TENANT CLIENTE 1       │  │   TENANT CLIENTE N     │
│   Imobiliária ABC        │  │   Imobiliária XYZ      │
│                          │  │                        │
│   - Acesso normal        │  │   - Acesso normal      │
│   - Seus imóveis         │  │   - Seus imóveis       │
│   - Seus corretores      │  │   - Seus corretores    │
│                          │  │                        │
└──────────────────────────┘  └────────────────────────┘
```

## 🔐 Tenant Master - ALTATECH Systems

### Dados do Tenant

```yaml
Nome: ALTATECH Systems
CNPJ: 36.077.869/0001-81
Slug: altatech-systems
Status: active
is_platform_admin: true  # Campo especial
```

### Usuário Admin

```yaml
Nome: Daniel Garcia
Email: daniel.garcia@altatechsystems.com
Senha: $%&AltatechSystems$%&
Phone: +5511999999999
Role: admin
Platform Admin: true  # Custom claim adicional
```

### Privilégios Especiais

✅ **Acesso total a todos os tenants**
- Visualizar dados de qualquer tenant
- Modificar configurações de qualquer tenant
- Acessar métricas globais da plataforma

✅ **Ferramentas de desenvolvimento**
- Console de debug
- Logs de sistema
- Métricas de performance

✅ **Gerenciamento de tenants**
- Criar/suspender/reativar tenants
- Modificar planos e limites
- Configurar whitelabel

✅ **Suporte técnico**
- Acessar dashboard de qualquer cliente
- Resolver problemas de importação
- Gerenciar leads de clientes (somente leitura)

### Custom Claims (Firebase)

```json
{
  "tenant_id": "altatech-systems-xxxxx",
  "role": "admin",
  "is_platform_admin": true,
  "permissions": [
    "access_all_tenants",
    "manage_tenants",
    "view_system_logs",
    "debug_mode"
  ]
}
```

### Firestore Structure

```
/tenants/altatech-systems-xxxxx
  - name: "ALTATECH Systems"
  - slug: "altatech-systems"
  - cnpj: "36.077.869/0001-81"
  - status: "active"
  - is_platform_admin: true  ⭐
  - created_at: timestamp
  - settings:
      - business_name: "ALTATECH Systems"
      - logo_url: ""
      - whatsapp_default: "+5511999999999"

/tenants/altatech-systems-xxxxx/brokers/broker-xxxxx
  - user_id: Firebase UID
  - email: "daniel.garcia@altatechsystems.com"
  - name: "Daniel Garcia"
  - phone: "+5511999999999"
  - role: "admin"
  - is_platform_admin: true  ⭐
  - status: "active"
```

## 🏠 Tenant Zero - ALTATECH Imóveis

### Dados do Tenant

```yaml
Nome: ALTATECH Imóveis
CNPJ: 26.517.873/0001-60
CRECI: 5733-J
Slug: altatech-imoveis
Status: active
is_default_tenant: true  # Campo especial
```

### Usuário Admin

```yaml
Nome: Administração
Email: administracao@altatechimoveis.com
Senha: $%&AltatechImoveis$%&
Phone: +5511988888888
Role: admin
```

### Responsabilidades

✅ **Proprietária dos imóveis importados**
- Todos os imóveis carregados via XML/XLS pertencem a este tenant
- CRECI vinculado: 5733-J
- Responsável legal pelos anúncios

✅ **Gestão de co-corretagem**
- Aprovar/rejeitar solicitações de selling_brokers
- Gerenciar comissões
- Controlar visibilidade dos imóveis

✅ **Operações normais**
- Gerenciar corretores
- Receber e distribuir leads
- Importar novos imóveis

### Custom Claims (Firebase)

```json
{
  "tenant_id": "altatech-imoveis-xxxxx",
  "role": "admin",
  "is_default_tenant": true,
  "creci": "5733-J"
}
```

### Firestore Structure

```
/tenants/altatech-imoveis-xxxxx
  - name: "ALTATECH Imóveis"
  - slug: "altatech-imoveis"
  - cnpj: "26.517.873/0001-60"
  - creci: "5733-J"
  - status: "active"
  - is_default_tenant: true  ⭐
  - created_at: timestamp
  - settings:
      - business_name: "ALTATECH Imóveis"
      - logo_url: ""
      - whatsapp_default: "+5511988888888"

/tenants/altatech-imoveis-xxxxx/brokers/broker-xxxxx
  - user_id: Firebase UID
  - email: "administracao@altatechimoveis.com"
  - name: "Administração"
  - phone: "+5511988888888"
  - role: "admin"
  - status: "active"

/tenants/altatech-imoveis-xxxxx/properties/...
  - [Todos os imóveis importados]
```

## 🔄 Migração de Imóveis Existentes

### Problema

Imóveis carregados antes da implementação multi-tenant estão na coleção raiz `/properties` sem `tenant_id`.

### Solução

Migrar todos os imóveis para `/tenants/altatech-imoveis-xxxxx/properties/`

### Scripts Disponíveis

**1. Script PowerShell de Migração**
```powershell
powershell -ExecutionPolicy Bypass -File scripts/migrate-properties-to-tenant-zero.ps1 -TenantZeroId "tenant-id-aqui"
```

**2. Script Node.js (Automático)**
```bash
node scripts/migrate-properties.js <TENANT_ZERO_ID>
```

**3. Via Firebase Console (Manual)**
- Acessar Firestore Database
- Copiar documentos de `/properties` para `/tenants/{tenantZeroId}/properties`
- Adicionar campo `tenant_id` em cada documento

### Campos Adicionados na Migração

```json
{
  "tenant_id": "altatech-imoveis-xxxxx",
  "migrated_at": "2025-12-22T...",
  "migrated_from": "root_collection"
}
```

## 🚀 Como Criar os Tenants Master

### Pré-requisitos

1. ✅ Backend rodando na porta 3000
2. ✅ Firebase Auth configurado
3. ✅ Firestore habilitado

### Passo 1: Executar Script de Criação

```powershell
powershell -ExecutionPolicy Bypass -File scripts/create-master-tenants.ps1
```

Este script irá:
- ✅ Criar TENANT MASTER (ALTATECH Systems)
- ✅ Criar TENANT ZERO (ALTATECH Imóveis)
- ✅ Criar usuários admin para cada tenant
- ⚠️ Mostrar instruções para configuração manual

### Passo 2: Configurar Metadados no Firestore

**Via Firebase Console:**

1. Acessar: https://console.firebase.google.com
2. Selecionar projeto
3. Ir em Firestore Database

**Para TENANT MASTER:**
```
Collection: tenants
Document: <tenant-id-do-master>

Adicionar campos:
  is_platform_admin: true (boolean)
  cnpj: "36.077.869/0001-81" (string)
```

**Para TENANT ZERO:**
```
Collection: tenants
Document: <tenant-id-do-zero>

Adicionar campos:
  is_default_tenant: true (boolean)
  cnpj: "26.517.873/0001-60" (string)
  creci: "5733-J" (string)
```

### Passo 3: Configurar Custom Claims

**Via Firebase Console > Authentication:**

1. Selecionar usuário `daniel.garcia@altatechsystems.com`
2. Adicionar Custom Claims:
```json
{
  "is_platform_admin": true,
  "permissions": ["access_all_tenants", "manage_tenants"]
}
```

Ou via Firebase Admin SDK (backend):
```go
client.SetCustomUserClaims(ctx, uid, map[string]interface{}{
    "tenant_id": tenantId,
    "role": "admin",
    "is_platform_admin": true,
})
```

### Passo 4: Migrar Imóveis

```bash
# 1. Instalar dependências
npm install firebase-admin

# 2. Executar migração
node scripts/migrate-properties.js <TENANT_ZERO_ID>
```

## 🔒 Segurança e Middleware

### Backend - Validação de Platform Admin

```go
func PlatformAdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        claims := c.MustGet("claims").(map[string]interface{})
        isPlatformAdmin, _ := claims["is_platform_admin"].(bool)

        if !isPlatformAdmin {
            c.JSON(403, gin.H{"error": "platform admin access required"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### Rotas Protegidas

```go
// Rotas de Platform Admin
platformAdmin := router.Group("/api/v1/platform")
platformAdmin.Use(AuthMiddleware())
platformAdmin.Use(PlatformAdminMiddleware())
{
    platformAdmin.GET("/tenants", handlers.ListAllTenants)
    platformAdmin.GET("/tenants/:tenantId", handlers.GetTenantDetails)
    platformAdmin.PATCH("/tenants/:tenantId/status", handlers.UpdateTenantStatus)
    platformAdmin.GET("/metrics/global", handlers.GetGlobalMetrics)
    platformAdmin.GET("/logs", handlers.GetSystemLogs)
}
```

## 📊 Diferenças Entre Tenants

| Característica | Tenant Master | Tenant Zero | Tenant Cliente |
|----------------|---------------|-------------|----------------|
| **is_platform_admin** | ✅ true | ❌ false | ❌ false |
| **is_default_tenant** | ❌ false | ✅ true | ❌ false |
| **Acesso a outros tenants** | ✅ Sim | ❌ Não | ❌ Não |
| **Gerenciar plataforma** | ✅ Sim | ❌ Não | ❌ Não |
| **Imóveis próprios** | ❌ Não | ✅ Sim | ✅ Sim |
| **Recebe leads** | ❌ Não | ✅ Sim | ✅ Sim |
| **Co-corretagem** | ❌ N/A | ✅ Sim | ✅ Sim |
| **CRECI obrigatório** | ❌ Não | ✅ Sim | ⚠️ Opcional |

## 📝 Credenciais de Acesso

### TENANT MASTER - ALTATECH Systems

```
URL: http://localhost:3002/login
Email: daniel.garcia@altatechsystems.com
Senha: $%&AltatechSystems$%&
```

### TENANT ZERO - ALTATECH Imóveis

```
URL: http://localhost:3002/login
Email: administracao@altatechimoveis.com
Senha: $%&AltatechImoveis$%&
```

## ⚠️ Importante

### Segurança

- ✅ Nunca commitar senhas no Git
- ✅ Usar variáveis de ambiente em produção
- ✅ Habilitar MFA para platform admin
- ✅ Rotacionar senhas periodicamente

### Backup

- ✅ Backup diário do Firestore
- ✅ Backup de custom claims
- ✅ Documentar qualquer mudança de privilégios

### Auditoria

- ✅ Registrar todas as ações do platform admin em ActivityLog
- ✅ Monitorar acessos cross-tenant
- ✅ Alertas para mudanças de status de tenants

---

**Última atualização**: 2025-12-22
**Status**: 📝 Documentado (Aguardando Criação)
