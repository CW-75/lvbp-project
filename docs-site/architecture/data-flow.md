# Flujo de Datos - LVBP GameCast

## Flujo Principal: Ingesta de Pitch → Fan Update

```mermaid
sequenceDiagram
    autonumber
    actor Scorekeeper
    participant FE as Frontend (Scorekeeper Panel)
    participant REST as REST Handler (Ingestion)
    participant Engine as Scorekeeper Engine
    participant Repo as Game Repository
    participant DB as PostgreSQL
    participant Bus as Redis Pub/Sub
    participant SSE as SSE Handler
    participant FanFE as Frontend (Fan View)

    Scorekeeper->>FE: Click zona strike + selecciona resultado
    FE->>REST: POST /games/{id}/pitches {pitchResult, coordinateX, coordinateY, velocityMPH}
    REST->>REST: Valida PitchInput (pitchResult en enum)
    REST->>Engine: RecordPitch(gameID, pitch)
    Engine->>Repo: GetGameByID(gameID)
    Repo->>DB: SELECT * FROM games WHERE id = $1
    DB-->>Repo: Game (status, inning, scores, etc.)
    Repo-->>Engine: Game domain
    Engine->>Engine: Valida game.Status == IN_PROGRESS
    Engine->>Repo: GetCurrentAtBat(gameID)
    Repo->>DB: SELECT * FROM at_bats WHERE game_id = $1 AND result IS NULL
    DB-->>Repo: AtBat actual
    Repo-->>Engine: AtBat domain
    Engine->>Repo: GetPitchesForAtBat(atBat.ID)
    Repo->>DB: SELECT * FROM pitches WHERE at_bat_id = $1 ORDER BY pitch_number
    DB-->>Repo: Lista pitches previos
    Repo-->>Engine: []Pitch
    Engine->>Engine: calculateCount(pitches) → balls, strikes
    Engine->>Engine: GetOutsForInning → outs
    Engine->>Engine: Construye Pitch domain (con ballsBefore, strikesBefore, outsBefore)
    Engine->>Repo: CreatePitch(pitch)
    Repo->>DB: INSERT INTO pitches (...) VALUES (...)
    DB-->>Repo: Pitch creado (con ID, CreatedAt)
    Repo-->>Engine: Pitch persistido
    Engine->>Bus: PublishEvent("game:{id}", {type: "PitchThrown", pitch})
    Bus-->>SSE: Entrega mensaje a suscriptores
    SSE->>FanFE: SSE data: {"type":"PitchThrown","pitch":{...}}
    FanFE->>FanFE: Zustand store update → UI re-render
    
    alt AtBat Completado (4 balls, 3 strikes, hit, out)
        Engine->>Engine: Evalúa pitchResult → atBatResult, outsRecorded
        Engine->>Repo: UpdateAtBatResult(atBat con result)
        Repo->>DB: UPDATE at_bats SET result = $1, runs_scored = $2, outs_recorded = $3 WHERE id = $4
        DB-->>Repo: OK
        Repo-->>Engine: AtBat actualizado
        Engine->>Bus: PublishEvent("game:{id}", {type: "AtBatCompleted", atBat})
        Bus-->>SSE: Entrega
        SSE->>FanFE: SSE AtBatCompleted
        
        alt 3 Outs → Cambio de Inning
            Engine->>Engine: Actualiza game.IsTopInning / game.CurrentInning
            Engine->>Repo: UpdateGame(game)
            Repo->>DB: UPDATE games SET is_top_inning = $1, current_inning = $2 WHERE id = $3
            DB-->>Repo: OK
            Engine->>Bus: PublishEvent("game:{id}", {type: "InningChanged", game})
            Bus-->>SSE: Entrega
            SSE->>FanFE: SSE InningChanged
        end
        
        alt Runs Scored
            Engine->>Repo: ScoreRuns(gameID, runs)
            Repo->>DB: UPDATE games SET home_score/away_score = ... WHERE id = $1
            Engine->>Bus: PublishEvent("game:{id}", {type: "ScoreUpdated", game, runs_added})
            Bus-->>SSE: Entrega
            SSE->>FanFE: SSE ScoreUpdated
        end
    end
    
    REST-->>FE: 200 OK PitchResponse {game, atBat, events}
    FE->>Scorekeeper: Confirmación visual
```

## Flujo: Consulta Standings (REST)

```mermaid
sequenceDiagram
    autonumber
    actor Fan
    participant FE as Frontend
    participant REST as Standings Handler
    participant Repo as Standings Repository
    participant DB as PostgreSQL

    Fan->>FE: Navega a /standings
    FE->>REST: GET /standings
    REST->>Repo: GetStandings(ctx)
    Repo->>DB: SELECT * FROM standings_view ORDER BY pct DESC
    DB-->>Repo: Rows (team_id, gp, w, l, pct, gb)
    Repo-->>REST: []TeamStandings domain
    REST-->>FE: 200 OK [TeamStandings...]
    FE->>FE: TanStack Query cache → UI render
```

## Flujo: Consulta Boxscore (REST)

