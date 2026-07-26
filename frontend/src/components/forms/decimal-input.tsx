'use client';

import * as React from 'react';

import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';

interface Props extends Omit<React.ComponentProps<'input'>, 'value' | 'onChange'> {
  /** Decimal string using a dot, e.g. "30" or "30.5". */
  value: string;
  onChange: (value: string) => void;
  suffix?: string;
}

function sanitize(input: string): string {
  const cleaned = input.replace(/[^\d.]/g, '');
  const parts = cleaned.split('.');
  if (parts.length <= 1) return cleaned;
  return `${parts[0]}.${parts.slice(1).join('')}`;
}

// Numeric input allowing one decimal point, with an optional suffix (e.g. "g").
export function DecimalInput({ value, onChange, suffix, className, ...props }: Props) {
  return (
    <div className="relative">
      <Input
        inputMode="decimal"
        className={cn(suffix && 'pr-9', className)}
        value={value}
        onChange={(e) => onChange(sanitize(e.target.value))}
        {...props}
      />
      {suffix && (
        <span className="pointer-events-none absolute top-1/2 right-2.5 -translate-y-1/2 text-sm text-muted-foreground">
          {suffix}
        </span>
      )}
    </div>
  );
}
