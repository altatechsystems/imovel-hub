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

async function checkDuplicates() {
  try {
    const propertiesRef = db.collection('properties');
    const snapshot = await propertiesRef.where('reference', '==', 'AP00330').get();

    console.log(`\n📊 Found ${snapshot.size} properties with reference AP00330:\n`);

    snapshot.forEach(doc => {
      const data = doc.data();
      console.log(`Property ID: ${doc.id}`);
      console.log(`  reference: ${data.reference}`);
      console.log(`  external_id: ${data.external_id}`);
      console.log(`  external_source: ${data.external_source}`);
      console.log(`  tenant_id: ${data.tenant_id}`);
      console.log(`  canonical_listing_id: ${data.canonical_listing_id || 'NOT SET'}`);
      console.log(`  created_at: ${data.created_at?.toDate?.() || data.created_at}`);
      console.log('');
    });

    // Check listings
    console.log('\n📋 Checking ALL listings:\n');
    const listingsRef = db.collection('listings');
    const listingsSnapshot = await listingsRef.limit(5).get();

    console.log(`Found ${listingsSnapshot.size} listings (showing first 5):\n`);
    listingsSnapshot.forEach(doc => {
      const data = doc.data();
      console.log(`Listing ID: ${doc.id}`);
      console.log(`  property_id: ${data.property_id}`);
      console.log(`  tenant_id: ${data.tenant_id}`);
      console.log(`  photos: ${data.photos?.length || 0}`);
      console.log('');
    });

  } catch (error) {
    console.error('❌ Error:', error);
  }
}

checkDuplicates().then(() => process.exit(0));
