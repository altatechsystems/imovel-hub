'use client';

import * as React from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { Property } from '@/types/property';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
  formatCurrency,
  formatArea,
  getPropertyTypeLabel,
  getTransactionTypeLabel,
} from '@/lib/utils';
import { Bed, Bath, Car, MapPin, Maximize2 } from 'lucide-react';

// Placeholder SVG for properties without images
const PLACEHOLDER_IMAGE = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAwIiBoZWlnaHQ9IjMwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iNDAwIiBoZWlnaHQ9IjMwMCIgZmlsbD0iI2YzZjRmNiIvPjx0ZXh0IHg9IjUwJSIgeT0iNTAlIiBmb250LWZhbWlseT0iQXJpYWwiIGZvbnQtc2l6ZT0iMTgiIGZpbGw9IiM5Y2EzYWYiIHRleHQtYW5jaG9yPSJtaWRkbGUiIGR5PSIuM2VtIj5TZW0gSW1hZ2VtPC90ZXh0Pjwvc3ZnPg==';

export interface PropertyCardProps {
  property: Property;
  variant?: 'grid' | 'list';
}

export const PropertyCard = React.memo(function PropertyCard({ property, variant = 'grid' }: PropertyCardProps) {
  const price = property.sale_price || property.rental_price || property.price_amount;
  const priceLabel = property.transaction_type === 'rent' ? 'Aluguel' : property.transaction_type === 'sale' ? 'Venda' : 'Valor';

  const features = [
    { icon: Bed, value: property.bedrooms, label: 'quartos' },
    { icon: Bath, value: property.bathrooms, label: 'banheiros' },
    { icon: Car, value: property.parking_spaces, label: 'vagas' },
    { icon: Maximize2, value: property.area_sqm ? formatArea(property.area_sqm) : null, label: '' },
  ].filter(f => f.value);

  if (variant === 'list') {
    return (
      <Card variant="bordered" padding="none" className="hover:shadow-md transition-shadow">
        <Link href={`/imoveis/${property.slug || property.id}`} className="flex flex-col sm:flex-row">
          {/* Image */}
          <div className="relative w-full sm:w-64 md:w-80 h-56 sm:h-64 md:h-auto flex-shrink-0">
            <Image
              src={property.cover_image_url || PLACEHOLDER_IMAGE}
              alt={property.title || `${getPropertyTypeLabel(property.property_type)} em ${property.city}`}
              fill
              sizes="(max-width: 640px) 100vw, (max-width: 768px) 256px, 320px"
              className="object-cover rounded-t-lg sm:rounded-l-lg sm:rounded-tr-none"
              loading="lazy"
              quality={60}
              placeholder="blur"
              blurDataURL="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAwIiBoZWlnaHQ9IjMwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iNDAwIiBoZWlnaHQ9IjMwMCIgZmlsbD0iI2RkZCIvPjwvc3ZnPg=="
            />
            <div className="absolute top-2 left-2 sm:top-3 sm:left-3 flex gap-1.5 sm:gap-2">
              {property.featured && (
                <Badge variant="featured" size="sm">Destaque</Badge>
              )}
              <Badge variant="info" size="sm">
                {getTransactionTypeLabel(property.transaction_type || 'sale')}
              </Badge>
            </div>
          </div>

          {/* Content */}
          <div className="flex-1 p-3 sm:p-4">
            <div className="flex flex-col h-full">
              <div className="flex-1">
                <h3 className="text-base sm:text-lg md:text-xl font-semibold text-gray-900 mb-2 line-clamp-2">
                  {property.title || `${getPropertyTypeLabel(property.property_type)} em ${property.neighborhood}`}
                </h3>

                <div className="flex items-center text-gray-600 text-xs sm:text-sm mb-2 sm:mb-3">
                  <MapPin className="w-3 h-3 sm:w-4 sm:h-4 mr-1 flex-shrink-0" />
                  <span className="line-clamp-1">{property.neighborhood}, {property.city} - {property.state}</span>
                </div>

                {property.description && (
                  <p className="text-gray-600 text-xs sm:text-sm mb-3 sm:mb-4 line-clamp-2 hidden sm:block">
                    {property.description}
                  </p>
                )}

                <div className="flex flex-wrap gap-2 sm:gap-3 md:gap-4 mb-3 sm:mb-4">
                  {features.map((feature, index) => (
                    <div key={index} className="flex items-center text-gray-700 text-sm">
                      <feature.icon className="w-4 h-4 mr-1.5 text-gray-500" />
                      <span>{typeof feature.value === 'number' ? feature.value : feature.value} {feature.label}</span>
                    </div>
                  ))}
                </div>
              </div>

              <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between pt-3 sm:pt-4 border-t gap-3">
                <div>
                  <p className="text-xs sm:text-sm text-gray-600">{priceLabel}</p>
                  <p className="text-xl sm:text-2xl font-bold text-blue-600">
                    {formatCurrency(price)}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </Link>
      </Card>
    );
  }

  // Grid variant (default)
  return (
    <Card variant="bordered" padding="none" className="hover:shadow-md transition-shadow">
      <Link href={`/imoveis/${property.slug || property.id}`}>
        {/* Image */}
        <div className="relative w-full h-48 sm:h-56">
          <Image
            src={property.cover_image_url || PLACEHOLDER_IMAGE}
            alt={property.title || `${getPropertyTypeLabel(property.property_type)} em ${property.city}`}
            fill
            sizes="(max-width: 640px) 100vw, (max-width: 1024px) 50vw, 33vw"
            className="object-cover rounded-t-lg"
            loading="lazy"
            quality={60}
            placeholder="blur"
            blurDataURL="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAwIiBoZWlnaHQ9IjMwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iNDAwIiBoZWlnaHQ9IjMwMCIgZmlsbD0iI2RkZCIvPjwvc3ZnPg=="
          />
          <div className="absolute top-2 left-2 sm:top-3 sm:left-3 flex gap-1.5 sm:gap-2">
            {property.featured && (
              <Badge variant="featured" size="sm">Destaque</Badge>
            )}
            <Badge variant="info" size="sm">
              {getTransactionTypeLabel(property.transaction_type || 'sale')}
            </Badge>
          </div>
        </div>

        <CardContent className="p-3 sm:p-4">
          {/* Price */}
          <div className="mb-2 sm:mb-3">
            <p className="text-xs sm:text-sm text-gray-600">{priceLabel}</p>
            <p className="text-xl sm:text-2xl font-bold text-blue-600">
              {formatCurrency(price)}
            </p>
          </div>

          {/* Title */}
          <h3 className="text-base sm:text-lg font-semibold text-gray-900 mb-2 line-clamp-2">
            {property.title || `${getPropertyTypeLabel(property.property_type)} em ${property.neighborhood}`}
          </h3>

          {/* Location */}
          <div className="flex items-center text-gray-600 text-xs sm:text-sm mb-2 sm:mb-3">
            <MapPin className="w-3 h-3 sm:w-4 sm:h-4 mr-1 flex-shrink-0" />
            <span className="line-clamp-1">{property.neighborhood}, {property.city} - {property.state}</span>
          </div>

          {/* Features */}
          <div className="flex flex-wrap gap-2 sm:gap-3 mb-3 sm:mb-4">
            {features.map((feature, index) => (
              <div key={index} className="flex items-center text-gray-700 text-xs sm:text-sm">
                <feature.icon className="w-3 h-3 sm:w-4 sm:h-4 mr-1 text-gray-500 flex-shrink-0" />
                <span>{typeof feature.value === 'number' ? feature.value : feature.value} {feature.label}</span>
              </div>
            ))}
          </div>
        </CardContent>
      </Link>
    </Card>
  );
});
