'use client';

import { ReactNode, useEffect, useState } from 'react';
import { AuthProvider } from '@/contexts/auth-context';

/**
 * Wrapper that ensures AuthProvider only renders on client side
 * Prevents SSR hydration issues with Firebase Auth
 */
export function ClientOnlyAuthProvider({ children }: { children: ReactNode }) {
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);
  }, []);

  // During SSR or initial render, don't render AuthProvider
  if (!isMounted) {
    return <>{children}</>;
  }

  // Only render AuthProvider after client-side mount
  return <AuthProvider>{children}</AuthProvider>;
}
