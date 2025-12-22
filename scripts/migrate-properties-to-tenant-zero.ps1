# Script para migrar imóveis sem tenant_id para o Tenant Zero (ALTATECH Imóveis)
# Uso: powershell -ExecutionPolicy Bypass -File scripts/migrate-properties-to-tenant-zero.ps1 -TenantZeroId "tenant-id-aqui"

param(
    [Parameter(Mandatory=$true)]
    [string]$TenantZeroId
)

$ErrorActionPreference = "Stop"

Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "Migração de Imóveis para Tenant Zero" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Tenant Zero ID: $TenantZeroId" -ForegroundColor Gray
Write-Host ""

# Instruções para usar o Firebase CLI ou Admin SDK
Write-Host "⚠️  Este script requer acesso ao Firestore." -ForegroundColor Yellow
Write-Host ""
Write-Host "OPÇÃO 1: Via Firebase Console (Manual)" -ForegroundColor Cyan
Write-Host "-----------------------------------------" -ForegroundColor Gray
Write-Host "1. Acesse: https://console.firebase.google.com" -ForegroundColor White
Write-Host "2. Selecione seu projeto" -ForegroundColor White
Write-Host "3. Vá em Firestore Database" -ForegroundColor White
Write-Host "4. Execute a seguinte query composta:" -ForegroundColor White
Write-Host ""
Write-Host "   // Encontrar propriedades sem tenant_id ou com tenant vazio" -ForegroundColor Gray
Write-Host "   Collection: properties" -ForegroundColor White
Write-Host "   Where: tenant_id == null OR tenant_id == ''" -ForegroundColor White
Write-Host ""
Write-Host "5. Para cada documento encontrado:" -ForegroundColor White
Write-Host "   - Adicionar campo: tenant_id = '$TenantZeroId'" -ForegroundColor White
Write-Host ""

Write-Host "OPÇÃO 2: Via Script Node.js (Automático)" -ForegroundColor Cyan
Write-Host "-----------------------------------------" -ForegroundColor Gray
Write-Host ""

# Criar script Node.js para migração
$nodeScript = @"
// migrate-properties.js
// Uso: node scripts/migrate-properties.js <TENANT_ZERO_ID>

const admin = require('firebase-admin');

// Inicializar Firebase Admin SDK
const serviceAccount = require('../backend/config/firebase-adminsdk.json');

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount)
});

const db = admin.firestore();

async function migrateProperties(tenantZeroId) {
  console.log('🔄 Iniciando migração de imóveis...');
  console.log('Tenant Zero ID:', tenantZeroId);
  console.log('');

  try {
    // Buscar todas as properties (coleção raiz - estrutura antiga)
    const propertiesSnapshot = await db.collection('properties').get();

    let migratedCount = 0;
    let errorCount = 0;

    console.log('Total de imóveis encontrados:', propertiesSnapshot.size);
    console.log('');

    // Processar cada property
    for (const doc of propertiesSnapshot.docs) {
      const propertyData = doc.data();
      const propertyId = doc.id;

      try {
        // 1. Criar property na nova estrutura multi-tenant
        const newPropertyRef = db
          .collection('tenants')
          .doc(tenantZeroId)
          .collection('properties')
          .doc(propertyId);

        // 2. Copiar dados e adicionar tenant_id
        const newPropertyData = {
          ...propertyData,
          tenant_id: tenantZeroId,
          migrated_at: admin.firestore.FieldValue.serverTimestamp(),
          migrated_from: 'root_collection'
        };

        await newPropertyRef.set(newPropertyData);

        console.log(\`✅ Migrado: \${propertyId} - \${propertyData.address?.street || 'Sem endereço'}\`);
        migratedCount++;

        // 3. (Opcional) Deletar da coleção raiz após confirmação
        // await doc.ref.delete();

      } catch (error) {
        console.error(\`❌ Erro ao migrar \${propertyId}:\`, error.message);
        errorCount++;
      }
    }

    console.log('');
    console.log('====================================');
    console.log('✅ Migração Concluída!');
    console.log('====================================');
    console.log(\`Total migrado: \${migratedCount}\`);
    console.log(\`Erros: \${errorCount}\`);
    console.log('');

  } catch (error) {
    console.error('❌ Erro na migração:', error);
    process.exit(1);
  }
}

// Executar migração
const tenantZeroId = process.argv[2];

if (!tenantZeroId) {
  console.error('❌ Erro: Tenant Zero ID não fornecido');
  console.log('Uso: node scripts/migrate-properties.js <TENANT_ZERO_ID>');
  process.exit(1);
}

migrateProperties(tenantZeroId)
  .then(() => {
    console.log('✅ Script finalizado com sucesso');
    process.exit(0);
  })
  .catch((error) => {
    console.error('❌ Erro fatal:', error);
    process.exit(1);
  });
"@

# Salvar script Node.js
$nodeScriptPath = "scripts/migrate-properties.js"
Set-Content -Path $nodeScriptPath -Value $nodeScript -Encoding UTF8

Write-Host "✅ Script Node.js criado em: $nodeScriptPath" -ForegroundColor Green
Write-Host ""
Write-Host "Para executar a migração automática:" -ForegroundColor Cyan
Write-Host "1. Instalar firebase-admin:" -ForegroundColor White
Write-Host "   npm install firebase-admin" -ForegroundColor Gray
Write-Host ""
Write-Host "2. Garantir que existe o arquivo de credenciais:" -ForegroundColor White
Write-Host "   backend/config/firebase-adminsdk.json" -ForegroundColor Gray
Write-Host ""
Write-Host "3. Executar o script:" -ForegroundColor White
Write-Host "   node $nodeScriptPath $TenantZeroId" -ForegroundColor Gray
Write-Host ""

Write-Host "OPÇÃO 3: Via Backend Go (Endpoint Administrativo)" -ForegroundColor Cyan
Write-Host "---------------------------------------------------" -ForegroundColor Gray
Write-Host ""
Write-Host "Criar endpoint temporário no backend:" -ForegroundColor White
Write-Host "POST /api/v1/admin/migrate-properties-to-tenant" -ForegroundColor Gray
Write-Host ""
Write-Host "Body:" -ForegroundColor White
Write-Host @"
{
  "tenant_id": "$TenantZeroId",
  "dry_run": false
}
"@ -ForegroundColor Gray
Write-Host ""

Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "⚠️  IMPORTANTE" -ForegroundColor Yellow
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Antes de executar a migração:" -ForegroundColor Yellow
Write-Host "1. ✅ Fazer backup do Firestore" -ForegroundColor White
Write-Host "2. ✅ Testar com dry_run=true primeiro" -ForegroundColor White
Write-Host "3. ✅ Verificar se o Tenant Zero foi criado corretamente" -ForegroundColor White
Write-Host "4. ✅ Confirmar estrutura de dados das properties" -ForegroundColor White
Write-Host ""
