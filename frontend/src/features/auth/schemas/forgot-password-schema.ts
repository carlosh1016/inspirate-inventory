import { z } from 'zod';

export const forgotPasswordSchema = z.object({
  correo: z.email('Correo inválido'),
});

export type ForgotPasswordInput = z.infer<typeof forgotPasswordSchema>;
