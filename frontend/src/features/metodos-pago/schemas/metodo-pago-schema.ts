import { z } from 'zod';

export const metodoPagoSchema = z.object({
  nombre: z.string().min(2, 'Mínimo 2 caracteres').max(100, 'Máximo 100 caracteres'),
  codigo: z.string().min(1, 'Requerido').max(50, 'Máximo 50 caracteres'),
});

export type MetodoPagoInput = z.infer<typeof metodoPagoSchema>;
