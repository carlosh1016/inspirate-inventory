#!/bin/bash
# Script para ejecutar en LOCAL (no en el servidor)
# Compila el binario y lo sube al servidor Oracle Cloud
# Uso: ./backend/deploy/deploy.sh <IP_DEL_SERVIDOR>

set -e

SERVER_IP=$1
SSH_USER="ubuntu"
SSH_KEY="~/.ssh/id_rsa"  # Ajustar si tu key tiene otro nombre
DEPLOY_DIR="/opt/inspirate"

if [ -z "$SERVER_IP" ]; then
  echo "Uso: $0 <IP_DEL_SERVIDOR>"
  exit 1
fi

echo "1. Compilando binario para Linux amd64..."
cd backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/server ./cmd/api
cd ..

echo "2. Subiendo binario al servidor..."
ssh -i $SSH_KEY $SSH_USER@$SERVER_IP "sudo mkdir -p $DEPLOY_DIR && sudo chown ubuntu:ubuntu $DEPLOY_DIR"
scp -i $SSH_KEY backend/bin/server $SSH_USER@$SERVER_IP:$DEPLOY_DIR/server
ssh -i $SSH_KEY $SSH_USER@$SERVER_IP "chmod +x $DEPLOY_DIR/server"

echo "3. Reiniciando servicio..."
ssh -i $SSH_KEY $SSH_USER@$SERVER_IP "sudo systemctl restart inspirate"

echo "4. Verificando estado..."
sleep 3
ssh -i $SSH_KEY $SSH_USER@$SERVER_IP "sudo systemctl status inspirate --no-pager"

echo "Deploy completado. Verificar: curl http://$SERVER_IP/api/v1/health"
