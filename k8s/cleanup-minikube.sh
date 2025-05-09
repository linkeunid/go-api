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

echo -e "${BOLD}${RED}🧹 Cleaning up Go API Kubernetes resources from Minikube...${NC}"

# Get the script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Define kubectl command (use minikube kubectl)
KUBECTL="minikube kubectl --"

# Delete Kubernetes resources
echo -e "${CYAN}🗑️ Deleting Kubernetes manifests...${NC}"
$KUBECTL delete -f deployment.yaml --ignore-not-found=true
$KUBECTL delete -f redis.yaml --ignore-not-found=true
$KUBECTL delete -f mysql.yaml --ignore-not-found=true
$KUBECTL delete -f secrets.yaml --ignore-not-found=true
$KUBECTL delete -f configmap.yaml --ignore-not-found=true

echo -e "${GREEN}✨ Cleanup completed!${NC}"
echo -e "${YELLOW}💡 If you want to stop Minikube completely, run:${NC}"
echo -e "${BOLD}${YELLOW}   minikube stop${NC}" 