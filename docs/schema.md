# Database Schema Document (DDL)
**Motor:** PostgreSQL 16
**Proyecto:** Live Baseball GameCast

```sql
-- Enums
CREATE TYPE game_status AS ENUM ('SCHEDULED', 'IN_PROGRESS', 'FINAL', 'POSTPONED', 'SUSPENDED');

-- Equipos
CREATE TABLE teams (
    id VARCHAR(10) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    short_name VARCHAR(20) NOT NULL,
    logo_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Juegos
CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    home_team_id VARCHAR(10) REFERENCES teams(id),
    away_team_id VARCHAR(10) REFERENCES teams(id),
    status game_status DEFAULT 'SCHEDULED',
    home_score INT DEFAULT 0,
    away_score INT DEFAULT 0,
    current_inning INT DEFAULT 1,
    is_top_inning BOOLEAN DEFAULT TRUE,
    start_time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Turnos al bate
CREATE TABLE at_bats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID REFERENCES games(id) ON DELETE CASCADE,
    inning INT NOT NULL,
    is_top_inning BOOLEAN NOT NULL,
    batter_id VARCHAR(50) NOT NULL,
    pitcher_id VARCHAR(50) NOT NULL,
    result VARCHAR(50), 
    runs_scored INT DEFAULT 0,
    outs_recorded INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Lanzamientos (Eventos de alta frecuencia)
CREATE TABLE pitches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    at_bat_id UUID REFERENCES at_bats(id) ON DELETE CASCADE,
    pitch_number INT NOT NULL,
    coordinate_x NUMERIC(4, 2) NOT NULL,
    coordinate_y NUMERIC(4, 2) NOT NULL,
    pitch_result VARCHAR(30) NOT NULL,
    balls_before INT NOT NULL,
    strikes_before INT NOT NULL,
    outs_before INT NOT NULL,
    velocity_mph NUMERIC(4, 1),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Índices recomendados
CREATE INDEX idx_games_status ON games(status);
CREATE INDEX idx_at_bats_game_id ON at_bats(game_id);
CREATE INDEX idx_pitches_at_bat_id ON pitches(at_bat_id);
