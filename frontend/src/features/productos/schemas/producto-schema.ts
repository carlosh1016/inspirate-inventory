import { z } from 'zod';

export const productoSchema = z.object({
  nombre: z.string().min(2, 'Mínimo 2 caracteres').max(200, 'Máximo 200 caracteres'),
  categoria: z.enum(['crema', 'feromona', 'hogar', 'regalo', 'otro']),
  precio: z.string().min(1, 'Requerido').regex(/^\d+(\.\d+)?$/, 'Número inválido'),
  stock_minimo: z.number({ error: 'Requerido' }).int('Debe ser entero').min(0, 'No puede ser negativo'),
});

export type ProductoInput = z.infer<typeof productoSchema>;
