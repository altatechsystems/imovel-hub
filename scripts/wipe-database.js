const { initializeApp, cert } = require('firebase-admin/app');
const { getFirestore } = require('firebase-admin/firestore');
const path = require('path');

// Initialize Firebase Admin
const serviceAccount = require(path.join(__dirname, '..', 'backend', 'config', 'firebase-adminsdk.json'));

initializeApp({
  credential: cert(serviceAccount),
  databaseURL: 'https://imob-dev.firebaseio.com'
});

const db = getFirestore();
db.settings({ databaseId: 'imob-dev' });

/**
 * SCRIPT DE LIMPEZA COMPLETA DA BASE DE DADOS
 *
 * ATENÇÃO: Este script irá DELETAR TODOS os dados das seguintes coleções:
 * - properties (imóveis)
 * - canonical_listings (anúncios canônicos)
 * - listing_references (referências de anúncios)
 * - Media files no Storage (fotos dos imóveis)
 *
 * USUÁRIOS, TENANTS E CONVITES SERÃO PRESERVADOS!
 *
 * Use apenas em ambiente de DESENVOLVIMENTO!
 */

const TENANT_ID = 'bd71c02b-5fa5-43df-8b46-a1df2206f1ef'; // Seu tenant de DEV
const BATCH_SIZE = 500; // Firestore limit

async function deleteCollection(collectionPath, batchSize = BATCH_SIZE) {
  const collectionRef = db.collection(collectionPath);
  const query = collectionRef.where('tenant_id', '==', TENANT_ID).limit(batchSize);

  return new Promise((resolve, reject) => {
    deleteQueryBatch(query, resolve, reject);
  });
}

async function deleteQueryBatch(query, resolve, reject) {
  try {
    const snapshot = await query.get();

    if (snapshot.size === 0) {
      resolve();
      return;
    }

    console.log(`  Deleting batch of ${snapshot.size} documents...`);

    const batch = db.batch();
    snapshot.docs.forEach((doc) => {
      batch.delete(doc.ref);
    });

    await batch.commit();

    // Recursively delete next batch
    process.nextTick(() => {
      deleteQueryBatch(query, resolve, reject);
    });
  } catch (error) {
    reject(error);
  }
}

async function wipeDatabase() {
  console.log('🚨 WIPE DATABASE SCRIPT 🚨\n');
  console.log('⚠️  WARNING: This will DELETE all property data!');
  console.log(`📍 Tenant ID: ${TENANT_ID}\n`);
  console.log('Collections to be wiped:');
  console.log('  - properties');
  console.log('  - canonical_listings');
  console.log('  - listing_references\n');
  console.log('Collections PRESERVED:');
  console.log('  ✅ tenants');
  console.log('  ✅ tenants/{tenant}/users');
  console.log('  ✅ tenants/{tenant}/brokers');
  console.log('  ✅ tenants/{tenant}/user_invitations\n');

  // Wait 5 seconds to allow cancellation
  console.log('⏳ Starting in 5 seconds... Press Ctrl+C to cancel!\n');
  await new Promise(resolve => setTimeout(resolve, 5000));

  try {
    // 1. Delete properties
    console.log('🗑️  Deleting properties...');
    const propertiesSnapshot = await db.collection('properties')
      .where('tenant_id', '==', TENANT_ID)
      .count()
      .get();
    const propertiesCount = propertiesSnapshot.data().count;
    console.log(`   Found ${propertiesCount} properties`);

    await deleteCollection('properties');
    console.log('   ✅ Properties deleted\n');

    // 2. Delete canonical_listings
    console.log('🗑️  Deleting canonical_listings...');
    const listingsSnapshot = await db.collection('canonical_listings')
      .where('tenant_id', '==', TENANT_ID)
      .count()
      .get();
    const listingsCount = listingsSnapshot.data().count;
    console.log(`   Found ${listingsCount} canonical listings`);

    await deleteCollection('canonical_listings');
    console.log('   ✅ Canonical listings deleted\n');

    // 3. Delete listing_references
    console.log('🗑️  Deleting listing_references...');
    const refsSnapshot = await db.collection('listing_references')
      .where('tenant_id', '==', TENANT_ID)
      .count()
      .get();
    const refsCount = refsSnapshot.data().count;
    console.log(`   Found ${refsCount} listing references`);

    await deleteCollection('listing_references');
    console.log('   ✅ Listing references deleted\n');

    // 4. Summary
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('✅ DATABASE WIPED SUCCESSFULLY!\n');
    console.log('📊 Summary:');
    console.log(`   - ${propertiesCount} properties deleted`);
    console.log(`   - ${listingsCount} canonical listings deleted`);
    console.log(`   - ${refsCount} listing references deleted`);
    console.log(`   - TOTAL: ${propertiesCount + listingsCount + refsCount} documents deleted\n`);
    console.log('🎯 Ready for fresh import!');
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n');

    // 5. Note about Storage
    console.log('📸 NOTE: Media files in Storage were NOT deleted.');
    console.log('   If you want to delete images, go to Firebase Console:');
    console.log('   Storage > ecosistema-imob-dev.firebasestorage.app > properties/');
    console.log('   And manually delete the folder.\n');

  } catch (error) {
    console.error('❌ Error wiping database:', error);
    process.exit(1);
  }
}

// Run the script
wipeDatabase()
  .then(() => {
    console.log('🏁 Script completed successfully!');
    process.exit(0);
  })
  .catch((error) => {
    console.error('💥 Fatal error:', error);
    process.exit(1);
  });
