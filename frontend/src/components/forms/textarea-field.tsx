import * as React from 'react';

import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { FormError } from './form-error';

interface Props extends React.ComponentProps<'textarea'> {
  id: string;
  label: string;
  error?: string;
}

export function TextareaField({ id, label, error, ...props }: Props) {
  return (
    <div className="space-y-2">
      <Label htmlFor={id}>{label}</Label>
      <Textarea id={id} {...props} />
      <FormError message={error} />
    </div>
  );
}
