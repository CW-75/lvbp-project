import { Menu } from "lucide-react";
import Image from "next/image";
import Link from "next/link";

const navLinks = [
  { href: "/", label: "Inicio" },
  { href: "/calendario", label: "Calendario" },
  { href: "/juegos", label: "Juegos" },
  { href: "/posiciones", label: "Posiciones" },
  { href: "/gamecast", label: "GameCast" },
];

export function Header() {
  return (
    <header className="sticky top-0 z-50 border-b border-white/5 bg-background/80 backdrop-blur-xl">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
        {/* Logo */}
        <Link href="/" className="flex h-full py-3 items-center gap-3">
          <Image 
            src="/logo-lvbp-darken.svg" 
            alt="LVBP GameCast Logo" 
            width={120} 
            height={36} 
            priority
            className="h-full w-auto object-contain"
          />
        </Link>

        {/* Desktop nav */}
        <nav className="hidden items-center gap-1 md:flex">
          {navLinks.map((link) => (
            <Link
              key={link.href}
              href={link.href}
              className="rounded-md px-3 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-white/5 hover:text-foreground"
            >
              {link.label}
            </Link>
          ))}
        </nav>

        {/* Mobile menu button */}
        <button
          type="button"
          className="inline-flex size-10 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-white/5 hover:text-foreground md:hidden"
          aria-label="Abrir menú"
        >
          <Menu className="size-5" />
        </button>
      </div>
    </header>
  );
}
