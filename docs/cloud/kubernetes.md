# SREDIAG Kubernetes Integration

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Installation](#installation)
4. [Configuration](#configuration)
5. [Deployment](#deployment)
6. [Monitoring](#monitoring)
7. [Troubleshooting](#troubleshooting)

## Overview

SREDIAG provides native Kubernetes integration, enabling efficient monitoring and diagnostics of Kubernetes clusters. This integration includes:

- Automatic cluster metrics collection
- Service and pod discovery
- Cluster resource monitoring
- Custom metrics integration
- Auto-scaling support

## Prerequisites

### Cluster Requirements

- Kubernetes 1.24+
- Helm 3.x
- Metrics Server enabled
- RBAC enabled

### Minimum Resources

```yaml
resources:
  requests:
    cpu: 200m
    memory: 256Mi
  limits:
    cpu: 1000m
    memory: 1Gi
```

## Installation

### Using Helm

```bash
# Add Helm repository
helm repo add srediag https://charts.srediag.io
helm repo update

# Install SREDIAG
helm install srediag srediag/srediag \
  --namespace monitoring \
  --create-namespace \
  --values values.yaml
```

### Basic values.yaml File

```yaml
srediag:
  image:
    repository: srediag/srediag
    tag: latest
    pullPolicy: IfNotPresent

  serviceAccount:
    create: true
    annotations: {}
    name: ""

  rbac:
    create: true

  config:
    telemetry:
      enabled: true
      endpoint: "http://collector:4317"
    
    kubernetes:
      enabled: true
      discovery:
        enabled: true
        interval: 10s

  resources:
    requests:
      cpu: 200m
      memory: 256Mi
    limits:
      cpu: 1000m
      memory: 1Gi

  nodeSelector: {}
  tolerations: []
  affinity: {}
```

## Configuration

### Collector Configuration

```yaml
collector:
  receivers:
    k8s_cluster:
      collection_interval: 10s
      node_conditions_to_report: ["Ready", "MemoryPressure", "DiskPressure"]
      allocatable_types_to_report: ["cpu", "memory", "pods"]
      
    kubeletstats:
      collection_interval: 10s
      auth_type: "serviceAccount"
      endpoint: "${K8S_NODE_IP}:10250"
      
  processors:
    k8sattributes:
      auth_type: "serviceAccount"
      passthrough: false
      extract:
        metadata:
          - k8s.pod.name
          - k8s.pod.uid
          - k8s.deployment.name
          - k8s.namespace.name
          - k8s.node.name
          
  exporters:
    prometheus:
      endpoint: "0.0.0.0:9090"
      namespace: srediag
      const_labels:
        cluster: "${CLUSTER_NAME}"
```

### RBAC

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: srediag
  namespace: monitoring

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: srediag
rules:
  - apiGroups: [""]
    resources:
      - nodes
      - pods
      - services
      - endpoints
      - namespaces
    verbs: ["get", "list", "watch"]
  - apiGroups: ["apps"]
    resources:
      - deployments
      - daemonsets
      - statefulsets
    verbs: ["get", "list", "watch"]
  - apiGroups: ["metrics.k8s.io"]
    resources:
      - pods
      - nodes
    verbs: ["get", "list", "watch"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: srediag
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: srediag
subjects:
  - kind: ServiceAccount
    name: srediag
    namespace: monitoring
```

## Deployment

### Basic Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: srediag
  namespace: monitoring
spec:
  replicas: 1
  selector:
    matchLabels:
      app: srediag
  template:
    metadata:
      labels:
        app: srediag
    spec:
      serviceAccountName: srediag
      containers:
        - name: srediag
          image: srediag/srediag:latest
          ports:
            - containerPort: 8080
              name: http
            - containerPort: 9090
              name: metrics
          env:
            - name: KUBERNETES_NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: KUBERNETES_POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
          resources:
            requests:
              cpu: 200m
              memory: 256Mi
            limits:
              cpu: 1000m
              memory: 1Gi
          livenessProbe:
            httpGet:
              path: /health
              port: http
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /ready
              port: http
            initialDelaySeconds: 5
            periodSeconds: 5
```

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: srediag
  namespace: monitoring
spec:
  type: ClusterIP
  ports:
    - port: 8080
      targetPort: http
      protocol: TCP
      name: http
    - port: 9090
      targetPort: metrics
      protocol: TCP
      name: metrics
  selector:
    app: srediag
```

## Monitoring

### Collected Metrics

1. **Cluster Metrics**
   - CPU/Memory utilization per node
   - Number of pods per node
   - Node status
   - Cluster events

2. **Pod Metrics**
   - Resource usage
   - Pod status
   - Restarts
   - Latency

3. **Application Metrics**
   - Custom metrics
   - Request latency
   - Error rates
   - Throughput

### Dashboards

SREDIAG provides pre-configured Kubernetes dashboards:

1. **Cluster Overview**
   - General status
   - Resource utilization
   - Important events

2. **Pod Analysis**
   - Performance
   - Logs
   - Events
   - Resources

3. **Application Metrics**
   - Business metrics
   - SLOs/SLIs
   - Alerts

## Troubleshooting

### Common Issues

1. **Metrics Collection Failure**

   ```bash
   # Check pod logs
   kubectl logs -n monitoring deploy/srediag
   
   # Verify permissions
   kubectl auth can-i get pods --as=system:serviceaccount:monitoring:srediag
   ```

2. **Resource Issues**

   ```bash
   # Check resource usage
   kubectl top pod -n monitoring
   
   # Describe pod
   kubectl describe pod -n monitoring -l app=srediag
   ```

3. **Connectivity Issues**

   ```bash
   # Test connectivity
   kubectl exec -n monitoring deploy/srediag -- curl -s http://localhost:8080/health
   
   # Check endpoints
   kubectl get endpoints -n monitoring srediag
   ```

### Useful Commands

```bash
# Restart deployment
kubectl rollout restart deploy/srediag -n monitoring

# Check logs
kubectl logs -f deploy/srediag -n monitoring

# Scale deployment
kubectl scale deploy/srediag -n monitoring --replicas=3

# Update configuration
kubectl create configmap srediag-config -n monitoring --from-file=config.yaml
```

## See Also

- [Advanced Configuration](../configuration/README.md)
- [Telemetry](../configuration/telemetry.md)
- [Security](../security/README.md)
- [Troubleshooting](../reference/troubleshooting.md)
