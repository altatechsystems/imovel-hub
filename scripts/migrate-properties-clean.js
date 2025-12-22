// migrate-properties-clean.js
// Uso: node scripts/migrate-properties-clean.js <TENANT_ZERO_ID>

const admin = require('firebase-admin');

// Inicializar Firebase Admin SDK
const serviceAccount = require('../backend/config/firebase-adminsdk.json');

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount),
  databaseURL: 'https://ecosistema-imob-dev.firebaseio.com'
});

// Conectar ao banco "imob-dev"
const db = admin.firestore();
db.settings({ databaseId: 'imob-dev' });

async function migrateProperties(tenantZeroId) {
  console.log('Iniciando migracao de imoveis...');
  console.log('Tenant Zero ID:', tenantZeroId);
  console.log('Database:', 'imob-dev');
  console.log('');

  try {
    // Buscar todas as properties (colecao raiz - estrutura antiga)
    const propertiesSnapshot = await db.collection('properties').get();

    let migratedCount = 0;
    let errorCount = 0;

    console.log('Total de imoveis encontrados na colecao raiz:', propertiesSnapshot.size);
    console.log('');

    if (propertiesSnapshot.size === 0) {
      console.log('Nenhum imovel encontrado na colecao raiz.');
      console.log('Verificando imoveis ja migrados...');

      const tenantPropertiesSnapshot = await db
        .collection('tenants')
        .doc(tenantZeroId)
        .collection('properties')
        .get();

      console.log('Imoveis ja na estrutura multi-tenant:', tenantPropertiesSnapshot.size);
      return;
    }

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

        const address = propertyData.address?.street || 'Sem endereco';
        console.log(`OK - Migrado: ${propertyId} - ${address}`);
        migratedCount++;

        // 3. (Opcional) Deletar da colecao raiz apos confirmacao
        // await doc.ref.delete();

      } catch (error) {
        console.error(`ERRO ao migrar ${propertyId}:`, error.message);
        errorCount++;
      }
    }

    console.log('');
    console.log('====================================');
    console.log('Migracao Concluida!');
    console.log('====================================');
    console.log(`Total migrado: ${migratedCount}`);
    console.log(`Erros: ${errorCount}`);
    console.log('');

  } catch (error) {
    console.error('ERRO na migracao:', error);
    process.exit(1);
  }
}

// Executar migracao
const tenantZeroId = process.argv[2];

if (!tenantZeroId) {
  console.error('ERRO: Tenant Zero ID nao fornecido');
  console.log('Uso: node scripts/migrate-properties-clean.js <TENANT_ZERO_ID>');
  process.exit(1);
}

migrateProperties(tenantZeroId)
  .then(() => {
    console.log('Script finalizado com sucesso');
    process.exit(0);
  })
  .catch((error) => {
    console.error('Erro fatal:', error);
    process.exit(1);
  });
