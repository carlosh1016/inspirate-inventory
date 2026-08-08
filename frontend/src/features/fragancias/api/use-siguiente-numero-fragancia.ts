import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';

// Fetched imperatively (on genero change in the create form), not via
// useQuery — it's a one-off UI suggestion, not data the page renders
// reactively, so there's nothing to keep in sync with a query cache.
export async function fetchSiguienteNumeroFragancia(genero: 'masculina' | 'femenina'): Promise<number> {
  const res = await api.get<ApiEnvelope<{ siguiente: number }>>('/fragancias/siguiente-numero', {
    params: { genero },
  });
  return res.data.data.siguiente;
}
