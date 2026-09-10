"use client";

import { CalendarDetailsStatsCard } from "@/components/calendar/calendar-details-stats-card";
import { CalendarPitcherInfo } from "@/components/calendar/calendar-pitcher-info";
import { CalendarTeamInfo } from "@/components/calendar/calendar-team-info";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { GameStatusBadge } from "@/components/ui/game-status-badge";
import type { Game } from "@/lib/mock-data";
import { games, teams } from "@/lib/mock-data";
import {
  Calendar as CalendarIcon,
  ChevronLeft,
  ChevronRight,
  Filter,
  MapPin,
  Trophy,
} from "lucide-react";
import { useState } from "react";

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
      <div className="flex items-center justify-between border-b border-white/5 bg-white/2 px-4 py-3 text-xs">
        <div className="flex items-center gap-2">
          <GameStatusBadge
            status={game.status}
            inning={game.inning}
            isTopInning={game.isTopInning}
            scheduledAt={game.scheduledAt}
          />
        </div>

        {game.stadium && (
          <span className="flex items-center gap-1 text-muted-foreground truncate max-w-50" title={game.stadium}>
            <MapPin className="size-3 shrink-0 text-primary" />
            {game.stadium.split(",")[0]}
          </span>
        )}
      </div>

      <CardContent className="p-4 sm:p-5">
        {/* Teams Matchup Grid */}
        <div className="grid grid-cols-12 items-center gap-2 sm:gap-4">
          {/* Away Team */}
          <CalendarTeamInfo team={game.awayTeam} align="away" />

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
          <CalendarTeamInfo team={game.homeTeam} align="home" />
        </div>

        {/* Detailed Stats (R - H - E) */}
        {(isLive || isFinal) && game.awayHits !== undefined && (
          <CalendarDetailsStatsCard game={game} />
        )}

        {/* Pitchers / Decision Info */}
        <CalendarPitcherInfo game={game} />
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
          <div className="flex items-center gap-1.5 overflow-x-auto scrollbar-none pb-1 md:pb-0">
            <button className="rounded-md pa-2 bg-white/5 text-muted-foreground hover:bg-white/10 hover:text-white cursor-pointer p-2">
              <ChevronLeft className="size-4 text-muted-foreground" />
            </button>
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
            <button className="rounded-md pa-2 bg-white/5 text-muted-foreground hover:bg-white/10 hover:text-white cursor-pointer p-2">
              <ChevronRight className="size-4 text-muted-foreground" />
            </button>
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
