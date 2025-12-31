import { Metadata } from 'next';

const SITE_URL = process.env.NEXT_PUBLIC_SITE_URL || 'http://localhost:3001';

export const metadata: Metadata = {
  title: 'Imóveis para Venda e Aluguel - Encontre sua Propriedade Ideal',
  description: 'Navegue por milhares de imóveis à venda e para alugar. Apartamentos, casas, terrenos, imóveis comerciais em todo o Brasil. Filtros avançados por preço, localização, tipo e características.',
  keywords: [
    'imóveis para venda',
    'imóveis para aluguel',
    'apartamentos à venda',
    'casas para alugar',
    'terrenos à venda',
    'imóveis comerciais',
    'comprar apartamento',
    'alugar casa',
    'busca de imóveis',
    'propriedades',
    'corretora de imóveis',
    'imóveis Brasil',
    'listings imobiliários',
  ],
  openGraph: {
    title: 'Imóveis para Venda e Aluguel - Encontre sua Propriedade Ideal',
    description: 'Navegue por milhares de imóveis à venda e para alugar. Apartamentos, casas, terrenos e imóveis comerciais com filtros avançados.',
    url: `${SITE_URL}/imoveis`,
    type: 'website',
    locale: 'pt_BR',
    siteName: 'Imobiliária',
    images: [
      {
        url: `${SITE_URL}/og-image.jpg`,
        width: 1200,
        height: 630,
        alt: 'Imóveis para Venda e Aluguel',
      },
    ],
  },
  twitter: {
    card: 'summary_large_image',
    title: 'Imóveis para Venda e Aluguel - Encontre sua Propriedade Ideal',
    description: 'Navegue por milhares de imóveis à venda e para alugar. Apartamentos, casas, terrenos e imóveis comerciais.',
    images: [`${SITE_URL}/twitter-image.jpg`],
  },
  alternates: {
    canonical: `${SITE_URL}/imoveis`,
  },
  robots: {
    index: true,
    follow: true,
  },
};

export default function PropertiesLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <>{children}</>;
}
