import type { LucideIcon } from 'lucide-react';

import { Card, CardContent } from '@/components/ui/card';
import { cn } from '@/lib/utils';

interface Props {
  icon: LucideIcon;
  label: string;
  value: string;
  variant?: 'default' | 'warning';
}

export function KpiCard({ icon: Icon, label, value, variant = 'default' }: Props) {
  return (
    <Card className={cn(variant === 'warning' && 'border-warning/30 bg-warning/5')}>
      <CardContent className="flex items-center gap-4 p-5">
        <div
          className={cn(
            'flex size-11 shrink-0 items-center justify-center rounded-md',
            variant === 'warning' ? 'bg-warning/10 text-warning' : 'bg-primary/10 text-primary',
          )}
        >
          <Icon size={20} />
        </div>
        <div className="min-w-0">
          <p className="truncate font-mono text-[11px] font-semibold tracking-wider text-muted-foreground uppercase">
            {label}
          </p>
          <p className="font-mono text-2xl font-bold tracking-tight">{value}</p>
        </div>
      </CardContent>
    </Card>
  );
}
