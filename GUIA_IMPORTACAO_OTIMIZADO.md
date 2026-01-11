# 📦 Guia Rápido - Importação Otimizada

## ✅ Problema Resolvido

O sistema agora processa arquivos XML/XLSX grandes **sem travar** o computador.

**Otimizações implementadas**:
- ✅ **Streaming XML** - não carrega tudo na memória
- ✅ **Batching** - processa 50 imóveis por vez
- ✅ **Concurrency Control** - máximo 3 goroutines simultâneas
- ✅ **Pausas automáticas** - permite Garbage Collector limpar memória

---

## 🚀 Como Usar

### 1. Limpar Base de Dados (Opcional)

Se quiser começar do zero:

```bash
cd scripts
node preview-wipe.js      # Ver o que será deletado
node wipe-database.js     # Deletar tudo
```

### 2. Iniciar Backend

```bash
cd backend
go run cmd/server/main.go
```

### 3. Fazer Importação

1. Acesse: http://localhost:3002/dashboard/importacao
2. Selecione arquivo XML (e XLS se tiver)
3. Clique em **"Iniciar Importação"**
4. Aguarde o processamento

### 4. Acompanhar Progresso

No terminal do backend, você verá:

```
📦 Processing batch 1-50 of 1234 properties
✅ Batch complete: 50/1234 properties processed
⏸️  Pausing 2s between batches to allow memory cleanup...

📦 Processing batch 51-100 of 1234 properties
✅ Batch complete: 100/1234 properties processed
```

No frontend, você verá:
- Status: "Importando Imóveis..."
- Progress bar animado
- Batch ID para tracking

---

## ⚙️ Configurações

### Servidor Fraco (2-4GB RAM)

Edite [import_handler.go:230-231](backend/internal/handlers/import_handler.go#L230-L231):

```go
const batchSize = 25     // Menor para economizar memória
const maxWorkers = 2     // Menos goroutines
```

### Servidor Normal (8-16GB RAM)

**Configuração padrão** (já implementada):
```go
const batchSize = 50
const maxWorkers = 3
```

### Servidor Potente (16GB+ RAM)

```go
const batchSize = 100    // Processa mais por vez
const maxWorkers = 5     // Mais concorrência
```

---

## 📊 Consumo de Memória Esperado

| Arquivo | Imóveis | Antes | Depois | Melhoria |
|---------|---------|-------|--------|----------|
| Pequeno | 100 | ~400MB | ~100MB | -75% |
| Médio | 500 | ~2GB | ~400MB | -80% |
| Grande | 1.000 | ~4GB+ (trava) | ~800MB | **-80%** |
| Muito Grande | 2.000+ | ❌ Trava | ~1.5GB | **Funciona!** |

---

## 🔍 Monitorar Uso de Memória

### Windows (PowerShell)

```powershell
# Abrir novo terminal e rodar:
while ($true) {
    $proc = Get-Process -Name "server" -ErrorAction SilentlyContinue
    if ($proc) {
        Write-Host "Memory: $([math]::Round($proc.WS/1MB,2)) MB" -ForegroundColor Green
    }
    Start-Sleep -Seconds 2
}
```

### Linux/Mac

```bash
watch -n 2 "ps aux | grep 'server' | grep -v grep | awk '{print \"Memory:\", \$6/1024, \"MB\"}'"
```

---

## ⏱️ Tempo Estimado de Importação

| Imóveis | Tempo Esperado | Observação |
|---------|----------------|------------|
| 100 | ~2 minutos | Rápido |
| 500 | ~10 minutos | Normal |
| 1.000 | ~20 minutos | Com pausas para GC |
| 2.000+ | ~40 minutos | Pode demorar, mas não trava! |

**Nota**: Tempo inclui:
- Parse XML/XLS
- Deduplicação
- Criação no Firestore
- Download de fotos (se habilitado)
- Pausas entre batches

---

## ❌ Troubleshooting

### "Ainda está usando muita memória"

1. **Reduza batch e workers**:
   ```go
   const batchSize = 25
   const maxWorkers = 2
   ```

2. **Aumente pausa entre batches**:
   ```go
   time.Sleep(5 * time.Second)  // De 2s para 5s
   ```

3. **Reinicie o backend** para limpar memória acumulada

### "Importação muito lenta"

1. **Aumente workers** (se tiver RAM):
   ```go
   const maxWorkers = 5
   ```

2. **Desabilite pausas** (apenas se não travar):
   ```go
   // time.Sleep(2 * time.Second)  // Comentar
   ```

### "Backend parou de responder"

1. Verifique logs do backend
2. Reinicie o servidor: `Ctrl+C` e rodar `go run cmd/server/main.go`
3. O batch continua de onde parou (se Firestore já salvou)

### "Frontend mostra erro 401"

1. Faça logout e login novamente
2. Token expirou durante importação longa
3. Backend precisa refresh automático de token

---

## 📋 Checklist Pré-Importação

Antes de importar arquivos grandes:

- [ ] Backend rodando (`go run cmd/server/main.go`)
- [ ] Frontend acessível (http://localhost:3002)
- [ ] Arquivos XML/XLS preparados
- [ ] Espaço em disco suficiente (Storage do Firebase)
- [ ] Verificar configuração de batch/workers
- [ ] Ter ~1-2GB RAM livre no sistema
- [ ] Não rodar outras tarefas pesadas durante importação

---

## 🎯 Próximos Passos

Depois da importação:

1. **Verificar resultados**:
   - Acesse: http://localhost:3002/dashboard
   - Veja imóveis importados

2. **Checar duplicatas** (se houver):
   ```bash
   node scripts/check-duplicates.js
   ```

3. **Limpar duplicatas** (se necessário):
   ```bash
   node scripts/cleanup-duplicates.js
   ```

4. **Verificar fotos**:
   - Acesse Firebase Console > Storage
   - Veja pasta `properties/`

---

## 📞 Suporte

Se encontrar problemas:

1. **Verifique logs do backend** (terminal onde rodou `go run`)
2. **Veja documentação técnica**: [MEMORY_OPTIMIZATION.md](backend/MEMORY_OPTIMIZATION.md)
3. **Abra issue** no repositório com logs e descrição

---

**Última atualização**: 2026-01-11
**Versão**: 1.0.0
