import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';

export function useDeleteProducto() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      await api.delete(`/productos/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['productos'] });
      queryClient.invalidateQueries({ queryKey: ['stock'] });
    },
  });
}
