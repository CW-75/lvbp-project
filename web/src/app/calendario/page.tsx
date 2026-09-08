import type { Metadata } from "next";
import { GamesCalendar } from "@/components/calendar/games-calendar";
import { Calendar, ShieldAlert } from "lucide-react";

export const metadata: Metadata = {
  title: "Calendario y Jornadas | LVBP GameCast",
  description:
    "Consulta los horarios, partidos en vivo, resultados y programación completa de la Liga Venezolana de Béisbol Profesional (LVBP).",
};

export default function CalendarioPage() {
  return (
    <div className="min-h-screen bg-[var(--background)] py-8 sm:py-12">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        {/* Page Title Header */}
        <div className="mb-8 flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between border-b border-white/5 pb-6">
          <div>
            <div className="flex items-center gap-2 text-xs font-semibold uppercase tracking-wider text-[var(--secondary)]">
              <Calendar className="size-4" />
              Temporada Regular 2026-2027
            </div>
            <h1 className="mt-1 font-black text-3xl sm:text-4xl tracking-tight text-[var(--foreground)]">
              Jornada y Calendario de Juegos
            </h1>
            <p className="mt-1 text-sm text-[var(--muted-foreground)]">
              Resultados en tiempo real, horarios y programación oficial de los 8 equipos de la LVBP.
            </p>
          </div>
        </div>

        {/* Main Games Calendar Component */}
        <GamesCalendar />
      </div>
    </div>
  );
}
