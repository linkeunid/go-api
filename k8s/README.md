# Kubernetes Deployment for Go API

This directory contains Kubernetes manifests for deploying the Go API application to Minikube or a production Kubernetes cluster.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Setup Minikube](#setup-minikube)
- [Deployment](#deployment)
  - [Automated Deployment](#automated-deployment)
  - [Manual Deployment](#manual-deployment)
- [Access the Application](#access-the-application)
- [Deployment Architecture](#deployment-architecture)
- [Troubleshooting](#troubleshooting)
  - [Common Issues](#common-issues)
  - [Debugging Tools](#debugging-tools)
- [Advanced Configuration](#advanced-configuration)
- [Production Deployment](#production-deployment)
- [Clean Up](#clean-up)

## Prerequisites

- [Minikube](https://minikube.sigs.k8s.io/docs/start/) v1.28+ for local development
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) v1.20+
- [Docker](https://docs.docker.com/get-docker/) v20+

## Setup Minikube

1. Start Minikube:

```bash
minikube start
```

2. Enable the Ingress addon:

```bash
minikube addons enable ingress
```

## Deployment

### Automated Deployment

The easiest way to deploy is using the provided deployment script:

```bash
# Make sure the script is executable
chmod +x deploy-minikube.sh

# Run the deployment script
./deploy-minikube.sh
```

The script will:
- Start Minikube if it's not running
- Enable the Ingress addon if needed
- Build the Docker image
- Load the image into Minikube
- Apply all Kubernetes manifests
- Wait for pods to be ready
- Display the Minikube IP and instructions for accessing the application

### Manual Deployment

If you prefer to deploy manually:

1. Build the Docker image:

```bash
docker build -t go-api:latest ..
```

2. Load the image into Minikube:

```bash
minikube image load go-api:latest
```

3. Apply the Kubernetes manifests:

```bash
# Apply all resources using kustomize
kubectl apply -k .

# Or apply each manifest individually
kubectl apply -f configmap.yaml
kubectl apply -f secrets.yaml
kubectl apply -f mysql.yaml
kubectl apply -f redis.yaml
kubectl apply -f deployment.yaml
```

4. Wait for pods to be ready:

```bash
kubectl wait --for=condition=ready pod -l app=go-api --timeout=180s
kubectl wait --for=condition=ready pod -l app=mysql --timeout=180s
kubectl wait --for=condition=ready pod -l app=redis --timeout=180s
```

## Access the Application

1. Get the Minikube IP:

```bash
minikube ip
```

2. Add an entry to your hosts file:

```
# /etc/hosts (Linux/Mac) or C:\Windows\System32\drivers\etc\hosts (Windows)
<minikube-ip> go-api.local
```

Example:
```
192.168.49.2 go-api.local
```

3. Access the application at:

```
http://go-api.local
```

## Deployment Architecture

The deployment consists of the following components:

- **API Service (go-api)**
  - Main REST API application
  - Uses init containers to wait for MySQL and Redis
  - Health checks on `/health` endpoint

- **MySQL Database**
  - Persistent storage for application data
  - Uses PersistentVolumeClaim for data persistence
  - Initialized with required databases and permissions

- **Redis Cache**
  - Used for caching and improving performance
  - Uses PersistentVolumeClaim for data persistence

- **ConfigMap**
  - Stores application configuration

- **Secrets**
  - Stores sensitive data like database credentials
  - Base64 encoded values

- **Ingress**
  - Provides external access to the API

## Troubleshooting

### Common Issues

1. **Pod in CrashLoopBackOff state**
   
   Check the pod logs:
   ```bash
   kubectl logs <pod-name>
   ```

   Common causes:
   - Database connection issues (incorrect dsn string)
   - Missing configuration
   - Resource limitations

2. **Init containers keep waiting**
   
   This usually means that the dependency services (MySQL/Redis) are not ready.
   Check their logs:
   ```bash
   kubectl logs <pod-name> -c <container-name>
   ```

3. **MySQL character set errors**
   
   Ensure the correct charset is used in the connection string:
   ```
   ?charset=utf8mb4
   ```
   Not `utf8mb` which is incomplete.

4. **Permission issues with PersistentVolumes**
   
   Check the PVC status and events:
   ```bash
   kubectl describe pvc <pvc-name>
   ```

### Debugging Tools

1. **Check pod status**
   ```bash
   kubectl get pods
   ```

2. **Describe a pod to see events**
   ```bash
   kubectl describe pod <pod-name>
   ```

3. **View logs**
   ```bash
   kubectl logs <pod-name>
   ```

4. **Shell into a pod**
   ```bash
   kubectl exec -it <pod-name> -- /bin/sh
   ```

5. **Port forwarding for direct testing**
   ```bash
   kubectl port-forward service/go-api 8080:80
   ```

## Advanced Configuration

### Scaling the API

You can scale the API horizontally:

```bash
kubectl scale deployment go-api --replicas=3
```

### Custom Resource Requirements

Modify the resource limits in `deployment.yaml` for different environments:

```yaml
resources:
  limits:
    cpu: "500m"
    memory: "512Mi"
  requests:
    cpu: "100m"
    memory: "128Mi"
```

### Adding Environment Variables

Add new environment variables to `configmap.yaml` or `secrets.yaml` as needed:

```yaml
data:
  NEW_VARIABLE: "value"
```

Then update the container env section in `deployment.yaml`:

```yaml
env:
  - name: NEW_VARIABLE
    valueFrom:
      configMapKeyRef:
        name: go-api-config
        key: NEW_VARIABLE
```

## Production Deployment

For production environments, consider:

1. **Secret Management**
   - Use a proper secrets manager instead of Kubernetes secrets
   - Consider using vault-injector or similar solutions

2. **External Database**
   - Use managed database services instead of in-cluster databases
   - Implement proper backup and disaster recovery

3. **Resource Tuning**
   - Set appropriate resource requests and limits
   - Implement autoscaling

4. **TLS**
   - Configure TLS termination
   - Use cert-manager for automated certificate management

5. **Monitoring**
   - Set up Prometheus and Grafana
   - Implement proper alerting

## Clean Up

To clean up all resources:

```bash
# Using the cleanup script
./cleanup-minikube.sh

# Or manually
kubectl delete -f deployment.yaml
kubectl delete -f redis.yaml
kubectl delete -f mysql.yaml
kubectl delete -f secrets.yaml
kubectl delete -f configmap.yaml
```

To stop Minikube:

```bash
minikube stop
``` 