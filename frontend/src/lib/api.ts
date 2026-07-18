import axios from 'axios';

export const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  withCredentials: true,
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // TODO(module-auth): redirigir a login / limpiar sesión cuando exista el flujo de auth.
    }
    return Promise.reject(error);
  },
);
