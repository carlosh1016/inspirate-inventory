'use client';

import { useEffect, useState } from 'react';

export interface Elapsed {
  hours: number;
  minutes: number;
  seconds: number;
  totalMs: number;
}

// Ticks every second while `startIso` is set. setState happens inside the
// interval callback (not synchronously in the effect), so it is lint-safe.
export function useElapsedTime(startIso: string | null | undefined): Elapsed | null {
  const [now, setNow] = useState(() => Date.now());

  useEffect(() => {
    if (!startIso) return;
    const id = setInterval(() => setNow(Date.now()), 1000);
    return () => clearInterval(id);
  }, [startIso]);

  if (!startIso) return null;

  const totalMs = Math.max(0, now - new Date(startIso).getTime());
  const totalSec = Math.floor(totalMs / 1000);
  return {
    hours: Math.floor(totalSec / 3600),
    minutes: Math.floor((totalSec % 3600) / 60),
    seconds: totalSec % 60,
    totalMs,
  };
}
