import Link from 'next/link';
import { ArrowLeft } from 'lucide-react';

interface Props {
  title: string;
  description?: string;
  action?: React.ReactNode;
  backHref?: string;
  backLabel?: string;
}

// Standard page header: optional back link, title + description, and an action
// slot on the right (e.g. a primary "+ Nuevo" button).
export function PageHeader({ title, description, action, backHref, backLabel = 'Volver' }: Props) {
  return (
    <div className="mb-6 space-y-3">
      {backHref && (
        <Link
          href={backHref}
          className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
        >
          <ArrowLeft className="size-4" />
          {backLabel}
        </Link>
      )}
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold">{title}</h1>
          {description && <p className="text-sm text-muted-foreground">{description}</p>}
        </div>
        {action}
      </div>
    </div>
  );
}
