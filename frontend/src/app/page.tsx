'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

// The root route has no auth context of its own (the access token lives only in
// the app area's AuthProvider), so it forwards to /dashboard, which validates
// the session and, if needed, redirects to /login.
export default function RootPage() {
  const router = useRouter();

  useEffect(() => {
    router.replace('/dashboard');
  }, [router]);

  return (
    <div className="flex min-h-screen items-center justify-center">
      <p className="text-sm text-muted-foreground">Cargando...</p>
    </div>
  );
}
