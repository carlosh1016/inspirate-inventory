import Image from 'next/image';

export default function PublicLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex min-h-screen items-center justify-center bg-muted/40 p-4">
      <div className="w-full max-w-md">
        <div className="mb-8 text-center">
          <Image
            src="/inspirate-logo.jpg"
            alt="Inspírate Perfumes & Cosmética"
            width={1542}
            height={688}
            priority
            className="mx-auto h-auto w-64"
          />
          <p className="mt-2 text-sm text-muted-foreground">Sistema de inventario</p>
        </div>
        {children}
      </div>
    </div>
  );
}
