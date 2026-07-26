'use client';

import * as React from 'react';

import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';

interface Props extends Omit<React.ComponentProps<'input'>, 'value' | 'onChange'> {
  /** Raw integer string (what the backend expects), e.g. "17000". */
  value: string;
  onChange: (value: string) => void;
}

function formatThousands(raw: string): string {
  if (!raw) return '';
  const n = Number.parseInt(raw, 10);
  return Number.isNaN(n) ? '' : new Intl.NumberFormat('es-CO').format(n);
}

// COP input: stores raw digits, shows thousands separators when not focused.
// No decimals (COP has none in this shop).
export function MoneyInput({ value, onChange, onFocus, onBlur, className, ...props }: Props) {
  const [focused, setFocused] = React.useState(false);
  const display = focused ? value : formatThousands(value);

  return (
    <div className="relative">
      <span className="pointer-events-none absolute top-1/2 left-2.5 -translate-y-1/2 text-sm text-muted-foreground">
        $
      </span>
      <Input
        inputMode="numeric"
        className={cn('pl-6', className)}
        value={display}
        onChange={(e) => onChange(e.target.value.replace(/\D/g, ''))}
        onFocus={(e) => {
          setFocused(true);
          onFocus?.(e);
        }}
        onBlur={(e) => {
          setFocused(false);
          onBlur?.(e);
        }}
        {...props}
      />
    </div>
  );
}
