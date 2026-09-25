# Runbook: Despliegue a Producción

## Arquitectura de Despliegue

```
┌─────────────────────────────────────────────────────────────────┐
│                     GitHub Actions CI/CD                        │
├─────────────────────────────────────────────────────────────────┤
│  Push main → Build → Test → Docker Build → Push Registry       │
│                              ↓                                  │
│                    Deploy to Kubernetes                         │
└─────────────────────────────────────────────────────────────────┘
```

## Variables de Entorno Producción

### Backend
```env
# Database (Managed PostgreSQL - Cloud SQL, RDS, etc.)
DB_HOST=prod-postgres.internal
DB_PORT=5432
DB_USER=lvbp_app
DB_PASSWORD=<from-secret-manager>
DB_NAME=lvbp_prod
DB_SSLMODE=require
DB_MAX_CONNS=25
DB_MIN_CONNS=5

# Redis (Managed - ElastiCache, Memorystore, etc.)
REDIS_ADDR=prod-redis.internal:6379
REDIS_PASSWORD=<from-secret-manager>
REDIS_DB=0
REDIS_POOL_SIZE=50

# Auth
JWT_SECRET=<from-secret-manager-256-bit>
JWT_EXPIRATION=24h

# Server
SERVER_PORT=8080
GIN_MODE=release
LOG_LEVEL=info

# CORS
CORS_ALLOWED_ORIGINS=https://lvbp.example.com,https://app.lvbp.example.com
```

### Frontend
```env
NEXT_PUBLIC_API_URL=https://api.lvbp.example.com
NEXT_PUBLIC_WS_URL=wss://api.lvbp.example.com
NEXT_TELEMETRY_DISABLED=1
```

## Docker Images

### Backend Dockerfile (`backend/Dockerfile`)
```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git make
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /lvbp-api ./cmd/api

# Runtime stage
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /lvbp-api .
COPY --from=builder /app/sql/migrations ./sql/migrations
EXPOSE 8080
ENTRYPOINT ["./lvbp-api"]
```

### Frontend Dockerfile (`web/Dockerfile`)
```dockerfile
# Build stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN corepack enable pnpm && pnpm install --frozen-lockfile
COPY . .
RUN pnpm build

# Runtime stage
FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
EXPOSE 3000
CMD ["node", "server.js"]
```

**Nota**: Requiere `output: 'standalone'` en `next.config.ts`

## Kubernetes Manifests

### Namespace
```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: lvbp-prod
  labels:
    name: lvbp-prod
```

### Backend Deployment
```yaml
# k8s/backend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: lvbp-api
  namespace: lvbp-prod
spec:
  replicas: 3
  selector:
    matchLabels:
      app: lvbp-api
  template:
    metadata:
      labels:
        app: lvbp-api
    spec:
      containers:
      - name: api
        image: ghcr.io/your-org/lvbp-api:v1.2.3
        ports:
        - containerPort: 8080
        envFrom:
        - secretRef:
            name: lvbp-api-secrets
        - configMapRef:
            name: lvbp-api-config
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: lvbp-api
  namespace: lvbp-prod
spec:
  selector:
    app: lvbp-api
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

### Frontend Deployment
```yaml
# k8s/frontend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: lvbp-web
  namespace: lvbp-prod
spec:
  replicas: 3
  selector:
    matchLabels:
      app: lvbp-web
  template:
    metadata:
      labels:
        app: lvbp-web
    spec:
      containers:
      - name: web
        image: ghcr.io/your-org/lvbp-web:v1.2.3
        ports:
        - containerPort: 3000
        envFrom:
        - configMapRef:
            name: lvbp-web-config
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "128Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 10
---
apiVersion: v1
kind: Service
metadata:
  name: lvbp-web
  namespace: lvbp-prod
spec:
  selector:
    app: lvbp-web
  ports:
  - port: 80
    targetPort: 3000
```

### Ingress (nginx)
```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: lvbp-ingress
  namespace: lvbp-prod
  annotations:
    nginx.ingress.kubernetes.io/proxy-buffering: "false"
    nginx.ingress.kubernetes.io/proxy-cache: "false"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - lvbp.example.com
    - api.lvbp.example.com
    secretName: lvbp-tls
  rules:
  - host: lvbp.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: lvbp-web
            port:
              number: 80
  - host: api.lvbp.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: lvbp-api
            port:
              number: 80
```

**Crítico para SSE**: `proxy-buffering: "false"` + `proxy-cache: "false"` (NFR-3)

### ConfigMap & Secrets
```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: lvbp-api-config
  namespace: lvbp-prod
data:
  DB_HOST: "prod-postgres.internal"
  DB_PORT: "5432"
  DB_NAME: "lvbp_prod"
  REDIS_ADDR: "prod-redis.internal:6379"
  SERVER_PORT: "8080"
  GIN_MODE: "release"
  LOG_LEVEL: "info"
---
apiVersion: v1
kind: Secret
metadata:
  name: lvbp-api-secrets
  namespace: lvbp-prod
type: Opaque
stringData:
  DB_USER: "lvbp_app"
  DB_PASSWORD: "<from-vault>"
  REDIS_PASSWORD: "<from-vault>"
  JWT_SECRET: "<from-vault-256-bit>"
```

## GitHub Actions Workflow

```yaml
# .github/workflows/deploy.yml
name: Deploy to Production

on:
  push:
    branches: [main]
  workflow_dispatch:

