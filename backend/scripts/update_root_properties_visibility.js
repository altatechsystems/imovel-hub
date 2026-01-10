const admin = require('firebase-admin');

// Initialize Firebase Admin
const serviceAccount = require('../config/firebase-adminsdk.json');

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount),
  databaseURL: `https://${serviceAccount.project_id}.firebaseio.com`
});

// Connect to named database "imob-dev"
const db = admin.firestore();
db.settings({ databaseId: 'imob-dev' });

async function updatePropertiesVisibility() {
  try {
    console.log('🔍 Buscando todos os imóveis na coleção root/properties...\n');

    // Get all properties from root collection
    const propertiesSnapshot = await db.collection('properties').get();

    console.log(`📊 Encontrados ${propertiesSnapshot.size} imóveis\n`);

    let updatedCount = 0;
    let alreadyPublicCount = 0;

    for (const propertyDoc of propertiesSnapshot.docs) {
      const propertyData = propertyDoc.data();
      const currentVisibility = propertyData.visibility || 'network';

      console.log(`📍 Imóvel: ${propertyDoc.id}`);
      console.log(`   Referência: ${propertyData.reference || 'N/A'}`);
      console.log(`   Tenant: ${propertyData.tenant_id}`);
      console.log(`   Visibilidade atual: ${currentVisibility}`);

      if (currentVisibility !== 'public') {
        console.log(`   ✏️  Atualizando para: public`);

        await propertyDoc.ref.update({
          visibility: 'public',
          updated_at: admin.firestore.FieldValue.serverTimestamp()
        });

        updatedCount++;
      } else {
        console.log(`   ✓  Já está público`);
        alreadyPublicCount++;
      }
      console.log('');
    }

    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('✅ Atualização concluída!');
    console.log(`📊 Total de imóveis: ${propertiesSnapshot.size}`);
    console.log(`📝 Imóveis atualizados: ${updatedCount}`);
    console.log(`✓  Já estavam públicos: ${alreadyPublicCount}`);
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');

  } catch (error) {
    console.error('❌ Erro ao atualizar imóveis:', error);
    process.exit(1);
  }

  process.exit(0);
}

// Run the update
updatePropertiesVisibility();
