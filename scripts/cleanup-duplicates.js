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

async function cleanupDuplicates() {
  try {
    const tenantId = 'bd71c02b-5fa5-43df-8b46-a1df2206f1ef';

    console.log('🔍 Finding all properties...\n');

    const propertiesRef = db.collection('properties');
    const snapshot = await propertiesRef.where('tenant_id', '==', tenantId).get();

    console.log(`Found ${snapshot.size} total properties\n`);

    // Group by external_id
    const grouped = {};
    snapshot.forEach(doc => {
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
    let totalToDelete = 0;

    console.log('📊 Duplicate Properties:\n');

    for (const [key, properties] of Object.entries(grouped)) {
      if (properties.length > 1) {
        totalDuplicates++;

        // Sort by created_at (most recent first)
        properties.sort((a, b) => {
          const dateA = a.created_at instanceof Date ? a.created_at : new Date(a.created_at);
          const dateB = b.created_at instanceof Date ? b.created_at : new Date(b.created_at);
          return dateB - dateA;
        });

        const newest = properties[0];
        const toDelete = properties.slice(1);

        console.log(`\n${properties[0].reference || 'NO-REF'} (${key}):`);
        console.log(`  Total copies: ${properties.length}`);
        console.log(`  ✅ KEEP: ${newest.id} (created: ${newest.created_at}, has_listing: ${!!newest.canonical_listing_id})`);

        for (const prop of toDelete) {
          console.log(`  ❌ DELETE: ${prop.id} (created: ${prop.created_at}, has_listing: ${!!prop.canonical_listing_id})`);
          totalToDelete++;
        }
      }
    }

    console.log(`\n\n📈 Summary:`);
    console.log(`  Total properties: ${snapshot.size}`);
    console.log(`  Duplicate groups: ${totalDuplicates}`);
    console.log(`  Properties to delete: ${totalToDelete}`);
    console.log(`  Properties to keep: ${snapshot.size - totalToDelete}`);

    console.log(`\n⚠️  This script is in DRY-RUN mode (no actual deletions)`);
    console.log(`To actually delete duplicates, you would need to uncomment the deletion code.`);

    // UNCOMMENT BELOW TO ACTUALLY DELETE
    /*
    console.log('\n🗑️  Starting deletion...\n');

    let deleted = 0;
    for (const [key, properties] of Object.entries(grouped)) {
      if (properties.length > 1) {
        properties.sort((a, b) => {
          const dateA = a.created_at instanceof Date ? a.created_at : new Date(a.created_at);
          const dateB = b.created_at instanceof Date ? b.created_at : new Date(b.created_at);
          return dateB - dateA;
        });

        const toDelete = properties.slice(1);

        for (const prop of toDelete) {
          await db.collection('properties').doc(prop.id).delete();
          console.log(`✅ Deleted ${prop.id} (${prop.reference})`);
          deleted++;
        }
      }
    }

    console.log(`\n✅ Deleted ${deleted} duplicate properties!`);
    */

  } catch (error) {
    console.error('❌ Error:', error);
  }
}

cleanupDuplicates().then(() => process.exit(0));
