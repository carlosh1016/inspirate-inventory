import { useMutation } from '@tanstack/react-query';

import { api } from '@/lib/api';

interface ConfirmPayload {
  token: string;
  password_nueva: string;
}

export function useConfirmPasswordReset() {
  return useMutation({
    mutationFn: async (input: ConfirmPayload) => {
      await api.post('/auth/password-reset/confirm', input);
    },
  });
}
