import type { Game } from "@/lib/mock-data";
import { cn } from "@/lib/utils";

interface CalendarPitcherInfoProps {
  game: Game;
  className?: string;
}

export function CalendarPitcherInfo({ game, className }: CalendarPitcherInfoProps) {
  const isFinal = game.status === "final";
  const hasPitcherInfo =
    (isFinal && (game.pitcherWin || game.pitcherLoss || game.pitcherSave)) ||
    (!isFinal && (game.pitcherProbableAway || game.pitcherProbableHome));

  if (!hasPitcherInfo) return null;

  return (
    <div className={cn("mt-3 pt-2 text-xs text-muted-foreground flex flex-wrap gap-y-1 gap-x-4", className)}>
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
  );
}
