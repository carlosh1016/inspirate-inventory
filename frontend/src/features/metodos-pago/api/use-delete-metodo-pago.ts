import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';

export function useDeleteMetodoPago() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      await api.delete(`/metodos-pago/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['metodos-pago'] });
    },
  });
}
