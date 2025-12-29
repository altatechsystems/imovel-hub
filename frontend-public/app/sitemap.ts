import { MetadataRoute } from 'next';

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const baseUrl = process.env.NEXT_PUBLIC_SITE_URL || 'http://localhost:3001';
  const apiUrl = process.env.NEXT_PUBLIC_API_URL;
  const tenantId = process.env.NEXT_PUBLIC_TENANT_ID;

  const routes: MetadataRoute.Sitemap = [
    {
      url: baseUrl,
      lastModified: new Date(),
      changeFrequency: 'daily',
      priority: 1,
    },
    {
      url: `${baseUrl}/imoveis`,
      lastModified: new Date(),
      changeFrequency: 'daily',
      priority: 0.9,
    },
    {
      url: `${baseUrl}/sobre`,
      lastModified: new Date(),
      changeFrequency: 'monthly',
      priority: 0.5,
    },
    {
      url: `${baseUrl}/contato`,
      lastModified: new Date(),
      changeFrequency: 'monthly',
      priority: 0.5,
    },
  ];

  // Fetch properties for dynamic sitemap
  try {
    if (apiUrl && tenantId) {
      const response = await fetch(
        `${apiUrl}/${tenantId}/properties?limit=1000&visibility=public&status=available`,
        {
          next: { revalidate: 3600 }, // Revalidate every hour
        }
      );

      if (response.ok) {
        const data = await response.json();
        const properties = data.data || [];

        const propertySitemapEntries: MetadataRoute.Sitemap = properties.map((property: any) => ({
          url: `${baseUrl}/imoveis/${property.slug || property.id}`,
          lastModified: property.updated_at ? new Date(property.updated_at) : new Date(),
          changeFrequency: 'weekly' as const,
          priority: property.featured ? 0.8 : 0.7,
        }));

        routes.push(...propertySitemapEntries);
      }
    }
  } catch (error) {
    console.error('Error generating sitemap:', error);
  }

  return routes;
}
