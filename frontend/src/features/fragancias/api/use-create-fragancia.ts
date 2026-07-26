import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { Fragancia } from '../types';

export interface FraganciaPayload {
  nombre_comercial: string;
  nombre_alternativo: string | null;
  genero: 'masculina' | 'femenina';
  gramos_minimo: string;
}

export function useCreateFragancia() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: FraganciaPayload) => {
      const res = await api.post<ApiEnvelope<Fragancia>>('/fragancias', input);
      return res.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['fragancias'] });
      queryClient.invalidateQueries({ queryKey: ['stock'] });
    },
  });
}
