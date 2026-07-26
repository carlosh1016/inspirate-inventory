import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiListEnvelope } from '@/types/api';
import type { MetodoPago } from '@/features/metodos-pago/types';

// Active métodos de pago for the venta and cuadre-cierre selectors. Readable
// by admin and vendedora (backend GET /metodos-pago read is open to both).
export function useMetodosPagoActivos() {
  return useQuery({
    queryKey: ['metodos-pago', 'activos'],
    queryFn: async () => {
      const res = await api.get<ApiListEnvelope<MetodoPago>>('/metodos-pago', {
        params: { activo: 'true', page_size: 100 },
      });
      return res.data.data;
    },
  });
}
