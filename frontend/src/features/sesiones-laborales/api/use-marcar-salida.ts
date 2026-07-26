import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { Sesion } from '../types';

export function useMarcarSalida() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async () => {
      const res = await api.post<ApiEnvelope<Sesion>>('/sesiones-laborales/salida');
      return res.data.data;
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['sesiones-laborales'] }),
  });
}
