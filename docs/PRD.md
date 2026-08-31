# Product Requirements Document (PRD)
**Proyecto:** Live Baseball GameCast & Standings Engine

## 1. Casos de Uso y Actores
*   **Anotador (Scorekeeper):** Usuario autenticado responsable de registrar la acción del juego en tiempo real desde el estadio.
*   **Espectador (Fan):** Usuario final que consume la pizarra en vivo y la tabla de posiciones desde cualquier dispositivo.

## 2. Especificación Funcional

### Módulo 1: Panel del Anotador (Ingesta)
*   **Interacción Principal:** El anotador debe poder hacer clic en una matriz visual (zona de strike) para definir las coordenadas X, Y del pitcheo.
*   **Resultados de Lanzamiento:** Tras registrar la coordenada, debe seleccionar el resultado explícito dictado por el umpire o la jugada: `Bola`, `Strike Cantado`, `Strike Abanicado`, `Foul`, `En Juego`, `Pelotazo`.
*   **Resultados de Jugada:** Si el pitcheo está "En Juego", se debe desplegar un submenú (Sencillo, Doble, Out elevado, etc.) y marcar el avance o eliminación de corredores.
*   **Control Manual (Override):** Botones para corregir manualmente el conteo de bolas/strikes o la posición de corredores ante un error de digitación o revisión arbitral.

### Módulo 2: Motor Lógico de Béisbol (State Engine)
*   El sistema no deduce bolas/strikes por coordenadas físicas; obedece la entrada del anotador.
*   El backend debe calcular transiciones automáticas:
    *   Llegar a 4 bolas genera avance a 1B.
    *   Llegar a 3 strikes genera 1 out.
    *   Llegar a 3 outs limpia las bases, resetea el conteo y cambia el turno (Top/Bottom Inning).

### Módulo 3: Interfaz del Espectador (GameCast)
*   **Pizarra en Vivo:** Mostrar el conteo actual (B-S-O), marcador, inning y ocupación de bases (diamante).
*   **Visualización de Pitcheos:** Renderizar los lanzamientos del turno actual en una zona de strike SVG interactiva.
*   **Tabla de Posiciones (Standings):** Mostrar métricas actualizadas (Juegos Jugados, Ganados, Perdidos, PCT, Diferencia) calculadas automáticamente al finalizar los encuentros.
