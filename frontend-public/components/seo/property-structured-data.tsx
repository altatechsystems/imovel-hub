'use client';

import { Property } from '@/types/property';
import { formatCurrency, getPropertyTypeLabel, getTransactionTypeLabel } from '@/lib/utils';

interface PropertyStructuredDataProps {
  property: Property;
}

export function PropertyStructuredData({ property }: PropertyStructuredDataProps) {
  const price = property.sale_price || property.rental_price || 0;
  const priceType = property.transaction_type === 'rent' ? 'rental' : 'sale';

  const structuredData = {
    '@context': 'https://schema.org',
    '@type': 'RealEstateListing',
    name: property.title || `${getPropertyTypeLabel(property.property_type)} em ${property.neighborhood}`,
    description: property.description || `${getPropertyTypeLabel(property.property_type)} para ${getTransactionTypeLabel(property.transaction_type)} em ${property.city}`,
    url: typeof window !== 'undefined' ? window.location.href : '',
    image: property.cover_image_url ? [property.cover_image_url] : [],
    address: {
      '@type': 'PostalAddress',
      streetAddress: property.street && property.number ? `${property.street}, ${property.number}` : property.street,
      addressLocality: property.city,
      addressRegion: property.state,
      postalCode: property.zip_code,
      addressCountry: 'BR',
    },
    geo: property.latitude && property.longitude ? {
      '@type': 'GeoCoordinates',
      latitude: property.latitude,
      longitude: property.longitude,
    } : undefined,
    offers: {
      '@type': 'Offer',
      price: price,
      priceCurrency: 'BRL',
      availability: property.status === 'available'
        ? 'https://schema.org/InStock'
        : 'https://schema.org/OutOfStock',
      priceSpecification: {
        '@type': 'UnitPriceSpecification',
        price: price,
        priceCurrency: 'BRL',
        ...(priceType === 'rental' && {
          unitText: 'MONTH',
          billingDuration: 'P1M',
        }),
      },
    },
    numberOfRooms: property.bedrooms,
    numberOfBathroomsTotal: property.bathrooms,
    floorSize: property.total_area ? {
      '@type': 'QuantitativeValue',
      value: property.total_area,
      unitCode: 'MTK',
    } : undefined,
    amenityFeature: [
      ...(property.bedrooms ? [{
        '@type': 'LocationFeatureSpecification',
        name: 'Quartos',
        value: property.bedrooms,
      }] : []),
      ...(property.bathrooms ? [{
        '@type': 'LocationFeatureSpecification',
        name: 'Banheiros',
        value: property.bathrooms,
      }] : []),
      ...(property.parking_spaces ? [{
        '@type': 'LocationFeatureSpecification',
        name: 'Vagas de Estacionamento',
        value: property.parking_spaces,
      }] : []),
      ...(property.pool ? [{
        '@type': 'LocationFeatureSpecification',
        name: 'Piscina',
        value: true,
      }] : []),
      ...(property.gym ? [{
        '@type': 'LocationFeatureSpecification',
        name: 'Academia',
        value: true,
      }] : []),
      ...(property.elevator ? [{
        '@type': 'LocationFeatureSpecification',
        name: 'Elevador',
        value: true,
      }] : []),
      ...(property.furnished ? [{
        '@type': 'LocationFeatureSpecification',
        name: 'Mobiliado',
        value: true,
      }] : []),
    ],
  };

  // Remove undefined fields
  const cleanedData = JSON.parse(JSON.stringify(structuredData));

  return (
    <script
      type="application/ld+json"
      dangerouslySetInnerHTML={{ __html: JSON.stringify(cleanedData) }}
    />
  );
}
