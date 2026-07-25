import { Label } from '@/components/ui/label';
import { FormError } from './form-error';

interface Props {
  id: string;
  label: string;
  error?: string;
  children: React.ReactNode;
}

// Label + control + error message, wired by a shared id for accessibility.
export function FormField({ id, label, error, children }: Props) {
  return (
    <div className="space-y-2">
      <Label htmlFor={id}>{label}</Label>
      {children}
      <FormError message={error} />
    </div>
  );
}
