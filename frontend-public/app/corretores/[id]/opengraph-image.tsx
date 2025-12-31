import { ImageResponse } from 'next/og';

export const runtime = 'edge';
export const alt = 'Perfil do Corretor - Imobiliária';
export const size = {
  width: 1200,
  height: 630,
};
export const contentType = 'image/png';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
const TENANT_ID = process.env.NEXT_PUBLIC_TENANT_ID || 'bd71c02b-5fa5-43df-8b46-a1df2206f1ef';

async function getBrokerById(brokerId: string) {
  try {
    const response = await fetch(`${API_URL}/${TENANT_ID}/brokers/${brokerId}/public`, {
      next: { revalidate: 3600 },
    });

    if (!response.ok) {
      return null;
    }

    const data = await response.json();
    return data.data;
  } catch (error) {
    console.error('Error fetching broker for OG image:', error);
    return null;
  }
}

export default async function Image({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const broker = await getBrokerById(id);

  if (!broker) {
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
          <div style={{ fontSize: 60, fontWeight: 'bold' }}>Corretor não encontrado</div>
        </div>
      ),
      { ...size }
    );
  }

  const firstName = broker.name.split(' ')[0];
  const specialtiesList = broker.specialties
    ? broker.specialties.split(',').slice(0, 3).map((s: string) => s.trim())
    : [];

  const stats = [
    broker.rating ? `${broker.rating.toFixed(1)} ⭐` : null,
    broker.experience ? `${broker.experience} anos` : null,
    broker.review_count ? `${broker.review_count} avaliações` : null,
  ].filter(Boolean);

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
        {/* Header Background */}
        <div
          style={{
            width: '100%',
            height: '250px',
            background: 'linear-gradient(135deg, #2563eb 0%, #1e40af 100%)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            position: 'relative',
          }}
        >
          {broker.photo_url ? (
            <div
              style={{
                position: 'absolute',
                bottom: '-60px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img
                src={broker.photo_url}
                alt={broker.name}
                style={{
                  width: '200px',
                  height: '200px',
                  borderRadius: '100px',
                  objectFit: 'cover',
                  border: '8px solid white',
                  boxShadow: '0 10px 40px rgba(0,0,0,0.15)',
                }}
              />
            </div>
          ) : (
            <div
              style={{
                position: 'absolute',
                bottom: '-60px',
                width: '200px',
                height: '200px',
                borderRadius: '100px',
                background: '#e5e7eb',
                border: '8px solid white',
                boxShadow: '0 10px 40px rgba(0,0,0,0.15)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontSize: '80px',
                color: '#9ca3af',
              }}
            >
              👤
            </div>
          )}
        </div>

        {/* Broker Info */}
        <div
          style={{
            flex: 1,
            display: 'flex',
            flexDirection: 'column',
            padding: '80px 60px 40px',
            background: '#ffffff',
          }}
        >
          {/* Name and Title */}
          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              marginBottom: '30px',
            }}
          >
            <div
              style={{
                fontSize: '48px',
                fontWeight: 'bold',
                color: '#111827',
                marginBottom: '12px',
                textAlign: 'center',
              }}
            >
              {broker.name}
            </div>

            <div
              style={{
                fontSize: '24px',
                color: '#6b7280',
                marginBottom: '8px',
                display: 'flex',
                alignItems: 'center',
              }}
            >
              🏅 CRECI {broker.creci}
            </div>

            {broker.company && (
              <div
                style={{
                  fontSize: '20px',
                  color: '#9ca3af',
                  display: 'flex',
                }}
              >
                {broker.company}
              </div>
            )}
          </div>

          {/* Stats */}
          {stats.length > 0 && (
            <div
              style={{
                display: 'flex',
                gap: '30px',
                justifyContent: 'center',
                marginBottom: '25px',
              }}
            >
              {stats.map((stat, index) => (
                <div
                  key={index}
                  style={{
                    display: 'flex',
                    padding: '12px 24px',
                    background: '#f3f4f6',
                    borderRadius: '8px',
                    fontSize: '20px',
                    fontWeight: '600',
                    color: '#374151',
                  }}
                >
                  {stat}
                </div>
              ))}
            </div>
          )}

          {/* Specialties */}
          {specialtiesList.length > 0 && (
            <div
              style={{
                display: 'flex',
                gap: '12px',
                justifyContent: 'center',
                flexWrap: 'wrap',
              }}
            >
              {specialtiesList.map((specialty, index) => (
                <div
                  key={index}
                  style={{
                    display: 'flex',
                    padding: '8px 20px',
                    background: '#dbeafe',
                    color: '#1e40af',
                    borderRadius: '20px',
                    fontSize: '18px',
                    fontWeight: '500',
                  }}
                >
                  {specialty}
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Footer */}
        <div
          style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            padding: '25px 60px',
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
            Corretor Profissional
          </div>
        </div>
      </div>
    ),
    { ...size }
  );
}
