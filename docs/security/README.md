# SREDIAG Security Guide

## Table of Contents

1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Authorization](#authorization)
4. [Network Security](#network-security)
5. [Data Protection](#data-protection)
6. [Compliance](#compliance)
7. [Best Practices](#best-practices)

## Overview

SREDIAG implements comprehensive security measures to protect your diagnostic data and system access. This guide covers all security aspects of the system.

## Authentication

### Methods

1. **JWT Authentication**

   ```yaml
   security:
     auth:
       type: jwt
       jwt:
         secret: ${JWT_SECRET}
         expiration: 24h
         public_key: /etc/srediag/keys/jwt.pub
         private_key: /etc/srediag/keys/jwt.key
   ```

2. **OAuth2 Integration**

   ```yaml
   security:
     auth:
       type: oauth2
       oauth2:
         provider_url: https://auth.example.com
         client_id: ${OAUTH_CLIENT_ID}
         client_secret: ${OAUTH_CLIENT_SECRET}
         scopes: [openid, profile, email]
   ```

3. **Certificate-based Authentication**

   ```yaml
   security:
     tls:
       enabled: true
       cert_file: /etc/srediag/certs/server.crt
       key_file: /etc/srediag/certs/server.key
       ca_file: /etc/srediag/certs/ca.crt
       verify_client: true
   ```

## Authorization

### Role-Based Access Control (RBAC)

```yaml
security:
  rbac:
    enabled: true
    default_role: viewer
    roles:
      admin:
        - "*"
      operator:
        - "read:*"
        - "write:config"
        - "write:telemetry"
      viewer:
        - "read:config"
        - "read:telemetry"
```

### Permission Levels

1. **Admin Permissions**
   - Full system access
   - User management
   - Configuration management
   - Security settings

2. **Operator Permissions**
   - Read all resources
   - Modify configurations
   - Manage telemetry
   - View audit logs

3. **Viewer Permissions**
   - Read-only access
   - View metrics
   - View logs
   - Access dashboards

## Network Security

### TLS Configuration

```yaml
security:
  tls:
    enabled: true
    min_version: TLS1.3
    cipher_suites:
      - TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
      - TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
    verify_client: true
```

### Network Policies

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: srediag-network-policy
spec:
  podSelector:
    matchLabels:
      app: srediag
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              name: monitoring
    ports:
      - protocol: TCP
        port: 8080
      - protocol: TCP
        port: 9090
```

## Data Protection

### Encryption at Rest

```yaml
security:
  encryption:
    enabled: true
    provider: aes
    key_file: /etc/srediag/keys/encryption.key
    algorithm: AES-256-GCM
```

### Data Retention

```yaml
security:
  retention:
    enabled: true
    metrics:
      duration: 30d
      size: 50GB
    logs:
      duration: 90d
      size: 100GB
    traces:
      duration: 7d
      size: 20GB
```

## Compliance

### Audit Logging

```yaml
security:
  audit:
    enabled: true
    log_path: /var/log/srediag/audit.log
    events:
      - authentication
      - authorization
      - configuration
      - data_access
```

### Compliance Standards

- SOC 2 Type II
- GDPR
- HIPAA
- PCI DSS

## Best Practices

### Security Hardening

1. **System Hardening**

   ```bash
   # Set proper file permissions
   chmod 600 /etc/srediag/keys/*
   chmod 644 /etc/srediag/config.yaml
   
   # Set ownership
   chown -R srediag:srediag /etc/srediag
   ```

2. **Container Security**

   ```yaml
   security:
     container:
       privileged: false
       read_only_root: true
       run_as_user: 1000
       run_as_group: 1000
       seccomp_profile: runtime/default
   ```

3. **Secret Management**

   ```yaml
   security:
     secrets:
       provider: vault
       path: secret/srediag
       auto_rotate: true
       rotation_period: 30d
   ```

### Security Monitoring

1. **Security Metrics**
   - Authentication attempts
   - Authorization failures
   - TLS handshake errors
   - Network policy violations

2. **Security Alerts**

   ```yaml
   alerts:
     security:
       - name: auth_failure
         condition: auth_failures > 10
         duration: 5m
         severity: critical
       - name: invalid_tokens
         condition: invalid_jwt_tokens > 5
         duration: 1m
         severity: warning
   ```

## See Also

- [Configuration Guide](../configuration/README.md)
- [Kubernetes Integration](../cloud/kubernetes.md)
- [Monitoring Guide](../configuration/telemetry.md)
- [Troubleshooting](../reference/troubleshooting.md)
