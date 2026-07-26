import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { MetodoPagoInput } from '../schemas/metodo-pago-schema';
import type { MetodoPago } from '../types';

export function useCreateMetodoPago() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: MetodoPagoInput) => {
      const res = await api.post<ApiEnvelope<MetodoPago>>('/metodos-pago', input);
      return res.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['metodos-pago'] });
    },
  });
}
