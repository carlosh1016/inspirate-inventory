import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { VarianteEnvase } from '../types';

export interface UpdateVarianteEnvasePayload {
  color?: string;
  stock_minimo?: number;
}

export function useUpdateVarianteEnvase(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: UpdateVarianteEnvasePayload) => {
      const res = await api.patch<ApiEnvelope<VarianteEnvase>>(`/variantes-envase/${id}`, input);
      return res.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['variantes-envase'] });
      queryClient.invalidateQueries({ queryKey: ['stock'] });
    },
  });
}
