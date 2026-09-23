-- Test User (password: mySuperSecretPassword123)
INSERT INTO users (email, password_hash) 
VALUES ('test@example.com', '$2a$10$ujzR0u6YtpbChC.hrG85s.zdNU1fItcPmoNIXYFraC6gjfgkLeuq.') 
ON CONFLICT (email) DO NOTHING;

-- Insert LVBP Teams
INSERT INTO teams (id, name, short_name) VALUES 
('CAR', 'Leones del Caracas', 'Caracas'),
('MAG', 'Navegantes del Magallanes', 'Magallanes'),
('TIG', 'Tigres de Aragua', 'Aragua'),
('TIB', 'Tiburones de La Guaira', 'La Guaira'),
('AGU', 'Águilas del Zulia', 'Zulia'),
('CARIB', 'Caribes de Anzoátegui', 'Anzoátegui'),
('BRA', 'Bravos de Margarita', 'Margarita'),
('LAR', 'Cardenales de Lara', 'Lara')
ON CONFLICT (id) DO NOTHING;

-- Insert some dummy games (One in-progress, one scheduled)
INSERT INTO games (id, home_team_id, away_team_id, status, home_score, away_score, current_inning, is_top_inning, start_time) 
VALUES 
('11111111-1111-1111-1111-111111111111', 'CAR', 'MAG', 'IN_PROGRESS', 2, 1, 3, false, NOW() - INTERVAL '1 hour'),
('22222222-2222-2222-2222-222222222222', 'TIG', 'TIB', 'SCHEDULED', 0, 0, 1, true, NOW() + INTERVAL '2 hours')
ON CONFLICT (id) DO NOTHING;

-- Insert an At-Bat for the in-progress game
INSERT INTO at_bats (id, game_id, inning, is_top_inning, batter_id, pitcher_id, result) 
VALUES 
('33333333-3333-3333-3333-333333333333', '11111111-1111-1111-1111-111111111111', 3, false, 'Altuve', 'Hernandez', NULL)
ON CONFLICT (id) DO NOTHING;

-- Insert some Pitches for that At-Bat
INSERT INTO pitches (at_bat_id, pitch_number, coordinate_x, coordinate_y, pitch_result, balls_before, strikes_before, outs_before, velocity_mph) 
VALUES 
('33333333-3333-3333-3333-333333333333', 1, 0.5, 0.5, 'Called Strike', 0, 0, 1, 95.5),
('33333333-3333-3333-3333-333333333333', 2, -1.0, 1.2, 'Ball', 0, 1, 1, 84.0)
ON CONFLICT (id) DO NOTHING;
