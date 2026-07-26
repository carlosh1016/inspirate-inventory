import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { CategoriaProducto } from '@/types/domain';
import type { Producto } from '../types';

export interface UpdateProductoPayload {
  nombre?: string;
  categoria?: CategoriaProducto;
  precio?: string;
  stock_minimo?: number;
}

export function useUpdateProducto(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: UpdateProductoPayload) => {
      const res = await api.patch<ApiEnvelope<Producto>>(`/productos/${id}`, input);
      return res.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['productos'] });
      queryClient.invalidateQueries({ queryKey: ['stock'] });
    },
  });
}
