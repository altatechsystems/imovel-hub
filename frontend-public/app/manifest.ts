import { MetadataRoute } from 'next';

export default function manifest(): MetadataRoute.Manifest {
  return {
    name: 'Imobiliária - Encontre seu Imóvel Ideal',
    short_name: 'Imobiliária',
    description: 'Encontre apartamentos, casas, terrenos e imóveis comerciais para venda e aluguel em todo o Brasil.',
    start_url: '/',
    display: 'standalone',
    background_color: '#ffffff',
    theme_color: '#2563eb',
    lang: 'pt-BR',
  };
}
