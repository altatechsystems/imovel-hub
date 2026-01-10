'use client';

import { ReactNode, useState } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ClientOnlyAuthProvider } from '@/components/client-only-auth-provider';
import { TenantProvider } from '@/contexts/tenant-context';

export interface ProvidersProps {
  children: ReactNode;
}

export function Providers({ children }: ProvidersProps) {
  // Create QueryClient inside component to avoid SSR issues
  // This ensures a new client is created for each request during SSR
  const [queryClient] = useState(() => new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
        retry: 1,
        staleTime: 5 * 60 * 1000, // 5 minutes
      },
    },
  }));

  return (
    <ClientOnlyAuthProvider>
      <TenantProvider>
        <QueryClientProvider client={queryClient}>
          {children}
        </QueryClientProvider>
      </TenantProvider>
    </ClientOnlyAuthProvider>
  );
}
