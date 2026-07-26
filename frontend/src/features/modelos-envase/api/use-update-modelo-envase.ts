import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { ModeloEnvasePayload } from './use-create-modelo-envase';
import type { ModeloEnvase } from '../types';

export function useUpdateModeloEnvase(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: Partial<ModeloEnvasePayload>) => {
      const res = await api.patch<ApiEnvelope<ModeloEnvase>>(`/modelos-envase/${id}`, input);
      return res.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['modelos-envase'] });
    },
  });
}
