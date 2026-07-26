'use client';

import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';

interface Props {
  id: string;
  label: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
}

export function BooleanFilter({ id, label, checked, onChange }: Props) {
  return (
    <div className="flex items-center gap-2">
      <Checkbox id={id} checked={checked} onCheckedChange={(next) => onChange(next === true)} />
      <Label htmlFor={id} className="cursor-pointer text-sm font-normal">
        {label}
      </Label>
    </div>
  );
}