env:
  REGISTRY: ghcr.io
  IMAGE_NAME_API: ${{ github.repository_owner }}/lvbp-api
  IMAGE_NAME_WEB: ${{ github.repository_owner }}/lvbp-web

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: Run Backend Tests
        run: |
          cd backend
          go test ./... -race -coverprofile=coverage.out
      
      - name: Set up Node
        uses: actions/setup-node@v4
        with:
          node-version: '20'
      
      - name: Run Frontend Tests
        run: |
          cd web
          pnpm install --frozen-lockfile
          pnpm test
          pnpm type-check

  build-api:
    needs: test
    runs-on: ubuntu-latest
    permissions:
      packages: write
      contents: read
    steps:
      - uses: actions/checkout@v4
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      - name: Log in to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME_API }}
          tags: |
            type=ref,event=branch
            type=sha
      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: ./backend
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

  build-web:
    needs: test
    runs-on: ubuntu-latest
    permissions:
      packages: write
      contents: read
    steps:
      - uses: actions/checkout@v4
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      - name: Log in to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: ./web
          push: true
          tags: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME_WEB }}:latest

  deploy:
    needs: [build-api, build-web]
    runs-on: ubuntu-latest
    environment: production
    steps:
      - uses: actions/checkout@v4
      - name: Set up kubectl
        uses: azure/k8s-set-context@v1
        with:
          kubeconfig: ${{ secrets.KUBECONFIG_PROD }}
      - name: Deploy to Kubernetes
        run: |
          kubectl set image deployment/lvbp-api \
            api=${{ env.REGISTRY }}/${{ env.IMAGE_NAME_API }}:${{ github.sha }} \
            -n lvbp-prod
          kubectl set image deployment/lvbp-web \
            web=${{ env.REGISTRY }}/${{ env.IMAGE_NAME_WEB }}:${{ github.sha }} \
            -n lvbp-prod
          kubectl rollout status deployment/lvbp-api -n lvbp-prod --timeout=300s
          kubectl rollout status deployment/lvbp-web -n lvbp-prod --timeout=300s
```

## Health Checks

### Backend Endpoints
```go
// cmd/api/main.go - añadir
router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Write([]byte("OK"))
})

router.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
    // Verificar DB + Redis
    if err := db.Ping(r.Context()); err != nil {
        http.Error(w, "DB not ready", 503)
        return
    }
    if err := redisClient.Ping(r.Context()).Err(); err != nil {
        http.Error(w, "Redis not ready", 503)
        return
    }
    w.WriteHeader(200)
    w.Write([]byte("READY"))
})
```

## Monitoreo y Alertas

### Métricas Clave (Prometheus)

| Métrica | Tipo | Alerta |
|---------|------|--------|
| `http_requests_total` | Counter | Rate > 1000/s |
| `http_request_duration_seconds` | Histogram | p99 > 500ms |
| `sse_connections_active` | Gauge | > 4500 (de 5000 max) |
| `db_connections_active` | Gauge | > 80% pool |
| `redis_memory_used_bytes` | Gauge | > 80% maxmemory |
| `go_goroutines` | Gauge | > 10000 |

### Logs Estructurados (Loki/ELK)
```json
{
  "timestamp": "2026-09-24T19:30:00Z",
  "level": "INFO",
  "service": "lvbp-api",
  "trace_id": "abc-123",
  "message": "pitch_processed",
  "game_id": "uuid",
  "pitch_result": "Swinging Strike",
  "latency_ms": 12
}
```

### Alertas Críticas (PrometheusRule)
```yaml
groups:
- name: lvbp-critical
  rules:
  - alert: APIHighLatency
    expr: histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m])) > 0.5
    for: 2m
    labels:
      severity: critical
    annotations:
      summary: "API p99 latency > 500ms"
      
  - alert: SSEConnectionsNearLimit
    expr: sse_connections_active > 4500
    for: 1m
    labels:
      severity: warning
    annotations:
      summary: "SSE connections near 5000 limit"
      
  - alert: DatabaseConnectionsExhausted
    expr: db_connections_active / db_connections_max > 0.8
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "DB connection pool > 80%"
```

## Rollback Procedure

```bash
# 1. Verificar deployment actual
kubectl rollout history deployment/lvbp-api -n lvbp-prod

# 2. Rollback a versión anterior
kubectl rollout undo deployment/lvbp-api -n lvbp-prod

# 3. Verificar rollback
kubectl rollout status deployment/lvbp-api -n lvbp-prod

# 4. Si hay migración DB involucrada
#    - Revertir migración (si reversible)
#    - O restaurar backup DB point-in-time
```

## Backup & Disaster Recovery

### PostgreSQL (Managed)
- **Automated**: Daily snapshots + Point-in-time recovery (7 days)
- **Manual**: `pg_dump` antes de migraciones riesgosas
- **RTO**: < 15 min, **RPO**: < 1 min

### Redis (Managed)
- **AOF/RDB**: Cada 1 segundo + snapshot cada hora
- **Datos**: Solo sessions + cache (recuperable)

### Documentación
- Runbooks en `docs-site/runbooks/`
- Diagramas arquitectura en `docs-site/architecture/`
- ADRs en `docs-site/architecture/adr/`

## Checklist Pre-Deploy

- [ ] Tests pasan (backend + frontend)
- [ ] Build Docker exitoso
- [ ] Migraciones DB probadas en staging
- [ ] Secrets actualizados en Vault/Secret Manager
- [ ] ConfigMap revisado
- [ ] Ingress TLS válido (cert-manager)
- [ ] Health checks responden
- [ ] Monitoreo/alertas configurados
- [ ] Rollback plan documentado
- [ ] Equipo notificado (Slack/email)

## Contactos de Escalación

| Rol | Contacto | Canal |
|-----|----------|-------|
| On-call Engineer | @dev-oncall | PagerDuty / Slack #alerts |
| Backend Lead | @backend-lead | Slack #backend |
| Frontend Lead | @frontend-lead | Slack #frontend |
| DBA | @dba-team | Slack #database |
| Infra/Platform | @platform-team | Slack #platform |