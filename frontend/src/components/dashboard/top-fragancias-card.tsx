import { EmptyState } from '@/components/feedback/empty-state';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { TopFragancia } from '@/features/dashboard/api/use-resumen-hoy';
import { formatCurrency, formatGramos } from '@/lib/formatters';

export function TopFraganciasCard({ fragancias }: { fragancias: TopFragancia[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Top fragancias del día</CardTitle>
      </CardHeader>
      <CardContent>
        {fragancias.length === 0 ? (
          <EmptyState title="Aún no hay fragancias vendidas hoy" />
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border text-left text-muted-foreground">
                <th className="py-2 font-medium">Fragancia</th>
                <th className="py-2 text-right font-medium">Gramos</th>
                <th className="py-2 text-right font-medium">Monto</th>
              </tr>
            </thead>
            <tbody>
              {fragancias.map((f) => (
                <tr key={f.id} className="border-b border-border last:border-0">
                  <td className="py-2">{f.nombre_comercial}</td>
                  <td className="py-2 text-right">{formatGramos(f.gramos_vendidos)}</td>
                  <td className="py-2 text-right font-medium">{formatCurrency(f.monto_vendido)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </CardContent>
    </Card>
  );
}
