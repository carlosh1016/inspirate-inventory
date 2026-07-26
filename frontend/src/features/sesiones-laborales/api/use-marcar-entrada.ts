import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { Sesion } from '../types';

export function useMarcarEntrada() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async () => {
      const res = await api.post<ApiEnvelope<Sesion>>('/sesiones-laborales/entrada');
      return res.data.data;
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['sesiones-laborales'] }),
  });
}
