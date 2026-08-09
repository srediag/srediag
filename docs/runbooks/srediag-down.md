# Runbook: SREDIAG Down

## Alert Name
`SrediagDown`

## Severity
Critical

## Description
One or more SREDIAG instances are not responding to health checks and are considered down.

## Impact
- Telemetry data (traces, metrics, logs) from affected sources will not be collected
- Data loss may occur if buffer capacity is exceeded
- Downstream monitoring and observability systems may have gaps in data

## Diagnosis

### 1. Check Pod Status
```bash
kubectl get pods -n observability -l app.kubernetes.io/name=srediag
kubectl describe pod <pod-name> -n observability
```

### 2. Check Recent Events
```bash
kubectl get events -n observability --sort-by='.lastTimestamp' | grep srediag
```

### 3. Check Pod Logs
```bash
kubectl logs <pod-name> -n observability --tail=100
kubectl logs <pod-name> -n observability --previous # if pod restarted
```

### 4. Check Resource Usage
```bash
kubectl top pod -n observability -l app.kubernetes.io/name=srediag
kubectl describe node <node-name> # check if node has resource pressure
```

### 5. Check Network Connectivity
```bash
kubectl exec -it <pod-name> -n observability -- wget -O- http://localhost:13133/health
```

## Resolution Steps

### If Pod is in CrashLoopBackOff:
1. Check logs for startup errors
2. Verify configuration is valid:
   ```bash
   kubectl get configmap srediag-config -n observability -o yaml
   ```
3. Check if image is accessible:
   ```bash
   kubectl describe pod <pod-name> -n observability | grep -A5 "Events:"
   ```

### If Pod is in Pending State:
1. Check for resource constraints:
   ```bash
   kubectl describe pod <pod-name> -n observability | grep -A10 "Events:"
   kubectl get nodes -o wide
   kubectl describe node <node-name>
   ```
2. Check for pod anti-affinity conflicts
3. Verify ResourceQuota is not exceeded:
   ```bash
   kubectl describe resourcequota -n observability
   ```

### If Pod is Running but Unhealthy:
1. Check health endpoint manually:
   ```bash
   kubectl port-forward <pod-name> -n observability 13133:13133
   curl http://localhost:13133/health
   ```
2. Check for memory/CPU pressure:
   ```bash
   kubectl top pod <pod-name> -n observability
   ```
3. Review recent configuration changes
4. Check for network policy issues

### Emergency Recovery:
1. Scale down and up the deployment:
   ```bash
   kubectl scale deployment srediag -n observability --replicas=0
   kubectl scale deployment srediag -n observability --replicas=3
   ```
2. If configuration is suspected, rollback:
   ```bash
   kubectl rollout undo deployment/srediag -n observability
   ```

## Escalation
If the issue persists after 15 minutes:
1. Contact the on-call SRE engineer
2. Check if this is part of a wider infrastructure issue
3. Consider failing over to backup telemetry collection if available

## Prevention
1. Ensure proper resource limits are set
2. Implement proper health checks
3. Use PodDisruptionBudgets to prevent all pods being evicted
4. Regular configuration validation
5. Implement canary deployments for configuration changes

## Related Documentation
- [SREDIAG Architecture](../architecture.md)
- [Configuration Guide](../configuration.md)
- [Troubleshooting Guide](../troubleshooting.md)