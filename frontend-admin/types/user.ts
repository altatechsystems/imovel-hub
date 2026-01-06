// User types for administrative users (NOT brokers)
// Administrative users do NOT have CRECI and are stored in /tenants/{id}/users/
// For brokers (with CRECI), see broker.ts

export type UserRole = 'admin' | 'manager';

export interface User {
  id: string;
  tenant_id: string;
  firebase_uid: string;
  name: string;
  email: string;
  phone?: string;
  document?: string;
  document_type?: 'cpf' | 'cnpj';
  role: UserRole;
  is_active: boolean;
  permissions?: string[];
  photo_url?: string;
  created_at: string | Date;
  updated_at: string | Date;
}

export interface CreateUserRequest {
  firebase_uid: string;
  name: string;
  email: string;
  phone?: string;
  document?: string;
  document_type?: 'cpf' | 'cnpj';
  role?: UserRole; // Defaults to 'admin' if not provided
  is_active?: boolean; // Defaults to true if not provided
  permissions?: string[];
  photo_url?: string;
}

export interface UpdateUserRequest {
  name?: string;
  email?: string;
  phone?: string;
  document?: string;
  document_type?: 'cpf' | 'cnpj';
  role?: UserRole;
  is_active?: boolean;
  photo_url?: string;
}

export interface GrantPermissionRequest {
  permission: string;
}

// Standard permissions for manager role
export const STANDARD_PERMISSIONS = {
  // Property permissions
  PROPERTY_VIEW: 'property.view',
  PROPERTY_CREATE: 'property.create',
  PROPERTY_UPDATE: 'property.update',
  PROPERTY_DELETE: 'property.delete',

  // Owner permissions
  OWNER_VIEW: 'owner.view',
  OWNER_CREATE: 'owner.create',
  OWNER_UPDATE: 'owner.update',
  OWNER_DELETE: 'owner.delete',

  // Broker permissions
  BROKER_VIEW: 'broker.view',
  BROKER_CREATE: 'broker.create',
  BROKER_UPDATE: 'broker.update',
  BROKER_DELETE: 'broker.delete',

  // Lead permissions
  LEAD_VIEW: 'lead.view',
  LEAD_CREATE: 'lead.create',
  LEAD_UPDATE: 'lead.update',
  LEAD_DELETE: 'lead.delete',

  // Listing permissions
  LISTING_VIEW: 'listing.view',
  LISTING_CREATE: 'listing.create',
  LISTING_UPDATE: 'listing.update',
  LISTING_DELETE: 'listing.delete',

  // User permissions (for managing other users)
  USER_VIEW: 'user.view',
  USER_CREATE: 'user.create',
  USER_UPDATE: 'user.update',
  USER_DELETE: 'user.delete',

  // Report permissions
  REPORT_VIEW: 'report.view',
  REPORT_EXPORT: 'report.export',

  // Settings permissions
  SETTINGS_VIEW: 'settings.view',
  SETTINGS_UPDATE: 'settings.update',
} as const;

export type Permission = typeof STANDARD_PERMISSIONS[keyof typeof STANDARD_PERMISSIONS];

// Helper function to check if user has permission
export function hasPermission(user: User, permission: Permission | string): boolean {
  // Admins have all permissions
  if (user.role === 'admin') {
    return true;
  }

  // Check if user has specific permission
  return user.permissions?.includes(permission) ?? false;
}

// Helper function to check if user is admin
export function isAdmin(user: User): boolean {
  return user.role === 'admin';
}

// Helper function to check if user is manager
export function isManager(user: User): boolean {
  return user.role === 'manager';
}

// Helper function to get role display name
export function getRoleDisplayName(role: UserRole): string {
  switch (role) {
    case 'admin':
      return 'Administrador';
    case 'manager':
      return 'Gerente';
    default:
      return role;
  }
}

// Helper function to get permission display name
export function getPermissionDisplayName(permission: string): string {
  const permissionMap: Record<string, string> = {
    'property.view': 'Visualizar imóveis',
    'property.create': 'Criar imóveis',
    'property.update': 'Editar imóveis',
    'property.delete': 'Excluir imóveis',
    'owner.view': 'Visualizar proprietários',
    'owner.create': 'Criar proprietários',
    'owner.update': 'Editar proprietários',
    'owner.delete': 'Excluir proprietários',
    'broker.view': 'Visualizar corretores',
    'broker.create': 'Criar corretores',
    'broker.update': 'Editar corretores',
    'broker.delete': 'Excluir corretores',
    'lead.view': 'Visualizar leads',
    'lead.create': 'Criar leads',
    'lead.update': 'Editar leads',
    'lead.delete': 'Excluir leads',
    'listing.view': 'Visualizar anúncios',
    'listing.create': 'Criar anúncios',
    'listing.update': 'Editar anúncios',
    'listing.delete': 'Excluir anúncios',
    'user.view': 'Visualizar usuários',
    'user.create': 'Criar usuários',
    'user.update': 'Editar usuários',
    'user.delete': 'Excluir usuários',
    'report.view': 'Visualizar relatórios',
    'report.export': 'Exportar relatórios',
    'settings.view': 'Visualizar configurações',
    'settings.update': 'Editar configurações',
  };

  return permissionMap[permission] || permission;
}