```mermaid
sequenceDiagram
    autonumber
    actor Fan
    participant FE as Frontend
    participant REST as Boxscore Handler
    participant Repo as Game Repository
    participant DB as PostgreSQL

    Fan->>FE: Click juego en calendario
    FE->>REST: GET /games/{id}/boxscore
    REST->>Repo: GetGameByID(gameID)
    Repo->>DB: SELECT * FROM games WHERE id = $1
    DB-->>Repo: Game row
    Repo-->>REST: Game domain (con at_bats, pitches via joins o queries separadas)
    REST-->>FE: 200 OK Game domain
    FE->>FE: Render boxscore completo
```

## Flujo: Autenticación Scorekeeper

```mermaid
sequenceDiagram
    autonumber
    actor Scorekeeper
    participant FE as Frontend
    participant Auth as Auth Handler
    participant Svc as Auth Service
    participant Repo as User Repository
    participant DB as PostgreSQL
    participant Redis as Redis Sessions

    Scorekeeper->>FE: Login form (email, password)
    FE->>Auth: POST /auth/login {email, password}
    Auth->>Svc: Authenticate(email, password)
    Svc->>Repo: GetUserByEmail(email)
    Repo->>DB: SELECT * FROM users WHERE email = $1
    DB-->>Repo: User (con password_hash)
    Repo-->>Svc: User
    Svc->>Svc: bcrypt.CompareHash(password, hash)
    alt Credenciales válidas
        Svc->>Redis: SET session:{token} {userID} EX 24h
        Svc-->>Auth: {token, user}
        Auth-->>FE: 200 OK {token, user}
        FE->>FE: Store token → Authenticated routes
    else Inválidas
        Svc-->>Auth: Error
        Auth-->>FE: 401 Unauthorized
    end
```

## Flujo: Middleware JWT en Requests Protegidos

```mermaid
sequenceDiagram
    autonumber
    participant Client
    participant Middleware as JWT Middleware
    participant Redis as Redis Sessions
    participant Handler as Protected Handler

    Client->>Middleware: Request + Authorization: Bearer <token>
    Middleware->>Middleware: Parse & Validate JWT signature
    Middleware->>Redis: GET session:{token}
    alt Session existe y no expirada
        Redis-->>Middleware: userID
        Middleware->>Handler: Request con userID en context
        Handler-->>Client: Response
    else Session no existe / expirada
        Redis-->>Middleware: nil
        Middleware-->>Client: 401 Unauthorized
    end
```

## Modelo de Eventos SSE (Game Channel)

```mermaid
classDiagram
    class GameEvent {
        +string type
        +json data
    }
    
    class PitchThrown {
        +Pitch pitch
    }
    
    class AtBatStarted {
        +AtBat atBat
    }
    
    class AtBatCompleted {
        +AtBat atBat
    }
    
    class InningChanged {
        +Game game
    }
    
    class ScoreUpdated {
        +Game game
        +int runs_added
    }
    
    class Connected {
        +string status
        +string gameId
    }
    
    GameEvent <|-- PitchThrown
    GameEvent <|-- AtBatStarted
    GameEvent <|-- AtBatCompleted
    GameEvent <|-- InningChanged
    GameEvent <|-- ScoreUpdated
    GameEvent <|-- Connected
```

## Canales Redis

| Canal | Patrón | Productores | Consumidores |
|-------|--------|-------------|--------------|
| `game:{gameId}:events` | Pub/Sub | Scorekeeper Engine | SSE StreamHandler |
| `session:{token}` | Key-Value (TTL) | Auth Service | JWT Middleware |

## Esquema de Base de Datos (Relacional)

```mermaid
erDiagram
    USERS ||--o{ GAMES : ""
    TEAMS ||--o{ GAMES : "home_team"
    TEAMS ||--o{ GAMES : "away_team"
    GAMES ||--o{ AT_BATS : "has"
    AT_BATS ||--o{ PITCHES : "contains"
    
    USERS {
        uuid id PK
        varchar email UK
        varchar password_hash
        timestamptz created_at
    }
    
    TEAMS {
        varchar(10) id PK
        varchar name
        varchar short_name
        text logo_url
        timestamptz created_at
    }
    
    GAMES {
        uuid id PK
        varchar(10) home_team_id FK
        varchar(10) away_team_id FK
        game_status status
        int home_score
        int away_score
        int current_inning
        boolean is_top_inning
        timestamptz start_time
        timestamptz created_at
    }
    
    AT_BATS {
        uuid id PK
        uuid game_id FK
        int inning
        boolean is_top_inning
        varchar batter_id
        varchar pitcher_id
        varchar result
        int runs_scored
        int outs_recorded
        timestamptz created_at
    }
    
    PITCHES {
        uuid id PK
        uuid at_bat_id FK
        int pitch_number
        numeric(4,2) coordinate_x
        numeric(4,2) coordinate_y
        varchar pitch_result
        int balls_before
        int strikes_before
        int outs_before
        numeric(4,1) velocity_mph
        timestamptz created_at
    }
```