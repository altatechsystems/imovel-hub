# 🗑️ Wipe Database Script

Script para limpar completamente a base de dados antes de uma nova importação.

## ⚠️ ATENÇÃO

Este script **DELETA PERMANENTEMENTE** todos os dados de imóveis do tenant especificado.

**NÃO** use em produção!

---

## 🎯 O que será deletado

- ✅ **properties** - Todos os imóveis
- ✅ **canonical_listings** - Todos os anúncios canônicos
- ✅ **listing_references** - Todas as referências de anúncios

## 🛡️ O que será preservado

- ✅ **tenants** - Dados da empresa
- ✅ **tenants/{tenant}/users** - Usuários administrativos
- ✅ **tenants/{tenant}/brokers** - Corretores
- ✅ **tenants/{tenant}/user_invitations** - Convites pendentes
- ✅ **Storage** - Imagens (podem ser deletadas manualmente depois)

---

## 🚀 Como usar

### 1. Certifique-se de estar em DEV

```bash
# Verifique o tenant ID no script (linha 30):
# const TENANT_ID = 'bd71c02b-5fa5-43df-8b46-a1df2206f1ef';
```

### 2. Execute o script

```bash
cd scripts
node wipe-database.js
```

### 3. Aguarde 5 segundos

O script tem um delay de segurança. Você pode cancelar com `Ctrl+C`.

### 4. Acompanhe o progresso

```
🗑️  Deleting properties...
   Found 1234 properties
   Deleting batch of 500 documents...
   Deleting batch of 500 documents...
   Deleting batch of 234 documents...
   ✅ Properties deleted

🗑️  Deleting canonical_listings...
   Found 567 canonical listings
   ✅ Canonical listings deleted

🗑️  Deleting listing_references...
   Found 890 listing references
   ✅ Listing references deleted

✅ DATABASE WIPED SUCCESSFULLY!
```

---

## 📊 Após a limpeza

### Importar novos dados

**Opção 1: XML (Imovelweb)**
```bash
# No backend
go run cmd/importer/main.go -file="path/to/file.xml" -tenant="bd71c02b..."
```

**Opção 2: XLSX (Excel)**
```bash
# No backend
go run cmd/importer/main.go -file="path/to/file.xlsx" -tenant="bd71c02b..."
```

### Verificar importação

```bash
# Contar imóveis importados
node scripts/check-property-listing.js

# Ver duplicatas (se houver)
node scripts/check-duplicates.js
```

---

## 🐛 Solução de Problemas

### Erro: "Permission denied"

Certifique-se de que o arquivo `firebase-adminsdk.json` existe:
```bash
ls backend/config/firebase-adminsdk.json
```

### Erro: "Tenant ID not found"

Verifique o tenant ID correto:
```bash
node scripts/check-user.js
```

### Script trava em "Deleting batch..."

É normal! Firestore processa em batches de 500. Aguarde alguns minutos.

---

## 💡 Dicas

### Backup antes de limpar (opcional)

Firestore não tem comando de backup simples, mas você pode:

1. **Firestore Console** > Exportar dados
2. Ou confiar no backup automático do Firebase (últimas 24h)

### Limpar também as imagens

Se quiser deletar as fotos dos imóveis:

1. Acesse: [Firebase Console - Storage](https://console.firebase.google.com/)
2. Vá em `ecosistema-imob-dev.firebasestorage.app`
3. Navegue até `properties/`
4. Selecione tudo e delete

---

## 📝 Histórico de Uso

Registre aqui cada vez que executar o script:

```
| Data       | Usuário | Tenant ID             | Docs Deletados | Motivo                    |
|------------|---------|------------------------|----------------|---------------------------|
| 2026-01-11 | Daniel  | bd71c02b...            | 2.691          | Reimportar arquivos XML   |
```

---

## 🔒 Segurança

- ✅ Delay de 5 segundos para cancelamento
- ✅ Apenas ambiente DEV (tenant ID hardcoded)
- ✅ Preserva usuários e configurações
- ✅ Logs detalhados de cada operação

---

**Desenvolvido por**: Altatech Systems
**Data**: Janeiro 2026
**Versão**: 1.0.0
