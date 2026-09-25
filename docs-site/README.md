# 📚 Documentation Scripts

Scripts para generar y gestionar la documentación técnica del proyecto LVBP GameCast.

## 📦 Requisitos Previos

| Herramienta | Versión | Instalación |
|-------------|---------|-------------|
| **Node.js** | 18+ | `winget install OpenJS.NodeJS` / `brew install node` |
| **pnpm** | 9+ | `npm install -g pnpm` |
| **Go** | 1.22+ | `winget install GoLang.Go` / `brew install go` |
| **Docker** | 24+ | Docker Desktop (para BD local) |

## 🚀 Comandos Principales

```bash
# Build completo (OpenAPI + TypeDoc + VitePress)
npm run build-docs

# Solo VitePress (docs-site ya generado)
npm run build-docs:vitepress

# Solo OpenAPI spec (backend Go)
npm run build-docs:openapi

# Solo TypeDoc (frontend TypeScript)
npm run build-docs:typedoc

# Modo CI/CD (falla rápido, sin dev server)
npm run build-docs:ci

# Modo desarrollo (build + servidor dev en :5173)
npm run dev-docs

# Limpiar artefactos
npm run clean-docs
```

## 🖥️ Ejecución Directa (Cross-platform)

```bash
# Windows (PowerShell/CMD)
node scripts\build-docs.js
node scripts\build-docs.js --ci
node scripts\build-docs.js --dev
node scripts\build-docs.js --skip-openapi --skip-typedoc

# Linux/macOS/WSL/Git Bash
node scripts/build-docs.js
node scripts/build-docs.js --ci
node scripts/build-docs.js --dev
```

## ⚙️ Opciones del Script

| Flag | Descripción |
|------|-------------|
| `--ci` | Modo CI: exit code ≠ 0 en error, no inicia dev server |
| `--dev` | Build + inicia VitePress dev server en `http://localhost:5173/lvbp-project/` |
| `--skip-openapi` | Omite generación OpenAPI (swag) |
| `--skip-typedoc` | Omite generación TypeDoc |
| `--skip-vitepress` | Omite build VitePress |
| `--help` | Muestra ayuda |

## 🐳 Docker (Desarrollo Local)

```bash
# Levantar infraestructura (PostgreSQL + Redis)
docker compose up -d

# Ver logs
docker compose logs -f postgres
docker compose logs -f redis

# Parar
docker compose down

# Parar + volúmenes (reset BD)
docker compose down -v
```

### Variables de Entorno (`.env` en raíz)

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=lvbp
DB_PASSWORD=lvbp
DB_NAME=lvbp
DB_SSLMODE=disable

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Auth
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION=24h

# Server
SERVER_PORT=8080
GIN_MODE=debug
```

## 🔧 Generación Manual de Artefactos

### OpenAPI Spec (Backend)

```bash
cd backend
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go -o ../docs-site/public/api --parseDependency --parseInternal
```

### TypeDoc (Frontend)

```bash
cd web
pnpm install --frozen-lockfile
npx typedoc --out ../docs-site/public/typedoc --entryPoints src/lib --plugin typedoc-plugin-markdown --readme none
```

### VitePress (Docs Site)

```bash
cd docs-site
pnpm install --frozen-lockfile
pnpm run build        # Producción
pnpm run dev          # Desarrollo
```

## 📁 Estructura de Salida

```
docs-site/
├── public/
│   ├── api/              # OpenAPI spec (openapi.yaml)
│   └── typedoc/          # TypeScript docs (markdown)
├── .vitepress/
│   └── dist/             # Sitio estático generado
└── .github/workflows/
    └── deploy-docs.yml   # GitHub Pages deploy
```

## 🌐 Despliegue Automático

El sitio se despliega a **GitHub Pages** automáticamente en push a `main`:

1. Push a `main` → GitHub Actions ejecuta `build-docs.js --ci`
2. Build artifacts → `actions/upload-pages-artifact`
3. Deploy → `actions/deploy-pages`

URL: `https://<user>.github.io/lvbp-project/`

## 🛠️ Troubleshooting

### Error: "swag not found"
```bash
cd backend && go install github.com/swaggo/swag/cmd/swag@latest
```

### Error: "typedoc not found"
```bash
cd web && pnpm install --frozen-lockfile
```

### Error: "vitepress not found"
```bash
cd docs-site && pnpm install --frozen-lockfile
```

### PowerShell: "running scripts is disabled"
```powershell
Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process
# O usa npm/pnpm directamente:
npm run build-docs
```

### Puerto 5173 ocupado
```bash
# El script usa --port 5173 --host 0.0.0.0
# Cambia en package.json si necesitas otro puerto
```

## 📝 Notas

- El script `build-docs.js` es **puro Node.js ES Modules** - funciona en Windows, Linux, macOS
- Auto-instala dependencias faltantes (Go modules, pnpm packages)
- Usa `spawnSync` con `shell: true` para comandos cross-platform
- En CI usa `--ci` para fallar rápido con exit codes correctos