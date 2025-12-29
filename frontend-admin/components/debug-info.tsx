'use client';

import { useEffect, useState } from 'react';

export function DebugInfo() {
  const [info, setInfo] = useState<Record<string, string>>({});
  const [isPlatformAdmin, setIsPlatformAdmin] = useState(false);

  useEffect(() => {
    const debugInfo = {
      tenant_id: localStorage.getItem('tenant_id') || 'not set',
      broker_id: localStorage.getItem('broker_id') || 'not set',
      broker_role: localStorage.getItem('broker_role') || 'not set',
      is_platform_admin: localStorage.getItem('is_platform_admin') || 'not set',
    };
    setInfo(debugInfo);
    setIsPlatformAdmin(localStorage.getItem('is_platform_admin') === 'true');
  }, []);

  // Only show for platform admins
  if (!isPlatformAdmin) {
    return null;
  }

  return (
    <div className="fixed bottom-4 right-4 bg-black/80 text-white p-4 rounded-lg text-xs font-mono z-50 max-w-xs">
      <div className="font-bold mb-2">Debug Info (Platform Admin)</div>
      {Object.entries(info).map(([key, value]) => (
        <div key={key}>
          <span className="text-blue-300">{key}:</span> {value}
        </div>
      ))}
    </div>
  );
}
