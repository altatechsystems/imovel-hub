import { ImageResponse } from 'next/og';

export const runtime = 'edge';
export const alt = 'Imóvel - Imobiliária';
export const size = {
  width: 1200,
  height: 630,
};
export const contentType = 'image/png';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
const TENANT_ID = process.env.NEXT_PUBLIC_TENANT_ID || 'bd71c02b-5fa5-43df-8b46-a1df2206f1ef';

async function getPropertyBySlug(slug: string) {
  try {
    const response = await fetch(`${API_URL}/${TENANT_ID}/properties/slug/${slug}`, {
      next: { revalidate: 3600 },
    });

    if (!response.ok) {
      return null;
    }

    const data = await response.json();
    return data.data;
  } catch (error) {
    console.error('Error fetching property for OG image:', error);
    return null;
  }
}

function getPropertyTypeLabel(type: string): string {
  const labels: Record<string, string> = {
    apartment: 'Apartamento',
    house: 'Casa',
    condo: 'Condomínio',
    land: 'Terreno',
    commercial: 'Comercial',
    farm: 'Fazenda',
    warehouse: 'Galpão',
    office: 'Escritório',
  };
  return labels[type] || 'Imóvel';
}

export default async function Image({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;
  const property = await getPropertyBySlug(slug);

  if (!property) {
    // Fallback image
    return new ImageResponse(
      (
        <div
          style={{
            background: 'linear-gradient(135deg, #2563eb 0%, #1e40af 100%)',
            width: '100%',
            height: '100%',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontFamily: 'system-ui',
            color: 'white',
          }}
        >
          <div style={{ fontSize: 60, fontWeight: 'bold' }}>Imóvel não encontrado</div>
        </div>
      ),
      { ...size }
    );
  }

  const propertyType = getPropertyTypeLabel(property.property_type);
  const location = property.neighborhood
    ? `${property.neighborhood}, ${property.city}`
    : property.city;

  const price = property.sale_price || property.rental_price || property.price_amount;
  const formattedPrice = price
    ? new Intl.NumberFormat('pt-BR', {
        style: 'currency',
        currency: 'BRL',
        minimumFractionDigits: 0,
      }).format(price)
    : '';

  const features = [
    property.bedrooms ? `${property.bedrooms} quartos` : null,
    property.bathrooms ? `${property.bathrooms} banheiros` : null,
    property.parking_spaces ? `${property.parking_spaces} vagas` : null,
    property.area_sqm ? `${property.area_sqm}m²` : null,
  ].filter(Boolean);

  // Try to use property image if available
  const hasImage = property.cover_image_url || (property.images && property.images.length > 0);
  const imageUrl = property.cover_image_url || property.images?.[0]?.large_url;

  return new ImageResponse(
    (
      <div
        style={{
          width: '100%',
          height: '100%',
          display: 'flex',
          flexDirection: 'column',
          background: '#ffffff',
        }}
      >
        {/* Property Image or Gradient Background */}
        {hasImage && imageUrl ? (
          <div
            style={{
              width: '100%',
              height: '400px',
              display: 'flex',
              position: 'relative',
            }}
          >
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={imageUrl}
              alt="Property"
              style={{
                width: '100%',
                height: '100%',
                objectFit: 'cover',
              }}
            />
            {/* Gradient overlay for better text visibility */}
            <div
              style={{
                position: 'absolute',
                bottom: 0,
                left: 0,
                right: 0,
                height: '150px',
                background: 'linear-gradient(to top, rgba(0,0,0,0.7), transparent)',
              }}
            />
          </div>
        ) : (
          <div
            style={{
              width: '100%',
              height: '400px',
              background: 'linear-gradient(135deg, #2563eb 0%, #1e40af 100%)',
            }}
          />
        )}

        {/* Property Info */}
        <div
          style={{
            flex: 1,
            display: 'flex',
            flexDirection: 'column',
            padding: '40px',
            background: '#ffffff',
          }}
        >
          {/* Type Badge */}
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              marginBottom: '20px',
            }}
          >
            <div
              style={{
                background: '#2563eb',
                color: 'white',
                padding: '8px 20px',
                borderRadius: '6px',
                fontSize: '20px',
                fontWeight: '600',
              }}
            >
              {propertyType}
            </div>
          </div>

          {/* Title/Location */}
          <div
            style={{
              fontSize: '36px',
              fontWeight: 'bold',
              color: '#111827',
              marginBottom: '12px',
              display: 'flex',
              lineHeight: 1.2,
            }}
          >
            {property.title || `${propertyType} em ${location}`}
          </div>

          {/* Features */}
          {features.length > 0 && (
            <div
              style={{
                display: 'flex',
                gap: '20px',
                fontSize: '24px',
                color: '#6b7280',
                marginBottom: '20px',
              }}
            >
              {features.join(' • ')}
            </div>
          )}

          {/* Price */}
          {formattedPrice && (
            <div
              style={{
                fontSize: '48px',
                fontWeight: 'bold',
                color: '#2563eb',
                display: 'flex',
              }}
            >
              {formattedPrice}
            </div>
          )}
        </div>

        {/* Footer */}
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            padding: '30px 40px',
            background: '#f9fafb',
            borderTop: '2px solid #e5e7eb',
          }}
        >
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '12px',
            }}
          >
            <div
              style={{
                fontSize: '28px',
                fontWeight: 'bold',
                color: '#2563eb',
              }}
            >
              🏠 Imobiliária
            </div>
          </div>
          <div
            style={{
              fontSize: '20px',
              color: '#6b7280',
            }}
          >
            {location}
          </div>
        </div>
      </div>
    ),
    { ...size }
  );
}
