'use client';

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { cn } from '@/lib/utils';

interface Option {
  value: string;
  label: string;
}

interface Props {
  value: string;
  onChange: (value: string) => void;
  options: Option[];
  placeholder?: string;
  ariaLabel?: string;
  className?: string;
}

// Compact select for toolbars. Options should include the "all" sentinel
// (e.g. { value: 'all', label: 'Todos' }); the caller's serializer keeps it
// out of the URL when it's the default.
export function SelectFilter({ value, onChange, options, placeholder, ariaLabel, className }: Props) {
  const items = Object.fromEntries(options.map((o) => [o.value, o.label]));

  return (
    <Select
      items={items}
      value={value}
      onValueChange={(next) => onChange((next as string) ?? '')}
    >
      <SelectTrigger aria-label={ariaLabel} className={cn('min-w-36', className)}>
        <SelectValue placeholder={placeholder} />
      </SelectTrigger>
      <SelectContent>
        {options.map((option) => (
          <SelectItem key={option.value} value={option.value}>
            {option.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
