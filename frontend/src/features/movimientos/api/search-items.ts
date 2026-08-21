import type { ComboboxOption } from '@/components/forms/combobox';
import { api } from '@/lib/api';
import { formatGramos } from '@/lib/formatters';
import type { ApiListEnvelope } from '@/types/api';
import type { TipoItem } from '@/types/domain';
import type { Fragancia } from '@/features/fragancias/types';
import { modeloEnvaseLabel } from '@/features/modelos-envase/api/search-modelos-envase';
import type { ModeloEnvase } from '@/features/modelos-envase/types';
import type { Producto } from '@/features/productos/types';
import type { VarianteEnvase } from '@/features/variantes-envase/types';

export interface CatalogItemOption extends ComboboxOption {
  /** Short stock hint shown under the label. */
  detail: string;
}

// Searches the catalog for one item type, returning combobox options with a
// stock hint. Used by movimiento item pickers and the movimientos filter.
// modelosMap (only needed for tipo_item='variante_envase') decorates the
// label with the modelo's name — required so the reader can tell which
// envase a grosor belongs to, and so a "sin variantes" modelo (ej. envase de
// lujo) is picked by its own name instead of the hidden grosor sentinel.
export async function searchCatalogItems(
  tipoItem: TipoItem,
  query: string,
  modelosMap?: Map<number, ModeloEnvase>,
): Promise<CatalogItemOption[]> {
  if (tipoItem === 'fragancia') {
    const res = await api.get<ApiListEnvelope<Fragancia>>('/fragancias', {
      params: { q: query, activo: 'true', page_size: 20 },
    });
    return res.data.data.map((f) => ({
      id: f.id,
      label: f.nombre_comercial,
      detail: `Vitrina ${formatGramos(f.stock.vitrina)} · Bodega ${formatGramos(f.stock.bodega)}`,
    }));
  }

  if (tipoItem === 'variante_envase') {
    const res = await api.get<ApiListEnvelope<VarianteEnvase>>('/variantes-envase', {
      params: { q: query, activo: 'true', page_size: 20 },
    });
    return res.data.data.map((v) => {
      const modelo = modelosMap?.get(v.modelo_envase_id);
      const modeloLabel = modelo ? modeloEnvaseLabel(modelo) : `Modelo #${v.modelo_envase_id}`;
      const label = modelo && !modelo.tiene_variantes ? modeloLabel : `${modeloLabel} · ${v.color}`;
      return {
        id: v.id,
        label,
        detail: `Vitrina ${v.stock.vitrina} · Bodega ${v.stock.bodega}`,
      };
    });
  }

  const res = await api.get<ApiListEnvelope<Producto>>('/productos', {
    params: { q: query, activo: 'true', page_size: 20 },
  });
  return res.data.data.map((p) => ({
    id: p.id,
    label: p.nombre,
    detail: `Vitrina ${p.stock.vitrina} · Bodega ${p.stock.bodega}`,
  }));
}
