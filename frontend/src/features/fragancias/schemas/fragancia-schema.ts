import { z } from 'zod';

export const fraganciaSchema = z.object({
  nombre_comercial: z.string().min(2, 'Mínimo 2 caracteres').max(200, 'Máximo 200 caracteres'),
  nombre_alternativo: z
    .string()
    .max(200, 'Máximo 200 caracteres')
    .optional()
    .or(z.literal('')),
  genero: z.enum(['masculina', 'femenina']),
  numero_genero: z.number({ error: 'Requerido' }).int('Debe ser un número entero').min(1, 'Debe ser mayor a 0'),
  gramos_minimo: z
    .string()
    .min(1, 'Requerido')
    .regex(/^\d+(\.\d+)?$/, 'Número inválido'),
});

export type FraganciaInput = z.infer<typeof fraganciaSchema>;
