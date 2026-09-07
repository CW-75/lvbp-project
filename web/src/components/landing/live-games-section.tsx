"use client";

import { games } from "@/lib/mock-data";
import type { Game } from "@/lib/mock-data";
import { ChevronRight, Radio } from "lucide-react";

function GameCard({ game }: { game: Game }) {
  const isLive = game.status === "live";
  const isFinal = game.status === "final";

  return (
    <div className="group relative flex min-w-[280px] flex-col gap-3 rounded-xl border border-white/5 bg-[var(--surface)] p-4 transition-all hover:border-white/10 hover:bg-[var(--surface-alt)]">
      {/* Status badge */}
      <div className="flex items-center justify-between">
        {isLive && (
          <span className="inline-flex items-center gap-1.5 rounded-full bg-[var(--primary)]/15 px-2.5 py-0.5 text-xs font-bold text-[var(--primary)]">
            <span className="animate-pulse-live size-1.5 rounded-full bg-[var(--primary)]" />
            EN VIVO
          </span>
        )}
        {isFinal && (
          <span className="rounded-full bg-white/5 px-2.5 py-0.5 text-xs font-medium text-[var(--muted-foreground)]">
            Final
          </span>
        )}
        {game.status === "scheduled" && (
          <span className="rounded-full bg-[var(--accent)]/10 px-2.5 py-0.5 text-xs font-medium text-[var(--accent)]">
            Programado
          </span>
        )}
        {isLive && (
          <span className="text-xs text-[var(--muted-foreground)]">
            {game.isTopInning ? "▲" : "▼"} {game.inning}°
          </span>
        )}
      </div>

      {/* Teams */}
      <div className="flex flex-col gap-2">
        <TeamRow
          initials={game.awayTeam.logoInitials}
          name={game.awayTeam.shortName}
          fullName={game.awayTeam.name}
          color={game.awayTeam.primaryColor}
          score={game.status !== "scheduled" ? game.awayScore : undefined}
          isWinning={game.awayScore > game.homeScore && !isLive}
        />
        <TeamRow
          initials={game.homeTeam.logoInitials}
          name={game.homeTeam.shortName}
          fullName={game.homeTeam.name}
          color={game.homeTeam.primaryColor}
          score={game.status !== "scheduled" ? game.homeScore : undefined}
          isWinning={game.homeScore > game.awayScore && !isLive}
        />
      </div>

      {/* Scheduled time */}
      {game.status === "scheduled" && (
        <p className="text-xs text-[var(--muted-foreground)]">
          {new Date(game.scheduledAt).toLocaleDateString("es-VE", {
            weekday: "short",
            day: "numeric",
            month: "short",
          })}{" "}
          —{" "}
          {new Date(game.scheduledAt).toLocaleTimeString("es-VE", {
            hour: "2-digit",
            minute: "2-digit",
          })}
        </p>
      )}
    </div>
  );
}

function TeamRow({
  initials,
  name,
  fullName,
  color,
  score,
  isWinning,
}: {
  initials: string;
  name: string;
  fullName: string;
  color: string;
  score?: number;
  isWinning: boolean;
}) {
  return (
    <div className="flex items-center justify-between">
      <div className="flex items-center gap-2.5">
        <div
          className="flex size-7 items-center justify-center rounded-md text-xs font-bold text-white"
          style={{ backgroundColor: color }}
        >
          {initials}
        </div>
        <div>
          <p
            className={`text-sm font-semibold ${isWinning ? "text-[var(--foreground)]" : "text-[var(--foreground)]"}`}
          >
            {name}
          </p>
          <p className="text-[10px] text-[var(--muted-foreground)]">{fullName}</p>
        </div>
      </div>
      {score !== undefined && (
        <span
          className={`text-lg font-bold tabular-nums ${isWinning ? "text-[var(--foreground)]" : "text-[var(--muted-foreground)]"}`}
        >
          {score}
        </span>
      )}
    </div>
  );
}

export function LiveGamesSection() {
  const liveGames = games.filter((g) => g.status === "live");
  const scheduledGames = games.filter((g) => g.status === "scheduled");
  const finalGames = games.filter((g) => g.status === "final");

  const orderedGames = [...liveGames, ...scheduledGames, ...finalGames];

  return (
    <section className="py-12">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        {/* Section header */}
        <div className="mb-6 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Radio className="size-5 text-[var(--primary)]" />
            <h2 className="text-xl font-bold text-[var(--foreground)]">Juegos</h2>
            {liveGames.length > 0 && (
              <span className="inline-flex items-center gap-1.5 rounded-full bg-[var(--primary)]/10 px-2 py-0.5 text-xs font-semibold text-[var(--primary)]">
                {liveGames.length} en vivo
              </span>
            )}
          </div>
          <button
            type="button"
            className="group/btn flex items-center gap-1 text-sm font-medium text-[var(--muted-foreground)] transition-colors hover:text-[var(--foreground)]"
          >
            Ver todos
            <ChevronRight className="size-4 transition-transform group-hover/btn:translate-x-0.5" />
          </button>
        </div>

        {/* Scrollable games row */}
        <div className="flex gap-4 overflow-x-auto pb-2 scrollbar-none">
          {orderedGames.map((game) => (
            <GameCard key={game.id} game={game} />
          ))}
        </div>
      </div>
    </section>
  );
}
