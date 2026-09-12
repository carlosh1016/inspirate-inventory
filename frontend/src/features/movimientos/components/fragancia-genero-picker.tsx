'use client';

import { useEffect, useState } from 'react';
import { Check } from 'lucide-react';

import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { SelectField } from '@/components/forms/select-field';
import { api } from '@/lib/api';
import type { ApiListEnvelope } from '@/types/api';
import type { Genero } from '@/types/domain';
import type { Fragancia } from '@/features/fragancias/types';

const GENERO_OPTIONS = [
  { value: 'femenina', label: 'Femenina' },
  { value: 'masculina', label: 'Masculina' },
];

interface Props {
  onChange: (id: number | null, label: string | null) => void;
  error?: string;
}

function isValid(genero: '' | Genero, numero: string): genero is Genero {
  const n = Number.parseInt(numero, 10);
  return !!genero && numero !== '' && !Number.isNaN(n) && n > 0;
}

// Encuentra una fragancia por Género + Número (el código "F-014" que ya
// conoce la vendedora de las hojas de Excel) en vez de tener que escribir
// texto en un buscador — dos campos estructurados que resuelven la
// fragancia exacta apenas ambos están completos. `status` se pone en
// 'loading'/'idle' desde los manejadores de los campos (no dentro del
// effect); el effect solo dispara la búsqueda y resuelve el estado final
// en sus callbacks async — mismo patrón que components/forms/combobox.tsx.
export function FraganciaGeneroPicker({ onChange, error }: Props) {
  const [genero, setGenero] = useState<'' | Genero>('');
  const [numero, setNumero] = useState('');
  const [status, setStatus] = useState<'idle' | 'loading' | 'found' | 'not-found'>('idle');
  const [foundLabel, setFoundLabel] = useState<string | null>(null);

  const handleGeneroChange = (v: string) => {
    const g = v as Genero;
    setGenero(g);
    setStatus(isValid(g, numero) ? 'loading' : 'idle');
    setFoundLabel(null);
    onChange(null, null);
  };

  const handleNumeroChange = (raw: string) => {
    const v = raw.replace(/\D/g, '');
    setNumero(v);
    setStatus(isValid(genero, v) ? 'loading' : 'idle');
    setFoundLabel(null);
    onChange(null, null);
  };

  useEffect(() => {
    if (!isValid(genero, numero)) return;
    const n = Number.parseInt(numero, 10);
    let cancelled = false;
    api
      .get<ApiListEnvelope<Fragancia>>('/fragancias', {
        params: { genero, numero_genero: n, activo: 'true', page_size: 1 },
      })
      .then((res) => {
        if (cancelled) return;
        const f = res.data.data[0];
        if (f) {
          setStatus('found');
          setFoundLabel(f.nombre_comercial);
          onChange(f.id, f.nombre_comercial);
        } else {
          setStatus('not-found');
          setFoundLabel(null);
        }
      })
      .catch(() => {
        if (cancelled) return;
        setStatus('not-found');
        setFoundLabel(null);
      });
    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [genero, numero]);

  return (
    <div className="grid grid-cols-2 gap-3">
      <SelectField
        label="Género"
        value={genero}
        onChange={handleGeneroChange}
        options={GENERO_OPTIONS}
        placeholder="Género…"
      />
      <div className="space-y-2">
        <Label htmlFor="fragancia-numero-genero">Número</Label>
        <Input
          id="fragancia-numero-genero"
          inputMode="numeric"
          placeholder="Ej. 14"
          value={numero}
          onChange={(e) => handleNumeroChange(e.target.value)}
        />
      </div>
      <div className="col-span-2 -mt-1 min-h-4 text-xs">
        {status === 'loading' && <p className="text-muted-foreground">Buscando…</p>}
        {status === 'found' && foundLabel && (
          <p className="flex items-center gap-1 text-success">
            <Check className="size-3.5" /> Fragancia encontrada: {foundLabel}
          </p>
        )}
        {status === 'not-found' && (
          <p className="text-destructive">No se encontró ninguna fragancia con ese género y número.</p>
        )}
        {error && status === 'idle' && <p className="text-destructive">{error}</p>}
      </div>
    </div>
  );
}
