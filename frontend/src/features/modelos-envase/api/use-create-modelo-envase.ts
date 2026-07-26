import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { ModeloEnvase } from '../types';

export interface ModeloEnvasePayload {
  tipo: string;
  tamano_oz: string;
  equiv_gramos: string;
  precio_solo: string;
  precio_con_fragancia: string;
  precio_recarga: string;
}

export function useCreateModeloEnvase() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: ModeloEnvasePayload) => {
      const res = await api.post<ApiEnvelope<ModeloEnvase>>('/modelos-envase', input);
      return res.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['modelos-envase'] });
    },
  });
}
