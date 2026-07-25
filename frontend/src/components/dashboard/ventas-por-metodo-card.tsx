import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { PorMetodoPago } from '@/features/dashboard/api/use-resumen-hoy';
import { formatCurrency } from '@/lib/formatters';

interface Props {
  totales: PorMetodoPago;
  totalDia: string;
}

export function VentasPorMetodoCard({ totales, totalDia }: Props) {
  const rows: Array<[string, string]> = [
    ['Efectivo', totales.efectivo],
    ['Nequi', totales.nequi],
    ['Daviplata', totales.daviplata],
    ['Transferencia', totales.transferencia],
    ['Otros', totales.otros],
  ];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Ventas por método de pago</CardTitle>
      </CardHeader>
      <CardContent>
        <table className="w-full text-sm">
          <tbody>
            {rows.map(([label, value]) => (
              <tr key={label} className="border-b border-border">
                <td className="py-2 text-muted-foreground">{label}</td>
                <td className="py-2 text-right font-medium">{formatCurrency(value)}</td>
              </tr>
            ))}
            <tr className="font-semibold">
              <td className="py-2">Total</td>
              <td className="py-2 text-right">{formatCurrency(totalDia)}</td>
            </tr>
          </tbody>
        </table>
      </CardContent>
    </Card>
  );
}
