package auditoria

// Audit action names, centralized here for reference and organized by domain.
// Only actions actually recorded by the codebase are listed (verified against
// the usecases that call auditoria.Insert); completeness-only constants would
// be noise.
const (
	// Auth
	AccionLoginSuccess           = "login_success"
	AccionLoginFailed            = "login_failed"
	AccionLogout                 = "logout"
	AccionPasswordResetRequested = "password_reset_requested"
	AccionPasswordResetCompleted = "password_reset_completed"
	AccionPasswordChanged        = "password_changed"

	// Usuarios
	AccionUsuarioCreado      = "usuario_creado"
	AccionUsuarioEditado     = "usuario_editado"
	AccionUsuarioActivado    = "usuario_activado"
	AccionUsuarioDesactivado = "usuario_desactivado"
	AccionUsuarioEliminado   = "usuario_eliminado"

	// Catálogo
	AccionFraganciaCreada       = "fragancia_creada"
	AccionFraganciaEditada      = "fragancia_editada"
	AccionFraganciaEliminada    = "fragancia_eliminada"
	AccionFraganciaRestaurada   = "fragancia_restaurada"
	AccionModeloEnvaseCreado    = "modelo_envase_creado"
	AccionModeloEnvaseEditado   = "modelo_envase_editado"
	AccionModeloEnvaseEliminado = "modelo_envase_eliminado"
	AccionVarianteEnvaseCreada  = "variante_envase_creada"
	AccionVarianteEnvaseEditada = "variante_envase_editada"
	AccionVarianteEnvaseElim    = "variante_envase_eliminada"
	AccionProductoCreado        = "producto_creado"
	AccionProductoEditado       = "producto_editado"
	AccionProductoEliminado     = "producto_eliminado"
	AccionMetodoPagoCreado      = "metodo_pago_creado"
	AccionMetodoPagoEditado     = "metodo_pago_editado"
	AccionMetodoPagoEliminado   = "metodo_pago_eliminado"

	// Inventario
	AccionAjusteInventario     = "ajuste_inventario"
	AccionCorreccionInventario = "correccion_inventario"

	// Ventas
	AccionVentaObservacionesEditadas = "venta_observaciones_editadas"

	// Cuadre
	AccionCuadreAbierto = "cuadre_abierto"
	AccionCuadreCerrado = "cuadre_cerrado"

	// Sesiones
	AccionSesionEditada = "sesion_editada"
)
