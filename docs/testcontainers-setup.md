# Testcontainers Setup for Local Development

## Prerequisites

1. **Docker Desktop** with WSL2 backend enabled
2. **WSL2** installed and configured
3. **Go 1.22+** installed

## Docker Desktop Configuration

### Enable WSL2 Backend
1. Open Docker Desktop Settings
2. Go to **General** → Check **"Use the WSL 2 based engine"**
3. Go to **Resources** → **WSL Integration** → Enable integration with your WSL2 distribution (e.g., Ubuntu)

### Verify Docker is Running
```powershell
# Check Docker version and daemon status
docker version

# Should show both Client and Server versions
```

## Testcontainers Configuration

The project includes a `.testcontainers.properties` file in the root and `backend/` directories with the following settings:

```properties
# Disable rootless mode (not supported on Windows)
testcontainers.rootless.enabled=false

# Reuse containers for faster test execution
testcontainers.reuse.enable=true

# Disable Ryuk (testcontainers reaper) to avoid issues in CI environments
testcontainers.ryuk.disabled=true
```

### Environment Variables (Alternative)
You can also configure testcontainers via environment variables:

```powershell
# PowerShell
$env:TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE = "npipe:////./pipe/dockerDesktopLinuxEngine"
$env:TESTCONTAINERS_RYUK_DISABLED = "true"
$env:TESTCONTAINERS_REUSE_ENABLE = "true"
```

```bash
# Bash/WSL
export TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE="unix:///var/run/docker.sock"
export TESTCONTAINERS_RYUK_DISABLED="true"
export TESTCONTAINERS_REUSE_ENABLE="true"
```

## Running Tests Locally

### Backend Tests
```bash
cd backend
go test -v ./...
```

### Run Specific Integration Tests
```bash
# Game repository tests
go test -v ./internal/infrastructure/repositories/... -run "TestGameRepository"

# Redis pubsub tests
go test -v ./internal/infrastructure/pubsub/... -run "TestRedisPublisher"

# Redis client tests
go test -v ./internal/pkg/redis/... -run "TestPing"

# SSE stream handler tests
go test -v ./internal/handlers/sse/... -run "TestStreamGameEvents"
```

## Troubleshooting

### Error: "rootless Docker is not supported on Windows"
**Cause**: Testcontainers is trying to use rootless Docker mode which isn't supported on Windows.

**Solution**: 
1. Ensure Docker Desktop is running with WSL2 backend
2. Verify `.testcontainers.properties` has `testcontainers.rootless.enabled=false`
3. Restart your IDE/terminal to reload environment variables

### Error: "failed to create Docker provider" or "docker info: failed to connect"
**Cause**: Docker daemon is not accessible.

**Solution**:
1. Ensure Docker Desktop is running (check system tray)
2. Run `docker version` to verify both client and server respond
3. If using WSL2, ensure WSL integration is enabled in Docker Desktop settings

### Tests Timeout or Hang
**Cause**: Ryuk (testcontainers reaper) may have issues in some environments.

**Solution**: Ensure `testcontainers.ryuk.disabled=true` is set in `.testcontainers.properties` or `TESTCONTAINERS_RYUK_DISABLED=true` environment variable.

## CI/CD Configuration

The GitHub Actions workflow (`.github/workflows/ci.yml`) is configured to:
- Use Ubuntu runners with Docker pre-installed
- Set `TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE=/var/run/docker.sock`
- Disable Ryuk with `TESTCONTAINERS_RYUK_DISABLED=true`
- Run both unit tests and integration tests with testcontainers

## VS Code / IDE Configuration

For VS Code, add to `.vscode/settings.json`:
```json
{
  "go.testEnvVars": {
    "TESTCONTAINERS_RYUK_DISABLED": "true",
    "TESTCONTAINERS_REUSE_ENABLE": "true"
  }
}
```

For GoLand/IntelliJ, add environment variables in Run Configuration:
- `TESTCONTAINERS_RYUK_DISABLED=true`
- `TESTCONTAINERS_REUSE_ENABLE=true`