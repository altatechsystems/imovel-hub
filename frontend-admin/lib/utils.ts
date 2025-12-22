import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

// Format currency to BRL
export function formatCurrency(value: number | undefined): string {
  if (!value) return 'Preço sob consulta';

  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(value);
}

// Format area (square meters)
export function formatArea(area: number | undefined): string {
  if (!area) return '-';
  return `${area.toFixed(0)} m²`;
}

// Format phone number
export function formatPhone(phone: string | undefined): string {
  if (!phone) return '';

  // Remove non-digits
  const digits = phone.replace(/\D/g, '');

  // Format: (11) 98765-4321 or (11) 3456-7890
  if (digits.length === 11) {
    return `(${digits.slice(0, 2)}) ${digits.slice(2, 7)}-${digits.slice(7)}`;
  } else if (digits.length === 10) {
    return `(${digits.slice(0, 2)}) ${digits.slice(2, 6)}-${digits.slice(6)}`;
  }

  return phone;
}

// Property type labels
export const propertyTypeLabels: Record<string, string> = {
  apartment: 'Apartamento',
  house: 'Casa',
  condo: 'Condomínio',
  commercial: 'Comercial',
  land: 'Terreno',
  farm: 'Chácara/Sítio',
  studio: 'Studio',
  penthouse: 'Cobertura',
  townhouse: 'Sobrado',
};

// Transaction type labels
export const transactionTypeLabels: Record<string, string> = {
  sale: 'Venda',
  rent: 'Aluguel',
  both: 'Venda/Aluguel',
};

// Status labels
export const statusLabels: Record<string, string> = {
  available: 'Disponível',
  sold: 'Vendido',
  rented: 'Alugado',
  reserved: 'Reservado',
  unavailable: 'Indisponível',
};

// Get property type label
export function getPropertyTypeLabel(type: string): string {
  return propertyTypeLabels[type] || type;
}

// Get transaction type label
export function getTransactionTypeLabel(type: string): string {
  return transactionTypeLabels[type] || type;
}

// Get status label
export function getStatusLabel(status: string): string {
  return statusLabels[status] || status;
}

// Build WhatsApp URL
export function buildWhatsAppUrl(phone: string, message: string): string {
  const cleanPhone = phone.replace(/\D/g, '');
  const encodedMessage = encodeURIComponent(message);
  return `https://wa.me/55${cleanPhone}?text=${encodedMessage}`;
}

// Generate slug from title
export function generateSlug(title: string): string {
  return title
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '') // Remove accents
    .replace(/[^\w\s-]/g, '') // Remove special chars
    .replace(/\s+/g, '-') // Replace spaces with -
    .replace(/-+/g, '-') // Replace multiple - with single -
    .trim();
}

// Truncate text
export function truncate(text: string, maxLength: number): string {
  if (text.length <= maxLength) return text;
  return text.substring(0, maxLength) + '...';
}

// Format date
export function formatDate(date: Date | string | undefined): string {
  if (!date) return '';

  const d = typeof date === 'string' ? new Date(date) : date;

  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  }).format(d);
}

// Format datetime
export function formatDateTime(date: Date | string | undefined): string {
  if (!date) return '';

  const d = typeof date === 'string' ? new Date(date) : date;

  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(d);
}

// Format relative time (e.g., "há 2 dias")
export function formatRelativeTime(date: Date | string | undefined): string {
  if (!date) return '';

  const d = typeof date === 'string' ? new Date(date) : date;
  const now = new Date();
  const diffMs = now.getTime() - d.getTime();
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffDays === 0) return 'Hoje';
  if (diffDays === 1) return 'Ontem';
  if (diffDays < 7) return `Há ${diffDays} dias`;
  if (diffDays < 30) return `Há ${Math.floor(diffDays / 7)} semanas`;
  if (diffDays < 365) return `Há ${Math.floor(diffDays / 30)} meses`;
  return `Há ${Math.floor(diffDays / 365)} anos`;
}

// Validate email
export function isValidEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

// Validate phone (BR)
export function isValidPhone(phone: string): boolean {
  const digits = phone.replace(/\D/g, '');
  return digits.length === 10 || digits.length === 11;
}

// Get property features as array
export function getPropertyFeatures(property: any): string[] {
  const features: string[] = [];

  if (property.bedrooms) features.push(`${property.bedrooms} quartos`);
  if (property.bathrooms) features.push(`${property.bathrooms} banheiros`);
  if (property.suites) features.push(`${property.suites} suítes`);
  if (property.parking_spaces) features.push(`${property.parking_spaces} vagas`);
  if (property.area_sqm) features.push(formatArea(property.area_sqm));

  return features;
}

// Get property amenities as array
export function getPropertyAmenities(property: any): string[] {
  const amenities: string[] = [];

  if (property.furnished) amenities.push('Mobiliado');
  if (property.pet_friendly) amenities.push('Aceita pets');
  if (property.has_pool) amenities.push('Piscina');
  if (property.has_gym) amenities.push('Academia');
  if (property.has_elevator) amenities.push('Elevador');
  if (property.has_security) amenities.push('Segurança 24h');
  if (property.has_garden) amenities.push('Jardim');
  if (property.has_balcony) amenities.push('Varanda');

  return amenities;
}

// Format file size
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes';

  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
}

// Format number with thousand separator
export function formatNumber(value: number | undefined): string {
  if (!value) return '0';
  return new Intl.NumberFormat('pt-BR').format(value);
}

// Format percentage
export function formatPercentage(value: number | undefined): string {
  if (!value) return '0%';
  return `${value.toFixed(1)}%`;
}
