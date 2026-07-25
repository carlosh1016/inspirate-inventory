import { useMutation } from '@tanstack/react-query';

import { api } from '@/lib/api';
import { useAuthStore } from '@/stores/auth-store';
import type { ApiEnvelope } from '@/types/api';
import { toUsuarioSession, type UsuarioApi } from '@/types/domain';
import type { LoginInput } from '../schemas/login-schema';

interface LoginData {
  access_token: string;
  expires_in: number;
  usuario: UsuarioApi;
}

export function useLogin() {
  const setSession = useAuthStore((s) => s.setSession);

  return useMutation({
    mutationFn: async (input: LoginInput) => {
      const res = await api.post<ApiEnvelope<LoginData>>('/auth/login', input);
      return res.data.data;
    },
    onSuccess: (data) => {
      setSession(data.access_token, toUsuarioSession(data.usuario));
    },
  });
}
