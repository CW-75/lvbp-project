import Link from "next/link";
import { Menu } from "lucide-react";

const navLinks = [
  { href: "/", label: "Inicio" },
  { href: "/calendario", label: "Calendario" },
  { href: "/juegos", label: "Juegos" },
  { href: "/posiciones", label: "Posiciones" },
  { href: "/gamecast", label: "GameCast" },
];

export function Header() {
  return (
    <header className="sticky top-0 z-50 border-b border-white/5 bg-[var(--background)]/80 backdrop-blur-xl">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
        {/* Logo */}
        <Link href="/" className="flex items-center gap-3">
          <div className="flex size-9 items-center justify-center rounded-lg bg-[var(--primary)] font-bold text-white text-sm">
            LV
          </div>
          <span className="text-lg font-bold tracking-tight text-[var(--foreground)]">
            LVBP <span className="font-normal text-[var(--muted-foreground)]">GameCast</span>
          </span>
        </Link>

        {/* Desktop nav */}
        <nav className="hidden items-center gap-1 md:flex">
          {navLinks.map((link) => (
            <Link
              key={link.href}
              href={link.href}
              className="rounded-md px-3 py-2 text-sm font-medium text-[var(--muted-foreground)] transition-colors hover:bg-white/5 hover:text-[var(--foreground)]"
            >
              {link.label}
            </Link>
          ))}
        </nav>

        {/* Mobile menu button */}
        <button
          type="button"
          className="inline-flex size-10 items-center justify-center rounded-md text-[var(--muted-foreground)] transition-colors hover:bg-white/5 hover:text-[var(--foreground)] md:hidden"
          aria-label="Abrir menú"
        >
          <Menu className="size-5" />
        </button>
      </div>
    </header>
  );
}
