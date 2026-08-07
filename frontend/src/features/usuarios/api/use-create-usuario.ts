import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { UsuarioApi } from '@/types/domain';
import type { CreateUsuarioInput } from '../schemas/usuario-schema';

export function useCreateUsuario() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: CreateUsuarioInput) => {
      const res = await api.post<ApiEnvelope<UsuarioApi>>('/usuarios', {
        nombre_completo: input.nombre_completo,
        correo: input.correo,
        password: input.password,
        rol: input.rol,
      });
      return res.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['usuarios'] });
    },
  });
}
