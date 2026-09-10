import type { Team } from "@/lib/mock-data";
import { CalendarTeamInfo } from "@/components/calendar/calendar-team-info";

interface CalendarHomeTeamProps {
  team: Team;
  className?: string;
}

export function CalendarHomeTeam({ team, className }: CalendarHomeTeamProps) {
  return <CalendarTeamInfo team={team} align="home" className={className} />;
}
