import type { ComboboxOption } from '@/components/forms/combobox';
import { api } from '@/lib/api';
import { formatGramos } from '@/lib/formatters';
import type { ApiListEnvelope } from '@/types/api';
import type { Fragancia } from '@/features/fragancias/types';
import type { ModeloEnvase } from '@/features/modelos-envase/types';
import { modeloEnvaseLabel } from '@/features/modelos-envase/api/search-modelos-envase';
import type { Producto } from '@/features/productos/types';
import type { VarianteEnvase } from '@/features/variantes-envase/types';

export interface EnvaseOption extends ComboboxOption {
  equiv_gramos: string;
  precio_solo: string;
  precio_con_fragancia: string;
  precio_recarga: string;
  detail: string;
}

export interface FraganciaOption extends ComboboxOption {
  detail: string;
}

export interface ProductoOption extends ComboboxOption {
  precio: string;
  detail: string;
}

// Active variantes joined with their modelo (for label + prices + equiv_gramos).
// modelosMap comes from useModelosFullLookup so it's fetched once per form.
export async function searchEnvasesVenta(
  query: string,
  modelosMap: Map<number, ModeloEnvase>,
): Promise<EnvaseOption[]> {
  const res = await api.get<ApiListEnvelope<VarianteEnvase>>('/variantes-envase', {
    params: { q: query, activo: 'true', page_size: 20 },
  });
  return res.data.data.map((v) => {
    const modelo = modelosMap.get(v.modelo_envase_id);
    const modeloLabel = modelo ? modeloEnvaseLabel(modelo) : `Modelo #${v.modelo_envase_id}`;
    // Un modelo "sin variantes" (ej. envase de lujo) tiene una única variante
    // oculta — no se le muestra el sufijo de grosor, se elige por su nombre.
    const label = modelo && !modelo.tiene_variantes ? modeloLabel : `${modeloLabel} · ${v.color}`;
    return {
      id: v.id,
      label,
      equiv_gramos: modelo?.equiv_gramos ?? '0',
      precio_solo: modelo?.precio_solo ?? '0',
      precio_con_fragancia: modelo?.precio_con_fragancia ?? '0',
      precio_recarga: modelo?.precio_recarga ?? '0',
      detail: `Stock ${v.stock.total}`,
    };
  });
}

export async function searchFraganciasVenta(query: string): Promise<FraganciaOption[]> {
  const res = await api.get<ApiListEnvelope<Fragancia>>('/fragancias', {
    params: { q: query, activo: 'true', page_size: 20 },
  });
  return res.data.data.map((f) => ({
    id: f.id,
    label: f.nombre_comercial,
    detail: `Vitrina ${formatGramos(f.stock.vitrina)} · Bodega ${formatGramos(f.stock.bodega)}`,
  }));
}

export async function searchProductosVenta(query: string): Promise<ProductoOption[]> {
  const res = await api.get<ApiListEnvelope<Producto>>('/productos', {
    params: { q: query, activo: 'true', page_size: 20 },
  });
  return res.data.data.map((p) => ({
    id: p.id,
    label: p.nombre,
    precio: p.precio,
    detail: `Stock ${p.stock.total}`,
  }));
}

export async function searchFeromonas(query: string): Promise<ProductoOption[]> {
  const res = await api.get<ApiListEnvelope<Producto>>('/productos', {
    params: { q: query, categoria: 'feromona', activo: 'true', page_size: 20 },
  });
  return res.data.data.map((p) => ({
    id: p.id,
    label: p.nombre,
    precio: p.precio,
    detail: '',
  }));
}
