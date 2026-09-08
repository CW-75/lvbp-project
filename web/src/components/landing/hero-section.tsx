import { newsArticles } from "@/lib/mock-data";
import { Clock, ArrowRight } from "lucide-react";

function formatRelativeDate(dateStr: string): string {
  const date = new Date(dateStr);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));

  if (diffHours < 1) return "Hace minutos";
  if (diffHours < 24) return `Hace ${diffHours}h`;
  const diffDays = Math.floor(diffHours / 24);
  if (diffDays === 1) return "Ayer";
  return `Hace ${diffDays} días`;
}

export function HeroSection() {
  const featured = newsArticles.find((a) => a.isFeatured) ?? newsArticles[0];

  return (
    <section className="relative overflow-hidden">
      {/* Background gradient that simulates a hero image */}
      <div className="absolute inset-0 bg-gradient-to-br from-primary/20 via-background to-accent/10" />
      <div className="absolute inset-0 bg-[url('/grid.svg')] opacity-[0.03]" />

      <div className="relative mx-auto max-w-7xl px-4 py-16 sm:px-6 sm:py-24 lg:px-8">
        <div className="flex flex-col gap-8 lg:flex-row lg:items-center lg:gap-12">
          {/* Content */}
          <div className="flex flex-1 flex-col gap-6">
            <div className="flex items-center gap-3">
              <span className="inline-flex items-center rounded-full bg-primary/15 px-3 py-1 text-xs font-semibold text-primary">
                {featured.category}
              </span>
              <span className="flex items-center gap-1.5 text-xs text-muted-foreground">
                <Clock className="size-3" />
                {formatRelativeDate(featured.publishedAt)}
              </span>
            </div>

            <h1 className="text-3xl font-extrabold leading-tight tracking-tight text-foreground sm:text-4xl lg:text-5xl">
              {featured.title}
            </h1>

            <p className="max-w-xl text-base leading-relaxed text-muted-foreground sm:text-lg">
              {featured.excerpt}
            </p>

            <div>
              <button
                type="button"
                className="group inline-flex items-center gap-2 rounded-lg bg-primary px-5 py-2.5 text-sm font-semibold text-white transition-all hover:bg-primary/90 hover:shadow-lg hover:shadow-primary/25"
              >
                Leer más
                <ArrowRight className="size-4 transition-transform group-hover:translate-x-0.5" />
              </button>
            </div>
          </div>

          {/* Visual accent — abstract diamond shape */}
          <div className="hidden flex-shrink-0 lg:block">
            <div className="relative size-72">
              <div className="absolute inset-0 rotate-45 rounded-3xl bg-gradient-to-br from-primary/20 to-accent/20 backdrop-blur-sm" />
              <div className="absolute inset-4 rotate-45 rounded-2xl border border-white/10 bg-surface/50 backdrop-blur-md" />
              <div className="absolute inset-0 flex items-center justify-center">
                <div className="text-center">
                  <div className="text-5xl font-black text-primary">⚾</div>
                  <p className="mt-2 text-xs font-semibold tracking-widest text-muted-foreground">
                    EN VIVO
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
