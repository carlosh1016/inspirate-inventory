import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { MetodoPagoInput } from '../schemas/metodo-pago-schema';
import type { MetodoPago } from '../types';

export function useUpdateMetodoPago(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: Partial<MetodoPagoInput>) => {
      const res = await api.patch<ApiEnvelope<MetodoPago>>(`/metodos-pago/${id}`, input);
      return res.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['metodos-pago'] });
    },
  });
}
