import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { FraganciaPayload } from './use-create-fragancia';
import type { Fragancia } from '../types';

export function useUpdateFragancia(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: Partial<FraganciaPayload>) => {
      const res = await api.patch<ApiEnvelope<Fragancia>>(`/fragancias/${id}`, input);
      return res.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['fragancias'] });
      queryClient.invalidateQueries({ queryKey: ['stock'] });
    },
  });
}
