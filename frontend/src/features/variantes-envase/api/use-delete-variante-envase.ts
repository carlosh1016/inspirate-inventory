import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';

export function useDeleteVarianteEnvase() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      await api.delete(`/variantes-envase/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['variantes-envase'] });
      queryClient.invalidateQueries({ queryKey: ['modelos-envase'] });
      queryClient.invalidateQueries({ queryKey: ['stock'] });
    },
  });
}
