import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { formatDate } from '@/lib/formatters';

interface Props {
  isAdmin: boolean;
  onAbrir: () => void;
}

export function CuadreNoAbierto({ isAdmin, onAbrir }: Props) {
  return (
    <Card>
      <CardContent className="flex flex-col items-center gap-4 py-16 text-center">
        <div className="space-y-1">
          <p className="text-sm text-muted-foreground">Caja del día · {formatDate(new Date().toISOString())}</p>
          <p className="text-lg font-medium">La caja aún no ha sido abierta</p>
        </div>
        {isAdmin ? (
          <Button onClick={onAbrir}>Abrir caja</Button>
        ) : (
          <p className="text-sm text-muted-foreground">
            Como vendedora, pide al admin que abra la caja.
          </p>
        )}
      </CardContent>
    </Card>
  );
}
