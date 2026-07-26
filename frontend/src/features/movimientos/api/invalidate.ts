import type { QueryClient } from '@tanstack/react-query';

// A movimiento touches stock and (potentially) every catalog item's snapshot.
export function invalidateInventory(queryClient: QueryClient) {
  ['stock', 'movimientos', 'fragancias', 'productos', 'variantes-envase'].forEach((key) =>
    queryClient.invalidateQueries({ queryKey: [key] }),
  );
}
