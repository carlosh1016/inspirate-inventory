# Deploy de Inspirate Inventory

## Stack
- Oracle Cloud VM (Ubuntu 22.04 Always Free) — backend Go
- Supabase — PostgreSQL
- Vercel — frontend Next.js

---

## PARTE 1 — Preparar la VM en Oracle Cloud

### 1.1 Crear la instancia
1. Login en cloud.oracle.com
2. Compute → Instances → Create Instance
3. Configuracion:
   - Name: inspirate-backend
   - Image: Canonical Ubuntu 22.04 (Always Free eligible)
   - Shape: VM.Standard.E2.1.Micro (Always Free) o Ampere A1 (ARM, tambien Always Free)
   - Network: dejar defaults (VCN nueva o existente)
   - SSH keys: subir tu clave publica (~/.ssh/id_rsa.pub) o generar una nueva
4. Click Create

Esperar 2-3 minutos hasta que el estado sea "Running".
Anotar la IP publica de la instancia.

> Si eliges la shape Ampere A1 (ARM), el binario debe compilarse para `arm64`, no `amd64`
> (ver nota en la Parte 3). `build.sh` y `deploy.sh` compilan para `amd64` por default.

### 1.2 Abrir puertos en Oracle Cloud
Oracle Cloud bloquea puertos por defecto en DOS lugares. Hay que abrirlos en ambos:

**A) Security List (firewall de Oracle):**
1. En la instancia → click en la subred → Security Lists → Default Security List
2. Add Ingress Rules:
   - Source: 0.0.0.0/0, Protocol: TCP, Port: 80
   - Source: 0.0.0.0/0, Protocol: TCP, Port: 443

**B) iptables dentro de la VM (Ubuntu bloquea por defecto en Oracle):**
```bash
sudo iptables -I INPUT -p tcp --dport 80 -j ACCEPT
sudo iptables -I INPUT -p tcp --dport 443 -j ACCEPT
sudo netfilter-persistent save
```

### 1.3 Conectarse a la VM
```bash
ssh ubuntu@<IP_PUBLICA>
```

### 1.4 Instalar dependencias en la VM
```bash
# Actualizar paquetes
sudo apt update && sudo apt upgrade -y

# Instalar herramientas basicas
sudo apt install -y curl wget git nano iptables-persistent

# Instalar Caddy (reverse proxy con HTTPS automatico)
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install caddy -y

# Instalar goose para migraciones (en la VM)
curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh
```

### 1.5 Crear directorio de la app
```bash
sudo mkdir -p /opt/inspirate
sudo chown ubuntu:ubuntu /opt/inspirate
```

---

## PARTE 2 — Configurar variables de entorno en produccion

Crear el archivo de entorno en la VM:
```bash
nano /opt/inspirate/.env
```

Contenido (reemplazar los valores reales — ver `backend/.env.example` para la lista completa
de variables que lee `internal/platform/config/config.go`; `DATABASE_URL` y `JWT_SECRET` son
obligatorias, el resto tiene defaults razonables):
```bash
DATABASE_URL=postgresql://postgres.<ref>:<password>@aws-0-us-east-1.pooler.supabase.com:6543/postgres
JWT_SECRET=<string-aleatorio-min-32-chars-genera-con: openssl rand -base64 32>
PORT=8080
ENVIRONMENT=production
LOG_LEVEL=info
FRONTEND_URL=https://<tu-app>.vercel.app
CORS_ALLOWED_ORIGINS=https://<tu-app>.vercel.app
```

**Sobre DATABASE_URL de Supabase:**
- Ve a Supabase → Project Settings → Database → Connection string
- Usar **Transaction mode** (puerto 6543) para el servidor en produccion
- Usar **Session mode** (puerto 5432) solo para correr migraciones con goose

Proteger el archivo:
```bash
chmod 600 /opt/inspirate/.env
```

---

## PARTE 3 — Primer deploy del backend

### 3.1 Compilar y subir el binario (desde tu maquina local)
```bash
# Desde la raiz del repo
./backend/deploy/deploy.sh <IP_PUBLICA_ORACLE>
```

> Este script compila para `GOOS=linux GOARCH=amd64`. Si tu instancia Oracle es
> Ampere A1 (ARM), cambia `GOARCH=amd64` por `GOARCH=arm64` en `backend/build.sh`
> y `backend/deploy/deploy.sh` antes de correrlo.

Si es el primer deploy y el servicio no existe todavia, el script fallara en el restart.
En ese caso, continua con 3.2 y 3.3 primero.

