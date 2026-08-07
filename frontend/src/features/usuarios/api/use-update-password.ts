import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';

// Admin resetting another usuario's password: password_actual is only
// required (and validated) when target === requester, never the case here.
export function useUpdatePassword(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (passwordNueva: string) => {
      await api.patch(`/usuarios/${id}/password`, { password_nueva: passwordNueva });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['usuarios'] });
    },
  });
}
