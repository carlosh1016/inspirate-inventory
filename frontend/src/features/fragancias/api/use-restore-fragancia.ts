import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';

// Un-soft-deletes a fragancia (POST /fragancias/:id/restaurar). Admin-only on
// the backend.
export function useRestoreFragancia() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      await api.post(`/fragancias/${id}/restaurar`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['fragancias'] });
      queryClient.invalidateQueries({ queryKey: ['stock'] });
    },
  });
}
