import { AxiosError } from 'axios';

// RFC 7807 Problem Details as returned by the backend (Content-Type
// application/problem+json). `errors` carries per-field validation messages;
// `extra` is an optional structured, error-specific payload.
export interface ProblemDetails {
  type?: string;
  title: string;
  status: number;
  detail?: string;
  instance?: string;
  // `extra` is an arbitrary error-specific payload. Its shape depends on the
  // error: for "stock insuficiente" it's an array of StockInsuficienteItem
  // (see features/movimientos), so it's typed `unknown` and narrowed at use.
  extra?: unknown;
  errors?: Record<string, string[]>;
}

export function parseApiError(error: unknown): ProblemDetails {
  if (error instanceof AxiosError && error.response?.data) {
    const data = error.response.data as Partial<ProblemDetails>;
    return {
      type: data.type,
      title: data.title ?? 'Error',
      status: data.status ?? error.response.status,
      detail: data.detail,
      instance: data.instance,
      extra: data.extra,
      errors: data.errors,
    };
  }

  if (error instanceof Error) {
    return { title: 'Error de conexión', status: 0, detail: error.message };
  }

  return { title: 'Error desconocido', status: 0 };
}

/** Best user-facing message: the Problem `detail`, falling back to `title`. */
export function getErrorMessage(error: unknown): string {
  const problem = parseApiError(error);
  return problem.detail ?? problem.title;
}
