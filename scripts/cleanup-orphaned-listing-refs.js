const { initializeApp, cert } = require('firebase-admin/app');
const { getFirestore } = require('firebase-admin/firestore');
const path = require('path');

const serviceAccount = require(path.join(__dirname, '..', 'backend', 'config', 'firebase-adminsdk.json'));

initializeApp({
  credential: cert(serviceAccount),
  databaseURL: 'https://imob-dev.firebaseio.com'
});

const db = getFirestore();
db.settings({ databaseId: 'imob-dev' });

async function cleanupOrphanedListingRefs() {
  try {
    const tenantId = 'bd71c02b-5fa5-43df-8b46-a1df2206f1ef';

    console.log('Finding properties with orphaned listing references...\n');

    const propertiesRef = db.collection('properties');
    const snapshot = await propertiesRef.where('tenant_id', '==', tenantId).get();

    console.log(`Checking ${snapshot.size} properties...\n`);

    let orphanedCount = 0;
    const updates = [];

    for (const propertyDoc of snapshot.docs) {
      const property = propertyDoc.data();

      if (property.canonical_listing_id) {
        // Check if listing exists
        const listingDoc = await db.collection('listings').doc(property.canonical_listing_id).get();

        if (!listingDoc.exists) {
          console.log(`Property ${property.reference} (${propertyDoc.id})`);
          console.log(`  has orphaned listing ID: ${property.canonical_listing_id}`);
          console.log(`  Cleaning...`);

          updates.push({
            id: propertyDoc.id,
            ref: property.reference,
            oldListingId: property.canonical_listing_id,
          });

          orphanedCount++;
        }
      }
    }

    console.log(`\nFound ${orphanedCount} properties with orphaned listing references`);

    if (updates.length > 0) {
      console.log('\nCleaning up...\n');

      const batch = db.batch();

      for (const update of updates) {
        const propertyRef = db.collection('properties').doc(update.id);
        batch.update(propertyRef, {
          canonical_listing_id: null,
        });
        console.log(`Cleared orphaned listing from ${update.ref}`);
      }

      await batch.commit();

      console.log(`\nSuccessfully cleaned ${updates.length} properties!`);
    } else {
      console.log('\nNo orphaned listing references found!');
    }

  } catch (error) {
    console.error('Error:', error);
  }
}

cleanupOrphanedListingRefs().then(() => process.exit(0));
