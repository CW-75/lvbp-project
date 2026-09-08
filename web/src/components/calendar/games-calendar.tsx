"use client";

import { useState } from "react";
import { games, teams } from "@/lib/mock-data";
import type { Game, Team } from "@/lib/mock-data";
import { GameStatusBadge } from "@/components/ui/game-status-badge";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  Calendar as CalendarIcon,
  ChevronLeft,
  ChevronRight,
  Filter,
  MapPin,
  Trophy,
} from "lucide-react";

// Fechas disponibles en los datos mock
const AVAILABLE_DATES = [
  { key: "2026-09-06", label: "Dom 6 Sep", short: "Ayer" },
  { key: "2026-09-07", label: "Lun 7 Sep", short: "Hoy", isToday: true },
  { key: "2026-09-08", label: "Mar 8 Sep", short: "Mañana" },
  { key: "2026-09-09", label: "Mié 9 Sep", short: "Próximo" },
];

function GameDetailCard({ game }: { game: Game }) {
  const isLive = game.status === "live";
  const isFinal = game.status === "final";

  return (
    <Card className="overflow-hidden border border-white/10 bg-surface transition-all hover:border-white/20 hover:shadow-lg">
      {/* Top Banner / Status */}
      <div className="flex items-center justify-between border-b border-white/5 bg-white/[0.02] px-4 py-3 text-xs">
        <div className="flex items-center gap-2">
          <GameStatusBadge
            status={game.status}
            inning={game.inning}
            isTopInning={game.isTopInning}
            scheduledAt={game.scheduledAt}
          />
        </div>

        {game.stadium && (
          <span className="flex items-center gap-1 text-muted-foreground truncate max-w-[200px]" title={game.stadium}>
            <MapPin className="size-3 shrink-0 text-primary" />
            {game.stadium.split(",")[0]}
          </span>
        )}
      </div>

      <CardContent className="p-4 sm:p-5">
        {/* Teams Matchup Grid */}
        <div className="grid grid-cols-12 items-center gap-2 sm:gap-4">
          {/* Away Team */}
          <div className="col-span-5 flex items-center justify-between sm:justify-start gap-3">
            <div
              className="flex size-10 sm:size-12 shrink-0 items-center justify-center rounded-xl font-bold text-white text-sm shadow-md"
              style={{ backgroundColor: game.awayTeam.primaryColor }}
            >
              {game.awayTeam.logoInitials}
            </div>
            <div className="min-w-0">
              <h4 className="truncate font-bold text-sm sm:text-base text-foreground">
                {game.awayTeam.name}
              </h4>
              <p className="text-xs text-muted-foreground">{game.awayTeam.city}</p>
            </div>
          </div>

          {/* Scores or VS */}
          <div className="col-span-2 text-center">
            {isLive || isFinal ? (
              <div className="flex items-center justify-center gap-2 font-mono font-black text-xl sm:text-2xl text-foreground">
                <span className={game.awayScore > game.homeScore ? "text-accent font-extrabold" : ""}>
                  {game.awayScore}
                </span>
                <span className="text-muted-foreground font-normal text-sm">-</span>
                <span className={game.homeScore > game.awayScore ? "text-accent font-extrabold" : ""}>
                  {game.homeScore}
                </span>
              </div>
            ) : (
              <span className="rounded-full bg-white/5 px-2.5 py-1 font-bold text-xs text-muted-foreground">
                VS
              </span>
            )}
          </div>

          {/* Home Team */}
          <div className="col-span-5 flex items-center justify-between sm:justify-end gap-3 text-right">
            <div className="min-w-0 order-2 sm:order-1">
              <h4 className="truncate font-bold text-sm sm:text-base text-foreground">
                {game.homeTeam.name}
              </h4>
              <p className="text-xs text-muted-foreground">{game.homeTeam.city}</p>
            </div>
            <div
              className="flex size-10 sm:size-12 shrink-0 items-center justify-center rounded-xl font-bold text-white text-sm shadow-md order-1 sm:order-2"
              style={{ backgroundColor: game.homeTeam.primaryColor }}
            >
              {game.homeTeam.logoInitials}
            </div>
          </div>
        </div>

        {/* Detailed Stats (R - H - E) */}
        {(isLive || isFinal) && game.awayHits !== undefined && (
          <div className="mt-4 border-t border-white/5 pt-3">
            <div className="grid grid-cols-4 items-center rounded-lg bg-black/20 p-2 text-center text-xs">
              <span className="font-semibold text-muted-foreground">Equipo</span>
              <span className="font-semibold text-muted-foreground">C (Runs)</span>
              <span className="font-semibold text-muted-foreground">H (Hits)</span>
              <span className="font-semibold text-muted-foreground">E (Err)</span>

              <span className="truncate font-bold text-left pl-2 text-foreground">{game.awayTeam.shortName}</span>
              <span className="font-mono font-bold text-foreground">{game.awayScore}</span>
              <span className="font-mono text-muted-foreground">{game.awayHits}</span>
              <span className="font-mono text-muted-foreground">{game.awayErrors ?? 0}</span>

              <span className="truncate font-bold text-left pl-2 text-foreground">{game.homeTeam.shortName}</span>
              <span className="font-mono font-bold text-foreground">{game.homeScore}</span>
              <span className="font-mono text-muted-foreground">{game.homeHits}</span>
              <span className="font-mono text-muted-foreground">{game.homeErrors ?? 0}</span>
            </div>
          </div>
        )}

        {/* Pitchers / Decision Info */}
        <div className="mt-3 pt-2 text-xs text-muted-foreground flex flex-wrap gap-y-1 gap-x-4">
          {isFinal && game.pitcherWin && (
            <span>
              <strong className="text-emerald-400">G:</strong> {game.pitcherWin}
            </span>
          )}
          {isFinal && game.pitcherLoss && (
            <span>
              <strong className="text-rose-400">P:</strong> {game.pitcherLoss}
            </span>
          )}
          {isFinal && game.pitcherSave && (
            <span>
              <strong className="text-accent">S:</strong> {game.pitcherSave}
            </span>
          )}
          {!isFinal && game.pitcherProbableAway && (
            <span>
              <strong>Lanzador {game.awayTeam.shortName}:</strong> {game.pitcherProbableAway}
            </span>
          )}
          {!isFinal && game.pitcherProbableHome && (
            <span>
              <strong>Lanzador {game.homeTeam.shortName}:</strong> {game.pitcherProbableHome}
            </span>
          )}
        </div>
      </CardContent>
    </Card>
  );
}

