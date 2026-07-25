import { useMutation } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ForgotPasswordInput } from '../schemas/forgot-password-schema';

export function useRequestPasswordReset() {
  return useMutation({
    mutationFn: async (input: ForgotPasswordInput) => {
      await api.post('/auth/password-reset/request', input);
    },
  });
}
