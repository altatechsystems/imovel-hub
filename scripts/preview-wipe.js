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

const TENANT_ID = 'bd71c02b-5fa5-43df-8b46-a1df2206f1ef';

async function previewWipe() {
  console.log('🔍 DATABASE WIPE PREVIEW\n');
  console.log(`📍 Tenant ID: ${TENANT_ID}\n`);
  console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n');

  try {
    // 1. Count properties
    console.log('📦 Properties:');
    const propertiesSnapshot = await db.collection('properties')
      .where('tenant_id', '==', TENANT_ID)
      .count()
      .get();
    const propertiesCount = propertiesSnapshot.data().count;
    console.log(`   Total: ${propertiesCount} documents`);

    // Get sample
    const propertiesSample = await db.collection('properties')
      .where('tenant_id', '==', TENANT_ID)
      .limit(3)
      .get();

    console.log('   Sample:');
    propertiesSample.docs.forEach((doc, i) => {
      const data = doc.data();
      console.log(`     ${i + 1}. ${data.reference || data.title || doc.id}`);
    });
    console.log('');

    // 2. Count canonical_listings
    console.log('🏷️  Canonical Listings:');
    const listingsSnapshot = await db.collection('canonical_listings')
      .where('tenant_id', '==', TENANT_ID)
      .count()
      .get();
    const listingsCount = listingsSnapshot.data().count;
    console.log(`   Total: ${listingsCount} documents`);

    // Get sample
    const listingsSample = await db.collection('canonical_listings')
      .where('tenant_id', '==', TENANT_ID)
      .limit(3)
      .get();

    console.log('   Sample:');
    listingsSample.docs.forEach((doc, i) => {
      const data = doc.data();
      console.log(`     ${i + 1}. ${data.title || doc.id}`);
    });
    console.log('');

    // 3. Count listing_references
    console.log('🔗 Listing References:');
    const refsSnapshot = await db.collection('listing_references')
      .where('tenant_id', '==', TENANT_ID)
      .count()
      .get();
    const refsCount = refsSnapshot.data().count;
    console.log(`   Total: ${refsCount} documents`);

    // Get sample
    const refsSample = await db.collection('listing_references')
      .where('tenant_id', '==', TENANT_ID)
      .limit(3)
      .get();

    console.log('   Sample:');
    refsSample.docs.forEach((doc, i) => {
      const data = doc.data();
      console.log(`     ${i + 1}. ${data.source || 'unknown'} - ${data.external_id || doc.id}`);
    });
    console.log('');

    // 4. Summary
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('📊 SUMMARY\n');
    console.log('Documents to be deleted:');
    console.log(`   🗑️  ${propertiesCount.toLocaleString()} properties`);
    console.log(`   🗑️  ${listingsCount.toLocaleString()} canonical listings`);
    console.log(`   🗑️  ${refsCount.toLocaleString()} listing references`);
    console.log(`   ═══════════════════════════════════`);
    console.log(`   🗑️  ${(propertiesCount + listingsCount + refsCount).toLocaleString()} TOTAL documents\n`);

    // 5. Estimate time
    const totalDocs = propertiesCount + listingsCount + refsCount;
    const estimatedMinutes = Math.ceil(totalDocs / 500 / 2); // ~2 batches per minute
    console.log(`⏱️  Estimated time: ~${estimatedMinutes} minute(s)\n`);

    // 6. What will be preserved
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('✅ PRESERVED DATA\n');
    console.log('These will NOT be deleted:');
    console.log('   ✅ Tenants configuration');
    console.log('   ✅ Users and brokers');
    console.log('   ✅ User invitations');
    console.log('   ✅ Media files in Storage\n');

    // 7. Next steps
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('📝 NEXT STEPS\n');
    console.log('To wipe the database, run:');
    console.log('   node wipe-database.js\n');
    console.log('⚠️  Make sure you have a backup if needed!');
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n');

  } catch (error) {
    console.error('❌ Error:', error);
    process.exit(1);
  }
}

// Run preview
previewWipe()
  .then(() => {
    process.exit(0);
  })
  .catch((error) => {
    console.error('💥 Fatal error:', error);
    process.exit(1);
  });
