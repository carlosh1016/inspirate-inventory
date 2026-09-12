import { useMutation, useQueryClient, type QueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { Conteo, ConteoItem } from '../types';

function invalidateConteos(queryClient: QueryClient) {
  queryClient.invalidateQueries({ queryKey: ['conteos-inventario'] });
}

// POST /conteos-inventario — sin body, el backend usa el mes actual.
export function useCrearConteo() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (): Promise<Conteo> => {
      const res = await api.post<ApiEnvelope<Conteo>>('/conteos-inventario', {});
      return res.data.data;
    },
    onSuccess: () => invalidateConteos(queryClient),
  });
}

// PATCH /conteos-inventario/:id/items/:itemId
export function useActualizarItemConteo(conteoId: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: { itemId: number; gramosFisico: string }): Promise<ConteoItem> => {
      const res = await api.patch<ApiEnvelope<ConteoItem>>(
        `/conteos-inventario/${conteoId}/items/${input.itemId}`,
        { gramos_fisico: input.gramosFisico },
      );
      return res.data.data;
    },
    onSuccess: () => invalidateConteos(queryClient),
  });
}

// POST /conteos-inventario/:id/cerrar
export function useCerrarConteo(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (): Promise<Conteo> => {
      const res = await api.post<ApiEnvelope<Conteo>>(`/conteos-inventario/${id}/cerrar`);
      return res.data.data;
    },
    onSuccess: () => invalidateConteos(queryClient),
  });
}

// GET /reportes/conteos-inventario/:id — descarga el .xlsx. La ruta va
// detrás de auth (Bearer), así que no puede ser un <a href> plano: se pide
// como blob y se dispara la descarga desde JS.
export function useExportarConteoInventario() {
  return useMutation({
    mutationFn: async ({ id, periodo }: { id: number; periodo: string }) => {
      const res = await api.get(`/reportes/conteos-inventario/${id}`, { responseType: 'blob' });
      const url = URL.createObjectURL(res.data as Blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `conteo-inventario-${periodo}.xlsx`;
      a.click();
      URL.revokeObjectURL(url);
    },
  });
}
