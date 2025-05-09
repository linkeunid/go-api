#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
RED='\033[0;31m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
NC='\033[0m' # No Color

echo -e "${BOLD}${MAGENTA}🚀 Starting Minikube deployment process for Go API...${NC}"

# Check if minikube is running
if ! minikube status | grep -q "Running"; then
  echo -e "${YELLOW}🔄 Minikube is not running. Starting minikube...${NC}"
  minikube start
fi

# Enable ingress addon if not already enabled
if ! minikube addons list | grep "ingress" | grep -q "enabled"; then
  echo -e "${YELLOW}🔌 Enabling ingress addon...${NC}"
  minikube addons enable ingress
fi

# Get the script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Define kubectl command (use minikube kubectl)
KUBECTL="minikube kubectl --"

# Build the Docker image
echo -e "${BLUE}🔨 Building Docker image...${NC}"
docker build -t go-api:latest "$PROJECT_ROOT"

# Load the image into Minikube
echo -e "${BLUE}📦 Loading image into Minikube...${NC}"
minikube image load go-api:latest

# Apply Kubernetes manifests
echo -e "${CYAN}📄 Applying Kubernetes manifests...${NC}"
cd "$SCRIPT_DIR"
$KUBECTL apply -f configmap.yaml
$KUBECTL apply -f secrets.yaml
$KUBECTL apply -f mysql.yaml
$KUBECTL apply -f redis.yaml
$KUBECTL apply -f deployment.yaml

# Wait for pods to be ready
echo -e "${CYAN}⏳ Waiting for pods to be ready...${NC}"
$KUBECTL wait --for=condition=ready pod -l app=go-api --timeout=180s || true
$KUBECTL wait --for=condition=ready pod -l app=mysql --timeout=180s || true
$KUBECTL wait --for=condition=ready pod -l app=redis --timeout=180s || true

# Get Minikube IP
MINIKUBE_IP=$(minikube ip)
echo -e "${GREEN}🌐 Minikube IP: ${MINIKUBE_IP}${NC}"

# Add entry to /etc/hosts if not already present
if ! grep -q "go-api.local" /etc/hosts; then
  echo -e "${YELLOW}📝 Please add the following entry to your /etc/hosts file:${NC}"
  echo -e "${BOLD}${YELLOW}   ${MINIKUBE_IP} go-api.local${NC}"
  echo -e "${YELLOW}💻 Run:${NC}"
  echo -e "${BOLD}${YELLOW}   echo '${MINIKUBE_IP} go-api.local' | sudo tee -a /etc/hosts${NC}"
else
  echo -e "${GREEN}✅ Host entry already exists.${NC}"
fi

echo -e "${GREEN}🎉 Deployment completed!${NC}"
echo -e "${GREEN}🔗 You can access the application at: ${BOLD}http://go-api.local${NC}"
echo -e "${CYAN}📊 To check pod status: ${BOLD}${KUBECTL} get pods${NC}"
echo -e "${CYAN}📋 To view logs: ${BOLD}${KUBECTL} logs deployment/go-api${NC}" 