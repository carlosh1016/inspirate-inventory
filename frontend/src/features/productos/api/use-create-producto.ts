import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { CategoriaProducto } from '@/types/domain';
import type { Producto } from '../types';

export interface CreateProductoPayload {
  nombre: string;
  categoria: CategoriaProducto;
  precio: string;
  stock_minimo: number;
}

export function useCreateProducto() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: CreateProductoPayload) => {
      const res = await api.post<ApiEnvelope<Producto>>('/productos', input);
      return res.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['productos'] });
      queryClient.invalidateQueries({ queryKey: ['stock'] });
    },
  });
}
