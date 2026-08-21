import { api } from '@/lib/api';
import type { ComboboxOption } from '@/components/forms/combobox';
import type { ApiListEnvelope } from '@/types/api';
import type { ModeloEnvase } from '../types';

export function modeloEnvaseLabel(m: Pick<ModeloEnvase, 'tipo' | 'tamano_oz'>): string {
  return `${m.tipo} · ${m.tamano_oz} oz`;
}

// Combobox searchFn for picking a modelo when creating a variante: active
// modelos matching the query, excluding modelos marked tiene_variantes=false
// (e.g. "envase de lujo") — those already have their single hidden variante
// and never need another one added manually.
export async function searchModelosEnvase(query: string): Promise<ComboboxOption[]> {
  const res = await api.get<ApiListEnvelope<ModeloEnvase>>('/modelos-envase', {
    params: { q: query, activo: 'true', page_size: 20 },
  });
  return res.data.data.filter((m) => m.tiene_variantes).map((m) => ({ id: m.id, label: modeloEnvaseLabel(m) }));
}
