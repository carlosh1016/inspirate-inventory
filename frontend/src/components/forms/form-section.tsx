import { cn } from '@/lib/utils';

interface Props {
  title?: string;
  description?: string;
  className?: string;
  children: React.ReactNode;
}

// Visual grouping for a set of related form fields.
export function FormSection({ title, description, className, children }: Props) {
  return (
    <section className={cn('space-y-4', className)}>
      {(title || description) && (
        <div className="space-y-1">
          {title && <h2 className="text-sm font-medium">{title}</h2>}
          {description && <p className="text-sm text-muted-foreground">{description}</p>}
        </div>
      )}
      {children}
    </section>
  );
}
