'use client';

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useAuth } from '@/hooks/use-auth';
import { adminApi } from '@/lib/api';

interface Tenant {
  id: string;
  name: string;
  slug: string;
  is_active: boolean;
}

interface TenantContextType {
  selectedTenantId: string | null;
  setSelectedTenantId: (tenantId: string | null) => void;
  effectiveTenantId: string | null;
  availableTenants: Tenant[];
  currentTenant: Tenant | null;
  isLoadingTenants: boolean;
  isPlatformAdmin: boolean;
}

const TenantContext = createContext<TenantContextType | undefined>(undefined);

export function TenantProvider({ children }: { children: ReactNode }) {
  const { tenantId, userProfile } = useAuth();
  const [selectedTenantId, setSelectedTenantId] = useState<string | null>(null);
  const [availableTenants, setAvailableTenants] = useState<Tenant[]>([]);
  const [isLoadingTenants, setIsLoadingTenants] = useState(false);

  // Check if user is platform admin
  const isPlatformAdmin = userProfile?.role === 'admin' && tenantId === 'tenant_master';

  // Determine which tenant to use: selected or user's tenant
  const effectiveTenantId = selectedTenantId || tenantId;

  // Find current tenant object
  const currentTenant = availableTenants.find(t => t.id === effectiveTenantId) || null;

  // Update API client when effective tenant changes
  useEffect(() => {
    if (effectiveTenantId) {
      console.log('[TenantContext] Updating API client tenant to:', effectiveTenantId);
      adminApi.setTenant(effectiveTenantId);
    }
  }, [effectiveTenantId]);

  // Fetch tenants if platform admin
  useEffect(() => {
    const fetchTenants = async () => {
      if (!isPlatformAdmin) return;

      try {
        setIsLoadingTenants(true);
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL?.replace('/api/v1', '')}/tenants`);

        if (response.ok) {
          const data = await response.json();
          setAvailableTenants(data.data || []);
        }
      } catch (error) {
        console.error('Error fetching tenants:', error);
      } finally {
        setIsLoadingTenants(false);
      }
    };

    fetchTenants();
  }, [isPlatformAdmin]);

  return (
    <TenantContext.Provider
      value={{
        selectedTenantId,
        setSelectedTenantId,
        effectiveTenantId,
        availableTenants,
        currentTenant,
        isLoadingTenants,
        isPlatformAdmin,
      }}
    >
      {children}
    </TenantContext.Provider>
  );
}

export function useTenant() {
  const context = useContext(TenantContext);
  if (context === undefined) {
    throw new Error('useTenant must be used within a TenantProvider');
  }
  return context;
}
