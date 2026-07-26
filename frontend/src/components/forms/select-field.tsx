'use client';

import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { cn } from '@/lib/utils';
import { FormError } from './form-error';

export interface SelectOption {
  value: string;
  label: string;
}

interface Props {
  label?: string;
  value: string;
  onChange: (value: string) => void;
  options: SelectOption[];
  placeholder?: string;
  error?: string;
  disabled?: boolean;
  className?: string;
  triggerClassName?: string;
}

// Label + base-ui Select + error. The `items` map lets SelectValue render the
// selected option's label (not its raw value).
export function SelectField({
  label,
  value,
  onChange,
  options,
  placeholder = 'Selecciona…',
  error,
  disabled,
  className,
  triggerClassName,
}: Props) {
  const items = Object.fromEntries(options.map((o) => [o.value, o.label]));

  return (
    <div className={cn('space-y-2', className)}>
      {label && <Label>{label}</Label>}
      <Select
        items={items}
        value={value}
        onValueChange={(next) => onChange((next as string) ?? '')}
        disabled={disabled}
      >
        <SelectTrigger className={cn('w-full', triggerClassName)}>
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
      <FormError message={error} />
    </div>
  );
}
