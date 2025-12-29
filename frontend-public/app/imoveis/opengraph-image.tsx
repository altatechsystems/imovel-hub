import { ImageResponse } from 'next/og';

export const runtime = 'edge';

export const alt = 'Imóveis Disponíveis - Imobiliária';
export const size = {
  width: 1200,
  height: 630,
};
export const contentType = 'image/png';

export default async function Image() {
  return new ImageResponse(
    (
      <div
        style={{
          background: 'linear-gradient(135deg, #2563eb 0%, #1e40af 100%)',
          width: '100%',
          height: '100%',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          fontFamily: 'system-ui',
          color: 'white',
        }}
      >
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            textAlign: 'center',
            padding: '60px',
          }}
        >
          <h1
            style={{
              fontSize: 72,
              fontWeight: 'bold',
              marginBottom: 24,
              lineHeight: 1.2,
            }}
          >
            Imóveis Disponíveis
          </h1>
          <p
            style={{
              fontSize: 36,
              opacity: 0.9,
              marginBottom: 48,
            }}
          >
            Encontre apartamentos, casas e terrenos
          </p>
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: 16,
              fontSize: 28,
            }}
          >
            <div
              style={{
                background: 'rgba(255, 255, 255, 0.2)',
                padding: '12px 24px',
                borderRadius: 12,
              }}
            >
              🏠 Imobiliária
            </div>
          </div>
        </div>
      </div>
    ),
    {
      ...size,
    }
  );
}
