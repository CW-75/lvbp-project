import Link from "next/link";

const footerLinks = [
  { href: "/calendario", label: "Calendario" },
  { href: "/juegos", label: "Juegos" },
  { href: "/posiciones", label: "Posiciones" },
  { href: "/gamecast", label: "GameCast" },
];

export function Footer() {
  return (
    <footer className="border-t border-white/5 bg-surface">
      <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        <div className="flex flex-col items-center gap-8 md:flex-row md:justify-between">
          {/* Brand */}
          <div className="flex flex-col items-center gap-2 md:items-start">
            <div className="flex items-center gap-3">
              <div className="flex size-8 items-center justify-center rounded-lg bg-primary font-bold text-white text-xs">
                LV
              </div>
              <span className="text-sm font-bold text-foreground">
                LVBP GameCast
              </span>
            </div>
            <p className="text-xs text-muted-foreground">
              Seguimiento en vivo del béisbol profesional venezolano
            </p>
          </div>

          {/* Links */}
          <nav className="flex items-center gap-6">
            {footerLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className="text-sm text-muted-foreground transition-colors hover:text-foreground"
              >
                {link.label}
              </Link>
            ))}
          </nav>
        </div>

        <div className="mt-8 border-t border-white/5 pt-6 text-center">
          <p className="text-xs text-muted-foreground">
            © {new Date().getFullYear()} LVBP GameCast. Datos con fines demostrativos.
          </p>
        </div>
      </div>
    </footer>
  );
}
