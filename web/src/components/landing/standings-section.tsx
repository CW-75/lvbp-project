import { standings } from "@/lib/mock-data";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Trophy, ChevronRight } from "lucide-react";

export function StandingsSection() {
  return (
    <section className="py-12">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        {/* Section header */}
        <div className="mb-6 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Trophy className="size-5 text-accent" />
            <h2 className="text-xl font-bold text-foreground">Posiciones</h2>
          </div>
          <button
            type="button"
            className="group/btn flex items-center gap-1 text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
          >
            Tabla completa
            <ChevronRight className="size-4 transition-transform group-hover/btn:translate-x-0.5" />
          </button>
        </div>

        {/* Standings table */}
        <div className="overflow-hidden rounded-xl border border-white/5 bg-surface">
          <Table>
            <TableHeader>
              <TableRow className="border-white/5 hover:bg-transparent">
                <TableHead className="w-12 text-center text-xs font-semibold text-muted-foreground">
                  #
                </TableHead>
                <TableHead className="text-xs font-semibold text-muted-foreground">
                  Equipo
                </TableHead>
                <TableHead className="text-center text-xs font-semibold text-muted-foreground">
                  JJ
                </TableHead>
                <TableHead className="text-center text-xs font-semibold text-muted-foreground">
                  JG
                </TableHead>
                <TableHead className="text-center text-xs font-semibold text-muted-foreground">
                  JP
                </TableHead>
                <TableHead className="text-center text-xs font-semibold text-muted-foreground">
                  PCT
                </TableHead>
                <TableHead className="text-center text-xs font-semibold text-muted-foreground">
                  DIF
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {standings.map((entry, index) => (
                <TableRow
                  key={entry.team.id}
                  className="border-white/5 transition-colors hover:bg-white/[0.02]"
                >
                  <TableCell className="text-center text-sm font-medium text-muted-foreground">
                    {index + 1}
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-2.5">
                      <div
                        className="flex size-6 items-center justify-center rounded text-[10px] font-bold text-white"
                        style={{ backgroundColor: entry.team.primaryColor }}
                      >
                        {entry.team.logoInitials}
                      </div>
                      <div>
                        <p className="text-sm font-semibold text-foreground">
                          {entry.team.shortName}
                        </p>
                        <p className="hidden text-[10px] text-muted-foreground sm:block">
                          {entry.team.city}
                        </p>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell className="text-center text-sm tabular-nums text-foreground">
                    {entry.gamesPlayed}
                  </TableCell>
                  <TableCell className="text-center text-sm font-semibold tabular-nums text-emerald-400">
                    {entry.wins}
                  </TableCell>
                  <TableCell className="text-center text-sm tabular-nums text-red-400">
                    {entry.losses}
                  </TableCell>
                  <TableCell className="text-center text-sm font-bold tabular-nums text-foreground">
                    {entry.pct}
                  </TableCell>
                  <TableCell className="text-center text-sm tabular-nums text-muted-foreground">
                    {entry.diff}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </div>
    </section>
  );
}
