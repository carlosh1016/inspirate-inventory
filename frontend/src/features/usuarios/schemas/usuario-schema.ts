import { z } from 'zod';

const rolSchema = z.enum(['admin', 'vendedora'], { error: 'Selecciona un rol' });

export const createUsuarioSchema = z
  .object({
    nombre_completo: z.string().min(3, 'Mínimo 3 caracteres').max(200, 'Máximo 200 caracteres'),
    correo: z.string().min(1, 'Requerido').email('Correo inválido'),
    password: z.string().min(8, 'La contraseña debe tener al menos 8 caracteres').max(100),
    confirmar_password: z.string().min(1, 'Confirma la contraseña'),
    rol: rolSchema,
  })
  .refine((d) => d.password === d.confirmar_password, {
    message: 'Las contraseñas no coinciden',
    path: ['confirmar_password'],
  });

export type CreateUsuarioInput = z.infer<typeof createUsuarioSchema>;

export const editUsuarioSchema = z
  .object({
    nombre_completo: z.string().min(3, 'Mínimo 3 caracteres').max(200, 'Máximo 200 caracteres'),
    rol: rolSchema,
    password_nueva: z.union([z.literal(''), z.string().min(8, 'La contraseña debe tener al menos 8 caracteres')]),
    confirmar_password: z.string(),
  })
  .refine((d) => d.password_nueva === '' || d.password_nueva === d.confirmar_password, {
    message: 'Las contraseñas no coinciden',
    path: ['confirmar_password'],
  });

export type EditUsuarioInput = z.infer<typeof editUsuarioSchema>;
