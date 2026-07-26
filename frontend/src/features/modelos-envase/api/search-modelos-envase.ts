import { api } from '@/lib/api';
import type { ComboboxOption } from '@/components/forms/combobox';
import type { ApiListEnvelope } from '@/types/api';
import type { ModeloEnvase } from '../types';

export function modeloEnvaseLabel(m: Pick<ModeloEnvase, 'tipo' | 'tamano_oz'>): string {
  return `${m.tipo} · ${m.tamano_oz} oz`;
}

// Combobox searchFn: active modelos matching the query.
export async function searchModelosEnvase(query: string): Promise<ComboboxOption[]> {
  const res = await api.get<ApiListEnvelope<ModeloEnvase>>('/modelos-envase', {
    params: { q: query, activo: 'true', page_size: 20 },
  });
  return res.data.data.map((m) => ({ id: m.id, label: modeloEnvaseLabel(m) }));
}
