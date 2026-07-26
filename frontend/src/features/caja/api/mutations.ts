import { useMutation, useQueryClient, type QueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { Consignacion, Cuadre, PagoCaja, Warning } from '../types';

function invalidateCuadres(queryClient: QueryClient) {
  queryClient.invalidateQueries({ queryKey: ['cuadres-caja'] });
}

export interface AbrirCuadreResult {
  cuadre: Cuadre;
  warnings: Warning[];
}

export function useAbrirCuadre() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (fondoBase: string | null): Promise<AbrirCuadreResult> => {
      const res = await api.post<{ data: Cuadre; warnings?: Warning[] }>(
        '/cuadres-caja/abrir',
        fondoBase ? { fondo_base: fondoBase } : {},
      );
      return { cuadre: res.data.data, warnings: res.data.warnings ?? [] };
    },
    onSuccess: () => invalidateCuadres(queryClient),
  });
}

export interface CerrarCuadreInput {
  valor_turno?: string;
  observaciones?: string | null;
}

export function useCerrarCuadre(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: CerrarCuadreInput) => {
      const res = await api.post<ApiEnvelope<Cuadre>>(`/cuadres-caja/${id}/cerrar`, input);
      return res.data.data;
    },
    onSuccess: () => invalidateCuadres(queryClient),
  });
}

export function useAddPago(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: { concepto: string; monto: string }) => {
      const res = await api.post<ApiEnvelope<PagoCaja>>(`/cuadres-caja/${id}/pagos`, input);
      return res.data.data;
    },
    onSuccess: () => invalidateCuadres(queryClient),
  });
}

export function useDeletePago(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (pagoId: number) => {
      await api.delete(`/cuadres-caja/${id}/pagos/${pagoId}`);
    },
    onSuccess: () => invalidateCuadres(queryClient),
  });
}

export function useAddConsignacion(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (input: { monto: string; banco?: string | null; referencia?: string | null }) => {
      const res = await api.post<ApiEnvelope<Consignacion>>(`/cuadres-caja/${id}/consignaciones`, input);
      return res.data.data;
    },
    onSuccess: () => invalidateCuadres(queryClient),
  });
}

export function useDeleteConsignacion(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (consignacionId: number) => {
      await api.delete(`/cuadres-caja/${id}/consignaciones/${consignacionId}`);
    },
    onSuccess: () => invalidateCuadres(queryClient),
  });
}
