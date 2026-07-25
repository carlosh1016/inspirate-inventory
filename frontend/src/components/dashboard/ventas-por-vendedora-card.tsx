import { EmptyState } from '@/components/feedback/empty-state';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { VentaPorVendedora } from '@/features/dashboard/api/use-resumen-hoy';
import { formatCurrency } from '@/lib/formatters';

export function VentasPorVendedoraCard({ vendedoras }: { vendedoras: VentaPorVendedora[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Ventas por vendedora</CardTitle>
      </CardHeader>
      <CardContent>
        {vendedoras.length === 0 ? (
          <EmptyState title="Sin ventas registradas hoy" />
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border text-left text-muted-foreground">
                <th className="py-2 font-medium">Vendedora</th>
                <th className="py-2 text-right font-medium">Ventas</th>
                <th className="py-2 text-right font-medium">Total</th>
              </tr>
            </thead>
            <tbody>
              {vendedoras.map((v) => (
                <tr key={v.usuario_id} className="border-b border-border last:border-0">
                  <td className="py-2">{v.nombre_completo}</td>
                  <td className="py-2 text-right">{v.ventas_count}</td>
                  <td className="py-2 text-right font-medium">{formatCurrency(v.total)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </CardContent>
    </Card>
  );
}
