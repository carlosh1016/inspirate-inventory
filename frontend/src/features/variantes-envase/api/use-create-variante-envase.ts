import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { VarianteEnvase } from '../types';

export interface CreateVarianteEnvasePayload {
  modelo_envase_id: number;
  color: string;
  stock_minimo: number;
}

export function useCreateVarianteEnvase() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: CreateVarianteEnvasePayload) => {
      const res = await api.post<ApiEnvelope<VarianteEnvase>>('/variantes-envase', input);
      return res.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['variantes-envase'] });
      queryClient.invalidateQueries({ queryKey: ['modelos-envase'] });
      queryClient.invalidateQueries({ queryKey: ['stock'] });
    },
  });
}
