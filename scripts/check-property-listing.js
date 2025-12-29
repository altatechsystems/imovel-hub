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

async function checkPropertyListingMatch() {
  try {
    // Get property AP00330
    const propertiesRef = db.collection('properties');
    const snapshot = await propertiesRef.where('reference', '==', 'AP00330').limit(1).get();

    if (snapshot.empty) {
      console.log('❌ Property AP00330 not found');
      return;
    }

    const propertyDoc = snapshot.docs[0];
    const property = propertyDoc.data();

    console.log('\n📄 Property AP00330:');
    console.log('  ID:', propertyDoc.id);
    console.log('  canonical_listing_id:', property.canonical_listing_id || 'NOT SET');

    // Check if canonical_listing_id exists
    if (property.canonical_listing_id) {
      const listingDoc = await db.collection('listings').doc(property.canonical_listing_id).get();

      if (listingDoc.exists) {
        const listing = listingDoc.data();
        console.log('\n✅ Canonical Listing EXISTS:');
        console.log('  ID:', listingDoc.id);
        console.log('  property_id:', listing.property_id);
        console.log('  tenant_id:', listing.tenant_id);
        console.log('  photos count:', listing.photos?.length || 0);

        if (listing.photos && listing.photos.length > 0) {
          console.log('  First photo thumb_url:', listing.photos[0].thumb_url);
        }

        // Check if property_id matches
        if (listing.property_id === propertyDoc.id) {
          console.log('\n✅ MATCH: listing.property_id === property.id');
        } else {
          console.log('\n❌ MISMATCH:');
          console.log('   listing.property_id:', listing.property_id);
          console.log('   property.id:', propertyDoc.id);
        }
      } else {
        console.log('\n❌ Canonical Listing DOES NOT EXIST in Firestore');
        console.log('   Expected ID:', property.canonical_listing_id);
      }
    } else {
      console.log('\n⚠️  Property has no canonical_listing_id set');

      // Search for listings with this property_id
      const listingsRef = db.collection('listings');
      const listingsSnapshot = await listingsRef.where('property_id', '==', propertyDoc.id).get();

      console.log('\n🔍 Searching for listings with property_id =', propertyDoc.id);
      console.log('   Found:', listingsSnapshot.size, 'listings');

      if (!listingsSnapshot.empty) {
        listingsSnapshot.forEach(doc => {
          console.log('   - Listing ID:', doc.id);
        });
      }
    }

  } catch (error) {
    console.error('❌ Error:', error);
  }
}

checkPropertyListingMatch().then(() => process.exit(0));
