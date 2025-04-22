# SREDIAG Security Configuration

## Overview

This document describes the security configurations available in SREDIAG, including authentication, authorization, encryption, and other protection measures.

## Configuration Structure

```yaml
security:
  # TLS Configuration
  tls:
    enabled: ${TLS_ENABLED:-true}
    cert_file: ${TLS_CERT_FILE:-/etc/srediag/certs/server.crt}
    key_file: ${TLS_KEY_FILE:-/etc/srediag/certs/server.key}
    ca_file: ${TLS_CA_FILE:-/etc/srediag/certs/ca.crt}
    min_version: ${TLS_MIN_VERSION:-TLS1.2}
    verify_client: ${TLS_VERIFY_CLIENT:-true}
    cipher_suites:
      - TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
      - TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
  
  # Authentication Configuration
  auth:
    enabled: ${AUTH_ENABLED:-true}
    type: ${AUTH_TYPE:-jwt}
    jwt:
      secret: ${JWT_SECRET:-}
      public_key: ${JWT_PUBLIC_KEY:-/etc/srediag/keys/jwt.pub}
      private_key: ${JWT_PRIVATE_KEY:-/etc/srediag/keys/jwt.key}
      expiration: ${JWT_EXPIRATION:-24h}
    oauth2:
      provider_url: ${OAUTH_PROVIDER_URL:-}
      client_id: ${OAUTH_CLIENT_ID:-}
      client_secret: ${OAUTH_CLIENT_SECRET:-}
      scopes: ${OAUTH_SCOPES:-[openid, profile, email]}
  
  # RBAC Configuration
  rbac:
    enabled: ${RBAC_ENABLED:-true}
    default_role: ${RBAC_DEFAULT_ROLE:-viewer}
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
  
  # Network Configuration
  network:
    allowed_origins: ${ALLOWED_ORIGINS:-["*"]}
    allowed_methods: ${ALLOWED_METHODS:-["GET", "POST", "PUT", "DELETE"]}
    allowed_headers: ${ALLOWED_HEADERS:-["Authorization", "Content-Type"]}
    exposed_headers: ${EXPOSED_HEADERS:-["Content-Length"]}
    allow_credentials: ${ALLOW_CREDENTIALS:-true}
    max_age: ${CORS_MAX_AGE:-86400}
  
  # Rate Limiting Configuration
  rate_limit:
    enabled: ${RATE_LIMIT_ENABLED:-true}
    requests_per_second: ${RATE_LIMIT_RPS:-100}
    burst: ${RATE_LIMIT_BURST:-200}
```

## Components

### TLS

1. **Certificates**
   - Server certificate
   - Private key
   - Certificate Authority (CA)
   - Client validation

2. **Versions and Ciphers**
   - Minimum TLS 1.2 version
   - Secure cipher suites
   - Priority configuration

### Authentication

1. **JWT**
   - Token generation
   - Token validation
   - Public/private keys
   - Expiration configuration

2. **OAuth2**
   - Supported providers
   - Authentication flows
   - Required scopes
   - Token validation

### RBAC (Role-Based Access Control)

1. **Roles**
   - Administrator
   - Operator
   - Viewer
   - Custom roles

2. **Permissions**
   - Configuration read
   - Configuration write
   - Telemetry access
   - User management

### Network Security

1. **CORS**
   - Allowed origins
   - Allowed methods
   - Allowed headers
   - Credentials

2. **Rate Limiting**
   - IP-based limits
   - User-based limits
   - Burst configuration
   - Blocking policies

## Configuration Examples

### Basic Configuration

```yaml
security:
  tls:
    enabled: true
    cert_file: /etc/srediag/certs/server.crt
    key_file: /etc/srediag/certs/server.key
  
  auth:
    enabled: true
    type: jwt
    jwt:
      secret: "your-secret-here"
  
  rbac:
    enabled: true
    default_role: viewer
  
  rate_limit:
    enabled: true
    requests_per_second: 100
```

### Production Configuration

```yaml
security:
  tls:
    enabled: true
    cert_file: /etc/srediag/certs/server.crt
    key_file: /etc/srediag/certs/server.key
    ca_file: /etc/srediag/certs/ca.crt
    min_version: TLS1.3
    verify_client: true
    cipher_suites:
      - TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
      - TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
  
  auth:
    enabled: true
    type: oauth2
    oauth2:
      provider_url: "https://your-oauth-provider.com"
      client_id: "your-client-id"
      client_secret: "your-client-secret"
      scopes: ["openid", "profile", "email"]
  
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
  
  network:
    allowed_origins: ["https://your-application.com"]
    allowed_methods: ["GET", "POST", "PUT", "DELETE"]
    allowed_headers: ["Authorization", "Content-Type"]
    exposed_headers: ["Content-Length"]
    allow_credentials: true
    max_age: 86400
  
  rate_limit:
    enabled: true
    requests_per_second: 100
    burst: 200
```

## Best Practices

1. **Certificate Management**
   - Use valid certificates
   - Rotate regularly
   - Protect private keys
   - Keep CAs updated

2. **Authentication and Authorization**
   - Use OAuth2 in production
   - Implement RBAC
   - Limit permissions
   - Audit access

3. **Network Security**
   - Configure CORS
   - Implement rate limiting
   - Use TLS 1.3
   - Monitor traffic

4. **Data Protection**
   - Encrypt sensitive data
   - Implement backup
   - Define retention
   - Monitor access

## See Also

- [Configuration Overview](README.md)
- [Collector Configuration](collector.md)
- [Telemetry Configuration](telemetry.md)
- [Troubleshooting](../reference/troubleshooting.md)
