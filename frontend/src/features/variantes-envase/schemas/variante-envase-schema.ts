import { z } from 'zod';

// modelo_envase_id is required at create (>0) but not editable afterwards, so
// the edit form omits it. stock_minimo is an integer >= 0.
export const varianteEnvaseSchema = z.object({
  modelo_envase_id: z.number().int().positive('Selecciona un modelo'),
  color: z.string().min(1, 'Requerido').max(100, 'Máximo 100 caracteres'),
  stock_minimo: z
    .number({ error: 'Requerido' })
    .int('Debe ser entero')
    .min(0, 'No puede ser negativo'),
});

export type VarianteEnvaseInput = z.infer<typeof varianteEnvaseSchema>;
