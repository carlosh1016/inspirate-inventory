import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { UsuarioApi } from '@/types/domain';

export interface UpdateUsuarioInput {
  nombre_completo: string;
  rol: 'admin' | 'vendedora';
}

export function useUpdateUsuario(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: UpdateUsuarioInput) => {
      const res = await api.patch<ApiEnvelope<UsuarioApi>>(`/usuarios/${id}`, input);
      return res.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['usuarios'] });
    },
  });
}
