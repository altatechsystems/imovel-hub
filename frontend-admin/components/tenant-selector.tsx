'use client';

import { useState, useEffect, useRef } from 'react';
import { Building2, ChevronDown, Check } from 'lucide-react';
import { useTenant } from '@/contexts/tenant-context';

export function TenantSelector() {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const {
    isPlatformAdmin,
    availableTenants,
    currentTenant,
    effectiveTenantId,
    setSelectedTenantId,
    isLoadingTenants,
  } = useTenant();

  useEffect(() => {
    // Close dropdown when clicking outside
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    }

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
      return () => document.removeEventListener('mousedown', handleClickOutside);
    }
  }, [isOpen]);

  const handleTenantChange = (tenantId: string) => {
    setSelectedTenantId(tenantId);
    setIsOpen(false);
  };

  // Don't show selector if not platform admin
  if (!isPlatformAdmin) {
    return null;
  }

  return (
    <div className="relative" ref={dropdownRef}>
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center gap-2 px-3 py-2 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
      >
        <Building2 className="w-4 h-4 text-gray-600" />
        <div className="text-left">
          <p className="text-xs text-gray-500">Tenant</p>
          <p className="text-sm font-medium text-gray-900">
            {isLoadingTenants ? 'Carregando...' : currentTenant?.name || 'Platform Admin'}
          </p>
        </div>
        <ChevronDown className={`w-4 h-4 text-gray-600 transition-transform ${isOpen ? 'rotate-180' : ''}`} />
      </button>

      {isOpen && (
        <div className="absolute top-full right-0 mt-2 w-72 bg-white border border-gray-200 rounded-lg shadow-lg z-50">
          <div className="p-2">
            <div className="px-3 py-2 text-xs font-semibold text-gray-500 uppercase">
              Selecionar Tenant
            </div>
            <div className="max-h-96 overflow-y-auto">
              {/* Platform Admin Option */}
              <button
                onClick={() => handleTenantChange('tenant_master')}
                className={`w-full flex items-center justify-between px-3 py-2 text-left rounded-lg transition-colors ${
                  effectiveTenantId === 'tenant_master'
                    ? 'bg-blue-50 text-blue-900'
                    : 'hover:bg-gray-50 text-gray-900'
                }`}
              >
                <div className="flex-1">
                  <p className="text-sm font-medium">Platform Admin</p>
                  <p className="text-xs text-gray-500">Visualizar como administrador</p>
                </div>
                {effectiveTenantId === 'tenant_master' && (
                  <Check className="w-4 h-4 text-blue-600" />
                )}
              </button>

              {availableTenants.length === 0 ? (
                <div className="px-3 py-4 text-sm text-gray-500 text-center">
                  Nenhum tenant disponível
                </div>
              ) : (
                availableTenants
                  .filter(tenant => tenant.id !== 'tenant_master')
                  .map((tenant) => (
                    <button
                      key={tenant.id}
                      onClick={() => handleTenantChange(tenant.id)}
                      className={`w-full flex items-center justify-between px-3 py-2 text-left rounded-lg transition-colors ${
                        effectiveTenantId === tenant.id
                          ? 'bg-blue-50 text-blue-900'
                          : 'hover:bg-gray-50 text-gray-900'
                      }`}
                    >
                      <div className="flex-1">
                        <p className="text-sm font-medium">{tenant.name}</p>
                        <p className="text-xs text-gray-500">{tenant.slug}</p>
                      </div>
                      {effectiveTenantId === tenant.id && (
                        <Check className="w-4 h-4 text-blue-600" />
                      )}
                      {!tenant.is_active && (
                        <span className="ml-2 px-2 py-0.5 bg-red-100 text-red-600 text-xs rounded-full">
                          Inativo
                        </span>
                      )}
                    </button>
                  ))
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
