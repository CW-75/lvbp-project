# DB Schema Summary (Telegraphic)
**Project:** LBGC (Live Baseball GameCast)
**Engine:** PostgreSQL 16

## Entities & Relations
- `teams` (1) -> (N) `games` (home/away)
- `games` (1) -> (N) `at_bats`
- `at_bats` (1) -> (N) `pitches`

## Key Fields
- `games`: status, scores, current_inning, is_top_inning.
- `at_bats`: batter_id, pitcher_id, result, runs_scored, outs_recorded.
- `pitches`: coord_x, coord_y, pitch_result, count_before (B/S/O), velocity_mph.

*Note: For full DDL including enums and indexes, refer to [schema.md](./schema.md).*
