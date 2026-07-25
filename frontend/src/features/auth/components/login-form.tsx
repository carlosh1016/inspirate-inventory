'use client';

import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';

import { FormField } from '@/components/forms/form-field';
import { PasswordInput } from '@/components/forms/password-input';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { getErrorMessage } from '@/lib/errors';
import { useLogin } from '../api/use-login';
import { loginSchema, type LoginInput } from '../schemas/login-schema';

export function LoginForm() {
  const router = useRouter();
  const params = useSearchParams();
  const login = useLogin();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginInput>({ resolver: zodResolver(loginSchema) });

  const onSubmit = async (data: LoginInput) => {
    try {
      const result = await login.mutateAsync(data);
      toast.success(`Bienvenida, ${result.usuario.nombre_completo}`);
      const from = params.get('from');
      router.push(from ?? '/dashboard');
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Inicia sesión</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4" noValidate>
          <FormField id="correo" label="Correo" error={errors.correo?.message}>
            <Input
              id="correo"
              type="email"
              autoComplete="email"
              disabled={isSubmitting}
              {...register('correo')}
            />
          </FormField>

          <FormField id="password" label="Contraseña" error={errors.password?.message}>
            <PasswordInput
              id="password"
              autoComplete="current-password"
              disabled={isSubmitting}
              {...register('password')}
            />
          </FormField>

          <Button type="submit" className="w-full" disabled={isSubmitting}>
            {isSubmitting ? 'Iniciando...' : 'Iniciar sesión'}
          </Button>

          <p className="text-center text-sm text-muted-foreground">
            <Link href="/forgot-password" className="hover:text-foreground">
              ¿Olvidaste tu contraseña?
            </Link>
          </p>
        </form>
      </CardContent>
    </Card>
  );
}
