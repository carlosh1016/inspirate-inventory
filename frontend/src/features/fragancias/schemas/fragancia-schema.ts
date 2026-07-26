import { z } from 'zod';

export const fraganciaSchema = z.object({
  nombre_comercial: z.string().min(2, 'Mínimo 2 caracteres').max(200, 'Máximo 200 caracteres'),
  nombre_alternativo: z
    .string()
    .max(200, 'Máximo 200 caracteres')
    .optional()
    .or(z.literal('')),
  genero: z.enum(['masculina', 'femenina']),
  gramos_minimo: z
    .string()
    .min(1, 'Requerido')
    .regex(/^\d+(\.\d+)?$/, 'Número inválido'),
});

export type FraganciaInput = z.infer<typeof fraganciaSchema>;
