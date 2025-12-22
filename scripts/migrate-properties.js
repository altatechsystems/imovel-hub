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
  console.log('ðŸ”„ Iniciando migraÃ§Ã£o de imÃ³veis...');
  console.log('Tenant Zero ID:', tenantZeroId);
  console.log('');

  try {
    // Buscar todas as properties (coleÃ§Ã£o raiz - estrutura antiga)
    const propertiesSnapshot = await db.collection('properties').get();

    let migratedCount = 0;
    let errorCount = 0;

    console.log('Total de imÃ³veis encontrados:', propertiesSnapshot.size);
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

        console.log(\âœ… Migrado: \ - \\);
        migratedCount++;

        // 3. (Opcional) Deletar da coleÃ§Ã£o raiz apÃ³s confirmaÃ§Ã£o
        // await doc.ref.delete();

      } catch (error) {
        console.error(\âŒ Erro ao migrar \:\, error.message);
        errorCount++;
      }
    }

    console.log('');
    console.log('====================================');
    console.log('âœ… MigraÃ§Ã£o ConcluÃ­da!');
    console.log('====================================');
    console.log(\Total migrado: \\);
    console.log(\Erros: \\);
    console.log('');

  } catch (error) {
    console.error('âŒ Erro na migraÃ§Ã£o:', error);
    process.exit(1);
  }
}

// Executar migraÃ§Ã£o
const tenantZeroId = process.argv[2];

if (!tenantZeroId) {
  console.error('âŒ Erro: Tenant Zero ID nÃ£o fornecido');
  console.log('Uso: node scripts/migrate-properties.js <TENANT_ZERO_ID>');
  process.exit(1);
}

migrateProperties(tenantZeroId)
  .then(() => {
    console.log('âœ… Script finalizado com sucesso');
    process.exit(0);
  })
  .catch((error) => {
    console.error('âŒ Erro fatal:', error);
    process.exit(1);
  });
