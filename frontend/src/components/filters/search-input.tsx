'use client';

import * as React from 'react';
import { Search } from 'lucide-react';

import { Input } from '@/components/ui/input';
import { useDebounce } from '@/hooks/use-debounce';
import { cn } from '@/lib/utils';

interface Props {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
}

// Text search with a 300ms debounce before propagating up (which typically
// writes to the URL). Local state keeps typing responsive.
export function SearchInput({ value, onChange, placeholder = 'Buscar…', className }: Props) {
  const [local, setLocal] = React.useState(value);
  const debounced = useDebounce(local, 300);

  React.useEffect(() => {
    if (debounced !== value) onChange(debounced);
  }, [debounced, value, onChange]);

  return (
    <div className={cn('relative', className)}>
      <Search className="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-muted-foreground" />
      <Input
        value={local}
        onChange={(e) => setLocal(e.target.value)}
        placeholder={placeholder}
        className="pl-8"
      />
    </div>
  );
}
