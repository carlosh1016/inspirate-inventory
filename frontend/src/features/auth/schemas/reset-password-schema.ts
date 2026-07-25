import { z } from 'zod';

export const resetPasswordSchema = z
  .object({
    password_nueva: z.string().min(8, 'La contraseña debe tener al menos 8 caracteres'),
    confirmar_password: z.string().min(1, 'Confirma la contraseña'),
  })
  .refine((d) => d.password_nueva === d.confirmar_password, {
    message: 'Las contraseñas no coinciden',
    path: ['confirmar_password'],
  });

export type ResetPasswordInput = z.infer<typeof resetPasswordSchema>;
