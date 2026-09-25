# ADR-002: sqlc en lugar de ORM

## Contexto

Necesitamos acceso a datos PostgreSQL type-safe desde Go. Opciones principales:
- ORM completo (GORM, Ent, sqlboiler)
- Query builder (sqlc, squirrel, bun)
- Raw `database/sql` + scanning manual

Requisitos: Type safety, performance, mantenimiento simple, migraciones controladas.

## Decisión

**Usar `sqlc` (SQL Compiler) para generar código Go desde queries SQL escritas a mano.**

- Esquema definido en `backend/sql/schema.sql` (DDL)
- Queries en `backend/sql/queries.sql` (SQL parametrizado)
- `sqlc.yaml` configura generación a `internal/infrastructure/db/`
- Migraciones manuales via archivos SQL versionados

## Consecuencias

### Positivas
- **SQL real**: Control total sobre queries, EXPLAIN plans, índices, hints
- **Type safety end-to-end**: Structs Go generados coinciden exactamente con schema
- **Zero runtime overhead**: Código generado usa `pgx` directamente, sin reflection
- **Migraciones explícitas**: Archivos SQL versionados en git, revisables en PR
- **Testing fácil**: `sqlc` genera interfaces para mocking en tests
- **Performance**: Queries optimizadas a mano, N+1 evitado por diseño

### Negativas
- **Boilerplate inicial**: Configuración `sqlc.yaml`, structs de queries
- **Cambios de schema = regeneración**: `sqlc compile` + `sqlc generate` tras migraciones
- **Menos "mágico"**: No auto-migrations, no relations loading automático

### Neutrales
- Curva de aprendizaje SQL puro para devs acostumbrados a ORM
- `sqlc` no valida SQL en compile-time (solo en generate-time)

## Alternativas Consideradas

| Opción | Por qué NO |
|--------|------------|
| **GORM** | Reflection overhead, magic callbacks, N+1 frecuente, migraciones opacas |
| **Ent** | Code-gen pesado, schema en Go (no SQL), learning curve alta |
| **sqlboiler** | Requiere DB viva para generar, regeneración lenta |
| **Raw database/sql** | Boilerplate scanning manual, propenso a errores, sin type safety |

## Referencias
- [backend/sqlc.yaml](../../../backend/sqlc.yaml)
- [backend/sql/schema.sql](../../../backend/sql/schema.sql)
- [backend/sql/queries.sql](../../../backend/sql/queries.sql)