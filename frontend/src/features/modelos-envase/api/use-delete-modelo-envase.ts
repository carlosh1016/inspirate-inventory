import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';

export function useDeleteModeloEnvase() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      await api.delete(`/modelos-envase/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['modelos-envase'] });
      queryClient.invalidateQueries({ queryKey: ['variantes-envase'] });
    },
  });
}
