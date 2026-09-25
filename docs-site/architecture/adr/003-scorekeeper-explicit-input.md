# ADR-003: Scorekeeper sigue Input Explícito (No Infiere de Coordenadas)

## Contexto

El PRD especifica que el scorekeeper registra:
1. Click en zona de strike → coordenadas X,Y
2. Selecciona resultado explícito: `Ball`, `Called Strike`, `Swinging Strike`, `Foul`, `In Play`, `Hit By Pitch`

**Regla clave del PRD**: "El sistema no deduce balls/strikes por coordenadas físicas; obedece el input del scorekeeper."

Esto significa: las coordenadas son solo para visualización (strike zone SVG), NO para lógica de reglas.

## Decisión

**El motor de reglas (Scorekeeper Engine) procesa ÚNICAMENTE el `pitchResult` explícito enviado por el scorekeeper.**

- Coordenadas X,Y y velocidad se persisten para analytics/visualización
- Lógica de conteo (balls/strikes/outs) deriva 100% de `pitchResult` enum
- No hay algoritmo de "zona de strike automática" ni ML

```go
// En scorekeeper.go - RecordPitch
switch pitch.PitchResult {
case "Ball":
    newBalls++
    if newBalls == 4 { atBatResult = "Base on Balls"; isCompleted = true }
case "Called Strike", "Swinging Strike":
    newStrikes++
    if newStrikes == 3 { atBatResult = "Strikeout"; outsRecorded = 1; isCompleted = true }
case "Foul Ball":
    if newStrikes < 2 { newStrikes++ }
case "In Play (Out)":
    atBatResult = "Out"; outsRecorded = 1; isCompleted = true
case "In Play (Hit)":
    atBatResult = "Hit"; isCompleted = true
}
```

## Consecuencias

### Positivas
- **Regla de negocio fiel al béisbol real**: El umpire decide, el sistema registra
- **Simplicidad extrema**: Sin geometría computacional, sin edge cases de zona
- **Auditoría clara**: Cada pitch tiene resultado explícito trazable
- **Testing trivial**: Casos de prueba = valores del enum `pitchResult`
- **Flexibilidad**: Soporta revisiones de umpire (override manual) sin cambiar motor

### Negativas
- **Carga en scorekeeper**: Debe conocer reglas (4 balls = walk, 3 strikes = K)
- **Posible error humano**: Scorekeeper puede equivocarse en resultado
- **Coordenadas "muertas" en lógica**: Solo para UI, no validan consistencia

### Neutrales
- Requiere training de scorekeepers (parte del dominio, no técnico)
- Futuro: podría añadirse validación opcional "coordenada vs resultado" como warning

## Alternativas Consideradas

| Opción | Por qué NO |
|--------|------------|
| **Auto-strike zone (ML/regla geométrica)** | Contradice PRD, complejidad alta, umpire tiene última palabra |
| **Híbrido: coords + override** | Complejidad híbrida, YAGNI, confusión responsabilidad |

## Referencias
- [PRD Module 2](../PRD.md#module-2-baseball-logic-engine-state-engine)
- [TRD 5.1 Constraint 3](../../../docs/TRD.md#51-backend-constraints-go--postgres--redis)
- [scorekeeper.go](../../../backend/internal/core/services/scorekeeper.go)