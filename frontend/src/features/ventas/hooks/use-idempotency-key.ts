'use client';

import { useState } from 'react';

// Generates a UUID v4 once when the form mounts and keeps it stable across
// renders, so an accidental double-submit replays the same POST /ventas.
export function useIdempotencyKey(): string {
  const [key] = useState(() => crypto.randomUUID());
  return key;
}
