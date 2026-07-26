import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';

export function useDeleteFragancia() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      await api.delete(`/fragancias/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['fragancias'] });
      queryClient.invalidateQueries({ queryKey: ['stock'] });
    },
  });
}
