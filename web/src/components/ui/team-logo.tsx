"use client";

import type { Team } from "@/lib/mock-data";
import { cn } from "@/lib/utils";
import { useState } from "react";

interface TeamLogoProps {
  team: Team;
  size?: "sm" | "md" | "lg";
  className?: string;
}

export function TeamLogo({ team, size = "md", className }: TeamLogoProps) {
  const [hasError, setHasError] = useState(false);

  return (
    <div
      className={cn(
        "flex size-16 sm:size-12 shrink-0 items-center justify-center rounded-xl font-bold text-white text-sm shadow-md overflow-hidden transition-all",
        size === "sm" && "size-8 text-xs rounded-lg",
        size === "lg" && "size-20 sm:size-16 text-base rounded-2xl",
        className
      )}
      style={{ backgroundColor: team.primaryColor }}
    >
      {team.logoUrl ? (
        <img
          src={team.logoUrl}
          alt={`Logo ${team.name}`}
          className="size-full object-contain p-0.75"
          onError={() => setHasError(true)}
        />
      ) : (
        <span>{team.logoInitials}</span>
      )}
    </div>
  );
}
