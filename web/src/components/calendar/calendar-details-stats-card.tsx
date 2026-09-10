import type { Game } from "@/lib/mock-data";
import { cn } from "@/lib/utils";

interface CalendarDetailsStatsCardProps {
  game: Game;
  className?: string;
}

export function CalendarDetailsStatsCard({ game, className }: CalendarDetailsStatsCardProps) {
  return (
    <div className={cn("mt-4 border-t border-white/5 pt-3", className)}>
      <div className="grid grid-cols-4 items-center rounded-lg bg-black/20 p-2 text-center text-xs">
        <span className="font-semibold text-muted-foreground">Equipo</span>
        <span className="font-semibold text-muted-foreground">C (Runs)</span>
        <span className="font-semibold text-muted-foreground">H (Hits)</span>
        <span className="font-semibold text-muted-foreground">E (Err)</span>

        <span className="truncate font-bold text-left pl-2 text-foreground">
          {game.awayTeam.shortName}
        </span>
        <span className="font-mono font-bold text-foreground">{game.awayScore}</span>
        <span className="font-mono text-muted-foreground">{game.awayHits}</span>
        <span className="font-mono text-muted-foreground">{game.awayErrors ?? 0}</span>

        <span className="truncate font-bold text-left pl-2 text-foreground">
          {game.homeTeam.shortName}
        </span>
        <span className="font-mono font-bold text-foreground">{game.homeScore}</span>
        <span className="font-mono text-muted-foreground">{game.homeHits}</span>
        <span className="font-mono text-muted-foreground">{game.homeErrors ?? 0}</span>
      </div>
    </div>
  );
}
