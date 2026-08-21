import { z } from 'zod';

const numeric = z.string().min(1, 'Requerido').regex(/^\d+(\.\d+)?$/, 'Número inválido');

export const modeloEnvaseSchema = z.object({
  tipo: z.string().min(2, 'Mínimo 2 caracteres').max(100, 'Máximo 100 caracteres'),
  tamano_oz: numeric,
  equiv_gramos: numeric,
  precio_solo: numeric,
  precio_con_fragancia: numeric,
  precio_recarga: numeric,
  tiene_variantes: z.boolean(),
});

export type ModeloEnvaseInput = z.infer<typeof modeloEnvaseSchema>;