### 3.2 Instalar el servicio systemd (solo la primera vez, en la VM)
```bash
# Desde tu maquina local:
scp backend/deploy/inspirate.service ubuntu@<IP>:/tmp/

# Luego en la VM:
sudo mv /tmp/inspirate.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable inspirate
sudo systemctl start inspirate
```

### 3.3 Aplicar migraciones (desde tu maquina local, una sola vez)
```bash
# Usar Session mode de Supabase (puerto 5432) para goose
DATABASE_URL="postgresql://postgres.<ref>:<password>@aws-0-us-east-1.pooler.supabase.com:5432/postgres" \
  goose -dir backend/db/migrations postgres up
```

### 3.4 Verificar que el backend corre
```bash
# En la VM
sudo systemctl status inspirate
sudo journalctl -u inspirate -f

# Desde local (antes de Caddy, directo al puerto de la app)
curl http://<IP_PUBLICA>:8080/api/v1/health
```

> El health check vive en `/api/v1/health` (montado dentro del grupo de rutas
> `/api/v1` en `internal/http/router.go`), no en `/health` sin prefijo. No requiere
> autenticacion. Responde `{"status":"ok","version":"...","checks":{"db":"ok"}}`
> con `200`, o `503` con `"status":"degraded"` si la base de datos no responde.

---

## PARTE 4 — Configurar Caddy (HTTPS)

### 4.1 Sin dominio propio (IP directa, solo HTTP)
```bash
sudo nano /etc/caddy/Caddyfile
```
Contenido:
```
:80 {
    reverse_proxy localhost:8080
}
```

### 4.2 Con dominio propio (HTTPS automatico)
Apunta un registro A de tu dominio a la IP publica de Oracle.
Espera 5-10 minutos para propagacion DNS. Luego:
```bash
sudo nano /etc/caddy/Caddyfile
```
Contenido:
```
api.tudominio.com {
    reverse_proxy localhost:8080
}
```

En ambos casos:
```bash
sudo systemctl restart caddy
sudo systemctl status caddy
```

Caddy obtiene el certificado SSL automaticamente si usas dominio.

Con Caddy al frente, el health check queda en `http://<IP>/api/v1/health`
(o `https://api.tudominio.com/api/v1/health` con dominio) — Caddy hace
reverse proxy de todo el trafico a `localhost:8080`, sin reescribir la ruta.

---

## PARTE 5 — Deploy del frontend en Vercel

### 5.1 Conectar el repo
1. vercel.com → New Project
2. Importar `github.com/carlosh1016/inspirate-inventory`
3. Root Directory: `frontend`
4. Framework: Next.js (se detecta automatico)

### 5.2 Variable de entorno
Antes del primer deploy, agregar:
- Nombre: `NEXT_PUBLIC_API_URL`
- Valor: `http://<IP_PUBLICA_ORACLE>/api/v1` (sin dominio) o `https://api.tudominio.com/api/v1` (con dominio)

> El cliente axios (`frontend/src/lib/api.ts`) usa este valor como `baseURL` y le
> concatena rutas como `/ventas`, `/usuarios`, etc. — **debe incluir el sufijo
> `/api/v1`**, igual que `frontend/.env.example`. Sin ese sufijo todas las
> peticiones del frontend fallaran con 404.

### 5.3 Deploy
Click en Deploy. La URL sera algo como `https://inspirate-inventory.vercel.app`.

---

## PARTE 6 — Actualizar CORS con la URL de Vercel

Editar `/opt/inspirate/.env` en la VM:
```bash
nano /opt/inspirate/.env
# Actualizar FRONTEND_URL y CORS_ALLOWED_ORIGINS con la URL de Vercel
sudo systemctl restart inspirate
```

---

## Deploys futuros (actualizaciones de codigo)

Cada vez que haya cambios en el backend:
```bash
./backend/deploy/deploy.sh <IP_PUBLICA_ORACLE>
```

El frontend en Vercel se despliega automaticamente con cada push a main.

---

## Troubleshooting

```bash
# Ver logs del backend en tiempo real
ssh ubuntu@<IP> "sudo journalctl -u inspirate -f"

# Reiniciar backend
ssh ubuntu@<IP> "sudo systemctl restart inspirate"

# Ver logs de Caddy
ssh ubuntu@<IP> "sudo journalctl -u caddy -f"

# Verificar puertos abiertos
ssh ubuntu@<IP> "sudo ss -tlnp"

# Verificar health desde la VM misma
ssh ubuntu@<IP> "curl localhost:8080/api/v1/health"
```
