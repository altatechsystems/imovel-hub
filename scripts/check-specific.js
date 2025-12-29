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

async function checkPropertyPhotos() {
  try {
    const tenantId = 'bd71c02b-5fa5-43df-8b46-a1df2206f1ef';
    const references = ['AP00330', 'TE00112', 'AP00331'];

    console.log('Checking properties and their listings\n');

    for (const ref of references) {
      console.log(`\nProperty ${ref}:`);

      const propertiesRef = db.collection('properties');
      const propSnapshot = await propertiesRef
        .where('tenant_id', '==', tenantId)
        .where('reference', '==', ref)
        .limit(1)
        .get();

      if (propSnapshot.empty) {
        console.log('  NOT FOUND');
        continue;
      }

      const propDoc = propSnapshot.docs[0];
      const property = propDoc.data();

      console.log('  ID:', propDoc.id);
      console.log('  canonical_listing_id:', property.canonical_listing_id || 'NOT SET');

      if (property.canonical_listing_id) {
        const listingDoc = await db.collection('listings').doc(property.canonical_listing_id).get();

        if (listingDoc.exists) {
          const listing = listingDoc.data();
          console.log('  Listing exists:', listingDoc.id);
          console.log('  Photos:', listing.photos ? listing.photos.length : 0);

          if (listing.photos && listing.photos.length > 0) {
            console.log('  First photo thumb:', listing.photos[0].thumb_url);
          }
        } else {
          console.log('  Listing NOT FOUND:', property.canonical_listing_id);
        }
      }
    }

  } catch (error) {
    console.error('Error:', error);
  }
}

checkPropertyPhotos().then(() => process.exit(0));
