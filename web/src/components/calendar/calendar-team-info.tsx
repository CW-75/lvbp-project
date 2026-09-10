import { TeamLogo } from "@/components/ui/team-logo";
import type { Team } from "@/lib/mock-data";
import { cn } from "@/lib/utils";

interface CalendarTeamInfoProps {
  team: Team;
  align?: "away" | "home";
  className?: string;
}

export function CalendarTeamInfo({
  team,
  align = "away",
  className,
}: CalendarTeamInfoProps) {
  const isHome = align === "home";

  return (
    <div
      className={cn(
        "col-span-5 flex items-center gap-3 justify-center",
        isHome ? "sm:justify-end text-right" : "sm:justify-start",
        className
      )}
    >
      <div className={cn("min-w-0 hidden sm:block", isHome && "order-2 sm:order-1")}>
        <h4 className="truncate font-bold sm:text-md text-foreground">
          {team.name}
        </h4>
        <p className="text-xs text-muted-foreground">{team.city}</p>
      </div>
      <TeamLogo
        team={team}
        size="lg"
        className={cn(isHome && "order-1 sm:order-2")}
      />
    </div>
  );
}
