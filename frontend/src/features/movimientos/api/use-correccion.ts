import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import { invalidateInventory } from './invalidate';
import type { AjustePayload, AjusteResult } from './use-ajuste';

// Correccion shares the payload/result shape with Ajuste, different endpoint.
export function useCorreccion() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (payload: AjustePayload): Promise<AjusteResult> => {
      const res = await api.post('/movimientos/correccion', payload);
      const data = res.data.data;
      if (data && typeof data === 'object' && 'mensaje' in data) {
        return { movimiento: null, mensaje: data.mensaje as string };
      }
      return { movimiento: data };
    },
    onSuccess: () => invalidateInventory(queryClient),
  });
}
