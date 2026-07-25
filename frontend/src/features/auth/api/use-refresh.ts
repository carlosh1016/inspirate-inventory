import { useMutation } from '@tanstack/react-query';

import { api } from '@/lib/api';
import { useAuthStore } from '@/stores/auth-store';
import type { ApiEnvelope } from '@/types/api';

interface RefreshData {
  access_token: string;
  expires_in: number;
}

export function useRefresh() {
  const setAccessToken = useAuthStore((s) => s.setAccessToken);

  return useMutation({
    mutationFn: async () => {
      const res = await api.post<ApiEnvelope<RefreshData>>('/auth/refresh');
      return res.data.data;
    },
    onSuccess: (data) => setAccessToken(data.access_token),
  });
}
