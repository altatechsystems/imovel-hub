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

async function cleanupOrphanedProperties() {
  try {
    const validTenantId = 'bd71c02b-5fa5-43df-8b46-a1df2206f1ef';
    const orphanedTenantId = 'DMsXDI6CcIsIE5LPICiW'; // Deleted tenant

    console.log('🔍 Searching for orphaned properties...\n');

    // Find all properties from the deleted tenant
    const propertiesRef = db.collection('properties');
    const orphanedSnapshot = await propertiesRef.where('tenant_id', '==', orphanedTenantId).get();

    console.log(`Found ${orphanedSnapshot.size} properties from deleted tenant ${orphanedTenantId}\n`);

    if (orphanedSnapshot.size === 0) {
      console.log('✅ No orphaned properties to delete!');
      return;
    }

    // List properties to delete
    console.log('📋 Properties to be deleted:\n');
    orphanedSnapshot.forEach(doc => {
      const data = doc.data();
      console.log(`  - ${data.reference || 'NO-REF'} (${doc.id})`);
    });

    console.log(`\n⚠️  WARNING: About to delete ${orphanedSnapshot.size} properties`);
    console.log('⚠️  This action is IRREVERSIBLE!\n');

    // Delete orphaned properties
    console.log('🗑️  Deleting orphaned properties...\n');

    let deleted = 0;
    const batch = db.batch();

    orphanedSnapshot.forEach(doc => {
      batch.delete(doc.ref);
      deleted++;
    });

    await batch.commit();

    console.log(`✅ Successfully deleted ${deleted} orphaned properties!`);

    // Check for remaining duplicates in valid tenant
    console.log('\n\n🔍 Checking for duplicates in valid tenant...\n');

    const validSnapshot = await propertiesRef.where('tenant_id', '==', validTenantId).get();
    console.log(`Found ${validSnapshot.size} properties in valid tenant\n`);

    // Group by external_id
    const grouped = {};
    validSnapshot.forEach(doc => {
      const data = doc.data();
      const key = `${data.external_source || 'no-source'}_${data.external_id || 'no-id'}`;

      if (!grouped[key]) {
        grouped[key] = [];
      }

      grouped[key].push({
        id: doc.id,
        reference: data.reference,
        external_id: data.external_id,
        canonical_listing_id: data.canonical_listing_id,
        created_at: data.created_at?.toDate?.() || data.created_at,
      });
    });

    // Find duplicates
    let totalDuplicates = 0;

    for (const [key, properties] of Object.entries(grouped)) {
      if (properties.length > 1) {
        totalDuplicates++;
        console.log(`\n⚠️  Duplicate: ${properties[0].reference || 'NO-REF'} (${key})`);
        console.log(`   Copies: ${properties.length}`);
        properties.forEach(prop => {
          console.log(`   - ${prop.id} (created: ${prop.created_at}, has_listing: ${!!prop.canonical_listing_id})`);
        });
      }
    }

    if (totalDuplicates === 0) {
      console.log('\n✅ No duplicates found in valid tenant!');
    } else {
      console.log(`\n⚠️  Found ${totalDuplicates} duplicate groups in valid tenant`);
      console.log('   These may need manual cleanup');
    }

  } catch (error) {
    console.error('❌ Error:', error);
  }
}

cleanupOrphanedProperties().then(() => process.exit(0));
