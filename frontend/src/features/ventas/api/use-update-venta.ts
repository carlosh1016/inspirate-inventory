import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { VentaDetallada } from '../types';

// Admin-only. Only the observaciones field is editable.
export function useUpdateVenta(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (observaciones: string | null) => {
      const res = await api.patch<ApiEnvelope<VentaDetallada>>(`/ventas/${id}`, { observaciones });
      return res.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['ventas'] });
    },
  });
}
