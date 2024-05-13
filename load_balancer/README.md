# Load balancer

## Tools
- Docker
- Kubernetes (minikube)

## Set up:
- Install Docker
- Install Minikube
- Install kubectl

### Set up database:
```bash
cd load_balancer
kubectl apply db/
```

### Set up service:
```bash
kubectl apply service/
```

### To get the service url:
```bash
minikube service job-service --url
```

## To test if that work
```bash
# Get all pods
kubectl get pods
# Select 1 pod name
kubetctl delete pod <name>
```
New pods will be created to match the replica numbers