export function GamesCalendar() {
  const [selectedDate, setSelectedDate] = useState<string>("2026-09-07");
  const [statusFilter, setStatusFilter] = useState<"all" | "live" | "final" | "scheduled">("all");
  const [selectedTeamId, setSelectedTeamId] = useState<string>("all");

  // Filtrar juegos por fecha, estado y equipo
  const filteredGames = games.filter((game) => {
    const gameDate = game.scheduledAt.split("T")[0];
    if (gameDate !== selectedDate) return false;

    if (statusFilter !== "all" && game.status !== statusFilter) return false;

    if (selectedTeamId !== "all") {
      if (game.homeTeam.id !== selectedTeamId && game.awayTeam.id !== selectedTeamId) {
        return false;
      }
    }

    return true;
  });

  return (
    <div className="space-y-8">
      {/* Date Selector Header */}
      <div className="rounded-2xl border border-white/10 bg-surface p-4 shadow-xl">
        <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div className="flex items-center gap-2">
            <CalendarIcon className="size-5 text-primary" />
            <h3 className="font-bold text-lg text-foreground">Jornada Regular LVBP</h3>
          </div>

          {/* Date Selector Pills */}
          <div className="flex items-center gap-1.5 overflow-x-auto pb-1 md:pb-0">
            {AVAILABLE_DATES.map((dateObj) => {
              const isSelected = selectedDate === dateObj.key;
              return (
                <button
                  key={dateObj.key}
                  onClick={() => setSelectedDate(dateObj.key)}
                  className={`flex flex-col items-center justify-center rounded-xl px-4 py-2 text-xs font-semibold transition-all shrink-0 ${
                    isSelected
                      ? "bg-primary text-white shadow-lg shadow-primary/30 scale-105"
                      : "bg-white/5 text-muted-foreground hover:bg-white/10 hover:text-white"
                  }`}
                >
                  <span className="text-[10px] uppercase opacity-80">{dateObj.short}</span>
                  <span>{dateObj.label}</span>
                </button>
              );
            })}
          </div>
        </div>

        {/* Filters Row */}
        <div className="mt-4 flex flex-wrap items-center justify-between gap-3 border-t border-white/5 pt-4">
          {/* Status Filter Tabs */}
          <div className="flex items-center gap-1.5 rounded-lg bg-black/30 p-1 border border-white/5">
            {[
              { id: "all", label: "Todos los juegos" },
              { id: "live", label: "En Vivo" },
              { id: "final", label: "Finalizados" },
              { id: "scheduled", label: "Programados" },
            ].map((tab) => (
              <button
                key={tab.id}
                onClick={() => setStatusFilter(tab.id as any)}
                className={`rounded-md px-3 py-1.5 text-xs font-medium transition-colors ${
                  statusFilter === tab.id
                    ? "bg-white/15 text-white font-bold"
                    : "text-muted-foreground hover:text-white"
                }`}
              >
                {tab.label}
              </button>
            ))}
          </div>

          {/* Team Filter Dropdown */}
          <div className="flex items-center gap-2">
            <Filter className="size-4 text-muted-foreground" />
            <select
              value={selectedTeamId}
              onChange={(e) => setSelectedTeamId(e.target.value)}
              className="rounded-lg border border-white/10 bg-black/40 px-3 py-1.5 text-xs font-medium text-foreground focus:border-primary focus:outline-none"
            >
              <option value="all">Todos los equipos</option>
              {teams.map((t) => (
                <option key={t.id} value={t.id}>
                  {t.name}
                </option>
              ))}
            </select>
          </div>
        </div>
      </div>

      {/* Games List Grid */}
      {filteredGames.length > 0 ? (
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
          {filteredGames.map((game) => (
            <GameDetailCard key={game.id} game={game} />
          ))}
        </div>
      ) : (
        <Card className="border border-white/10 bg-surface p-12 text-center">
          <div className="mx-auto flex size-12 items-center justify-center rounded-full bg-white/5">
            <Trophy className="size-6 text-muted-foreground" />
          </div>
          <h3 className="mt-4 font-bold text-lg text-foreground">
            No hay partidos registrados
          </h3>
          <p className="mt-1 text-sm text-muted-foreground">
            No se encontraron juegos para los filtros seleccionados en esta fecha.
          </p>
          <Button
            variant="outline"
            className="mt-4"
            onClick={() => {
              setStatusFilter("all");
              setSelectedTeamId("all");
            }}
          >
            Restablecer Filtros
          </Button>
        </Card>
      )}
    </div>
  );
}
