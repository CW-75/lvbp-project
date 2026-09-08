import { Badge } from "@/components/ui/badge";
import { Radio } from "lucide-react";
import type { Game } from "@/lib/mock-data";

export interface LiveBadgeProps {
  inning?: number;
  isTopInning?: boolean;
  className?: string;
}

export interface ScheduledBadgeProps {
  scheduledAt: string;
  className?: string;
}

export interface FinalBadgeProps {
  label?: string;
  className?: string;
}

export interface GameStatusBadgeProps {
  status: Game["status"] | "postponed";
  inning?: number;
  isTopInning?: boolean;
  scheduledAt?: string;
  className?: string;
}

/**
  Badge especifico para partidos EN VIVO
 */
export function LiveBadge({ inning = 1, isTopInning = true, className }: LiveBadgeProps) {
  const inningHalf = isTopInning ? "Alta" : "Baja";
  return (
    <Badge variant="live" className={className}>
      <Radio className="size-3" />
      <span>
        EN VIVO — Inning {inning} ({inningHalf})
      </span>
    </Badge>
  );
}

/**
  Badge especifico para partidos FINALIZADOS
 */
export function FinalBadge({ label = "FINALIZADO", className }: FinalBadgeProps) {
  return (
    <Badge variant="final" className={className}>
      {label}
    </Badge>
  );
}

/**
  Badge especifico para partidos PROGRAMADOS
 */
export function ScheduledBadge({ scheduledAt, className }: ScheduledBadgeProps) {
  const formattedTime = new Date(scheduledAt).toLocaleTimeString("es-VE", {
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
  });

  return (
    <Badge variant="scheduled" className={className}>
      {formattedTime} HRS
    </Badge>
  );
}

/**
  Badge especifico para partidos POSPUESTOS
 */
export function PostponedBadge({ className }: { className?: string }) {
  return (
    <Badge variant="postponed" className={className}>
      POSPUESTO
    </Badge>
  );
}

/**
  Componente polimorfico GameStatusBadge que delega la renderizacion
  segun el estado del juego (SOLID - Single Responsibility / Open-Closed)
 */
export function GameStatusBadge({
  status,
  inning,
  isTopInning,
  scheduledAt,
  className,
}: GameStatusBadgeProps) {
  switch (status) {
    case "live":
      return <LiveBadge inning={inning} isTopInning={isTopInning} className={className} />;
    case "final":
      return <FinalBadge className={className} />;
    case "scheduled":
      return scheduledAt ? (
        <ScheduledBadge scheduledAt={scheduledAt} className={className} />
      ) : (
        <Badge variant="scheduled" className={className}>PROGRAMADO</Badge>
      );
    case "postponed":
      return <PostponedBadge className={className} />;
    default:
      return <Badge variant="outline" className={className}>PROGRAMADO</Badge>;
  }
}
