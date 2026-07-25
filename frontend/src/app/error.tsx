'use client';

import Link from 'next/link';

import { Button, buttonVariants } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';

export default function GlobalError({ reset }: { error: Error & { digest?: string }; reset: () => void }) {
  return (
    <div className="flex min-h-screen items-center justify-center p-4">
      <Card className="w-full max-w-md text-center">
        <CardHeader>
          <CardTitle>Algo salió mal</CardTitle>
          <CardDescription>Ocurrió un error inesperado. Intenta de nuevo.</CardDescription>
        </CardHeader>
        <CardContent className="text-sm text-muted-foreground">
          Si el problema persiste, contacta al administrador.
        </CardContent>
        <CardFooter className="justify-center gap-2">
          <Button onClick={reset}>Recargar</Button>
          <Link href="/dashboard" className={buttonVariants({ variant: 'outline' })}>
            Ir al inicio
          </Link>
        </CardFooter>
      </Card>
    </div>
  );
}
