import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiListEnvelope } from '@/types/api';
import type { CuadreListItem } from '../types';

// Admin-only, reload-safe detection of a previous day's cuadre left open. Used
// to show the persistent warning banner on /caja. Returns the earliest open
// cuadre whose fecha is before today, or null.
export function useCuadreAnteriorAbierto(enabled: boolean) {
  return useQuery({
    queryKey: ['cuadres-caja', 'anterior-abierto'],
    queryFn: async () => {
      const res = await api.get<ApiListEnvelope<CuadreListItem>>('/cuadres-caja', {
        params: { estado: 'abierto', page_size: 50 },
      });
      const hoy = new Date().toISOString().slice(0, 10);
      const anteriores = res.data.data
        .filter((c) => c.fecha < hoy)
        .sort((a, b) => a.fecha.localeCompare(b.fecha));
      return anteriores[0] ?? null;
    },
    enabled,
  });
}
