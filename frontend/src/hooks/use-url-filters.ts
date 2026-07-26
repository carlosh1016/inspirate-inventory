'use client';

import { useCallback, useMemo } from 'react';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';

export interface UrlFiltersConfig<T> {
  defaults: T;
  parsers: {
    [K in keyof T]: (value: string | null) => T[K];
  };
  serializers: {
    [K in keyof T]: (value: T[K]) => string | null;
  };
}

// Reads filters from the URL query string (falling back to defaults) and
// writes them back so filters survive reloads and can be shared by URL. A
// value equal to its default serializes to null and is omitted, keeping the
// URL clean. Any change other than `page` resets `page` to 1 when present.
export function useUrlFilters<T extends Record<string, unknown>>(config: UrlFiltersConfig<T>) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const filters = useMemo(() => {
    const result = {} as T;
    (Object.keys(config.defaults) as (keyof T)[]).forEach((key) => {
      result[key] = config.parsers[key](searchParams.get(key as string));
    });
    return result;
  }, [searchParams, config]);

  const commit = useCallback(
    (next: T) => {
      const params = new URLSearchParams();
      (Object.keys(next) as (keyof T)[]).forEach((key) => {
        const serialized = config.serializers[key](next[key]);
        if (serialized !== null && serialized !== '') {
          params.set(key as string, serialized);
        }
      });
      const qs = params.toString();
      router.replace(qs ? `${pathname}?${qs}` : pathname, { scroll: false });
    },
    [router, pathname, config],
  );

  const setFilter = useCallback(
    <K extends keyof T>(key: K, value: T[K]) => {
      const next = { ...filters, [key]: value };
      if (key !== 'page' && 'page' in next) {
        (next as Record<string, unknown>).page = 1;
      }
      commit(next);
    },
    [filters, commit],
  );

  // Set several keys atomically (a single URL write). Resets page to 1 unless
  // page itself is among the changed keys.
  const setFilters = useCallback(
    (partial: Partial<T>) => {
      const next = { ...filters, ...partial };
      if (!('page' in partial) && 'page' in next) {
        (next as Record<string, unknown>).page = 1;
      }
      commit(next);
    },
    [filters, commit],
  );

  const resetFilters = useCallback(() => commit(config.defaults), [commit, config]);

  return { filters, setFilter, setFilters, resetFilters };
}
