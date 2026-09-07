# Product Requirements Document (PRD)
**Project:** Live Baseball GameCast & Standings Engine

## 1. Use Cases and Actors
*   **Scorekeeper:** Authenticated user responsible for registering game action in real-time from the stadium.
*   **Fan (Spectator):** End-user who consumes the live scoreboard and standings from any device.

## 2. Functional Specification

### Module 1: Scorekeeper Panel (Ingestion)
*   **Main Interaction:** The scorekeeper must be able to click on a visual matrix (strike zone) to define the X, Y coordinates of the pitch.
*   **Pitch Results:** After registering the coordinate, they must select the explicit result dictated by the umpire or the play: `Ball`, `Called Strike`, `Swinging Strike`, `Foul`, `In Play`, `Hit By Pitch`.
*   **Play Results:** If the pitch is "In Play", a submenu must be displayed (Single, Double, Fly out, etc.) and mark the advance or elimination of runners.
*   **Manual Control (Override):** Buttons to manually correct the ball/strike count or runners' position due to a typing error or umpire review.

### Module 2: Baseball Logic Engine (State Engine)
*   The system does not deduce balls/strikes by physical coordinates; it obeys the scorekeeper's input.
*   The backend must calculate automatic transitions:
    *   Reaching 4 balls generates an advance to 1B.
    *   Reaching 3 strikes generates 1 out.
    *   Reaching 3 outs clears the bases, resets the count, and changes the turn (Top/Bottom Inning).

### Module 3: Spectator Interface (GameCast)
*   **Live Scoreboard:** Display the current count (B-S-O), score, inning, and base occupancy (diamond).
*   **Pitch Visualization:** Render the pitches of the current at-bat in an interactive SVG strike zone.
*   **Standings:** Display updated metrics (Games Played, Won, Lost, PCT, Difference) automatically calculated at the end of the games.
