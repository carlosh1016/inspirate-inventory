'use client';

import Link from 'next/link';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';

import { FormField } from '@/components/forms/form-field';
import { Button, buttonVariants } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { getErrorMessage } from '@/lib/errors';
import { cn } from '@/lib/utils';
import { useRequestPasswordReset } from '../api/use-request-password-reset';
import { forgotPasswordSchema, type ForgotPasswordInput } from '../schemas/forgot-password-schema';

export function ForgotPasswordForm() {
  const request = useRequestPasswordReset();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting, isSubmitSuccessful },
  } = useForm<ForgotPasswordInput>({ resolver: zodResolver(forgotPasswordSchema) });

  const onSubmit = async (data: ForgotPasswordInput) => {
    try {
      await request.mutateAsync(data);
      toast.success('Si el correo existe, te enviamos instrucciones para restablecer tu contraseña.');
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Restablecer contraseña</CardTitle>
        <CardDescription>Te enviaremos un enlace a tu correo para crear una nueva contraseña.</CardDescription>
      </CardHeader>
      <CardContent>
        {isSubmitSuccessful ? (
          <div className="space-y-4">
            <p className="text-sm text-muted-foreground">
              Revisa tu correo. Si la dirección está registrada, recibirás un enlace para continuar.
            </p>
            <Link href="/login" className={cn(buttonVariants(), 'w-full')}>
              Volver a iniciar sesión
            </Link>
          </div>
        ) : (
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

            <Button type="submit" className="w-full" disabled={isSubmitting}>
              {isSubmitting ? 'Enviando...' : 'Enviar instrucciones'}
            </Button>

            <p className="text-center text-sm text-muted-foreground">
              <Link href="/login" className="hover:text-foreground">
                Volver a iniciar sesión
              </Link>
            </p>
          </form>
        )}
      </CardContent>
    </Card>
  );
}
