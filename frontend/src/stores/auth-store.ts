import { create } from 'zustand';

import type { UsuarioSession } from '@/types/domain';

interface AuthState {
  accessToken: string | null;
  usuario: UsuarioSession | null;
  /** True until the initial refresh-on-mount resolves (see AuthProvider). */
  isInitializing: boolean;
  setAccessToken: (token: string) => void;
  setSession: (token: string, usuario: UsuarioSession) => void;
  setUsuario: (usuario: UsuarioSession) => void;
  clearSession: () => void;
  setInitializing: (v: boolean) => void;
}

// Access token lives only in memory (never localStorage). It is lost on reload
// and restored by the AuthProvider via POST /auth/refresh.
export const useAuthStore = create<AuthState>((set) => ({
  accessToken: null,
  usuario: null,
  isInitializing: true,
  setAccessToken: (token) => set({ accessToken: token }),
  setSession: (token, usuario) => set({ accessToken: token, usuario }),
  setUsuario: (usuario) => set({ usuario }),
  clearSession: () => set({ accessToken: null, usuario: null }),
  setInitializing: (v) => set({ isInitializing: v }),
}));
