'use client';

import * as React from 'react';
import { Check, ChevronsUpDown, Loader2, X } from 'lucide-react';

import { Input } from '@/components/ui/input';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { useDebounce } from '@/hooks/use-debounce';
import { cn } from '@/lib/utils';

export interface ComboboxOption {
  id: number;
  label: string;
}

interface Props<O extends ComboboxOption> {
  value: number | null;
  /** Label of the current selection, shown on the trigger. */
  selectedLabel?: string | null;
  onChange: (id: number | null, label: string | null, option: O | null) => void;
  searchFn: (query: string) => Promise<O[]>;
  placeholder?: string;
  renderOption?: (option: O) => React.ReactNode;
  emptyMessage?: string;
  disabled?: boolean;
  clearable?: boolean;
  className?: string;
}

// Async single-select combobox built on Popover + Input (not cmdk): the search
// is server-side, so options are fetched with a 300ms debounce whenever the
// popup is open. Selecting an option stores its id and label.
export function Combobox<O extends ComboboxOption>({
  value,
  selectedLabel,
  onChange,
  searchFn,
  placeholder = 'Buscar…',
  renderOption,
  emptyMessage = 'No se encontraron resultados',
  disabled,
  clearable = true,
  className,
}: Props<O>) {
  const [open, setOpen] = React.useState(false);
  const [query, setQuery] = React.useState('');
  const [options, setOptions] = React.useState<O[]>([]);
  const [loading, setLoading] = React.useState(false);
  const debounced = useDebounce(query, 300);

  // Keep searchFn in a ref so an inline callback from the caller doesn't
  // re-trigger the fetch effect on every render.
  const searchRef = React.useRef(searchFn);
  React.useEffect(() => {
    searchRef.current = searchFn;
  });

  // `loading` is turned on from the event handlers (typing / opening) so it is
  // never set synchronously inside this effect; the effect only runs the fetch
  // and clears loading in its async callbacks.
  React.useEffect(() => {
    if (!open) return;
    let active = true;
    searchRef
      .current(debounced)
      .then((res) => {
        if (active) setOptions(res);
      })
      .catch(() => {
        if (active) setOptions([]);
      })
      .finally(() => {
        if (active) setLoading(false);
      });
    return () => {
      active = false;
    };
  }, [debounced, open]);

  const handleOpenChange = (next: boolean) => {
    setOpen(next);
    if (next) setLoading(true);
  };

  const handleSelect = (option: O) => {
    onChange(option.id, option.label, option);
    setOpen(false);
    setQuery('');
  };

  return (
    <Popover open={open} onOpenChange={handleOpenChange}>
      <PopoverTrigger
        disabled={disabled}
        className={cn(
          'flex h-8 w-full items-center justify-between gap-1.5 rounded-lg border border-input bg-transparent px-2.5 text-sm outline-none transition-colors focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50',
          className,
        )}
      >
        <span className={cn('truncate', value === null && 'text-muted-foreground')}>
          {value !== null ? (selectedLabel ?? `#${value}`) : placeholder}
        </span>
        <span className="flex items-center gap-1">
          {clearable && value !== null && (
            <X
              className="size-3.5 text-muted-foreground hover:text-foreground"
              onClick={(e) => {
                e.stopPropagation();
                onChange(null, null, null);
              }}
            />
          )}
          <ChevronsUpDown className="size-4 shrink-0 text-muted-foreground" />
        </span>
      </PopoverTrigger>
      <PopoverContent align="start" className="w-(--anchor-width) min-w-72 p-1.5">
        <Input
          autoFocus
          value={query}
          onChange={(e) => {
            setQuery(e.target.value);
            setLoading(true);
          }}
          placeholder="Escribe para buscar…"
          className="mb-1.5"
        />
        <div className="max-h-64 overflow-y-auto">
          {loading ? (
            <div className="flex items-center justify-center gap-2 py-6 text-sm text-muted-foreground">
              <Loader2 className="size-4 animate-spin" />
              Buscando…
            </div>
          ) : options.length === 0 ? (
            <p className="py-6 text-center text-sm text-muted-foreground">{emptyMessage}</p>
          ) : (
            options.map((option) => (
              <button
                key={option.id}
                type="button"
                onClick={() => handleSelect(option)}
                className="flex w-full items-center justify-between gap-2 rounded-md px-2 py-1.5 text-left text-sm hover:bg-muted"
              >
                <span className="min-w-0 flex-1">
                  {renderOption ? renderOption(option) : option.label}
                </span>
                {value === option.id && <Check className="size-4 shrink-0 text-primary" />}
              </button>
            ))
          )}
        </div>
      </PopoverContent>
    </Popover>
  );
}
