import { z } from 'zod';

const gramos = z.string().regex(/^\d+(\.\d+)?$/, 'Gramos inválidos');
const cantidad = z.number().int('Cantidad inválida').positive('Cantidad mínima 1');

const ventaItemSchema = z.discriminatedUnion('tipo_linea', [
  z.object({
    tipo_linea: z.literal('envase_con_fragancia'),
    variante_envase_id: z.number().positive('Selecciona un envase'),
    fragancia_id: z.number().positive('Selecciona una fragancia'),
    gramos_fragancia: gramos,
    cantidad,
    feromona_producto_id: z.number().positive('Selecciona la feromona').optional(),
  }),
  z.object({
    tipo_linea: z.literal('recarga'),
    variante_envase_id: z.number().positive('Selecciona un envase'),
    fragancia_id: z.number().positive('Selecciona una fragancia'),
    gramos_fragancia: gramos,
    cantidad,
    feromona_producto_id: z.number().positive('Selecciona la feromona').optional(),
  }),
  z.object({
    tipo_linea: z.literal('envase_solo'),
    variante_envase_id: z.number().positive('Selecciona un envase'),
    cantidad,
  }),
  z.object({
    tipo_linea: z.literal('producto_otro'),
    producto_id: z.number().positive('Selecciona un producto'),
    cantidad,
  }),
]);

export const nuevaVentaSchema = z.object({
  metodo_pago_id: z.number().positive('Selecciona un método de pago'),
  observaciones: z.string().max(1000).optional().nullable(),
  items: z.array(ventaItemSchema).min(1, 'Agrega al menos un ítem'),
});

export type NuevaVentaInput = z.infer<typeof nuevaVentaSchema>;
