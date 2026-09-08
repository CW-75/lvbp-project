import { newsArticles } from "@/lib/mock-data";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Clock, Newspaper } from "lucide-react";

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

const categoryColors: Record<string, string> = {
  Resultados: "bg-emerald-500/15 text-emerald-400",
  Fichajes: "bg-blue-500/15 text-blue-400",
  Equipos: "bg-purple-500/15 text-purple-400",
  Temporada: "bg-accent/15 text-accent",
  Destacado: "bg-primary/15 text-primary",
};

export function NewsFeedSection() {
  // Show all non-featured articles
  const articles = newsArticles.filter((a) => !a.isFeatured);

  return (
    <section className="py-12">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        {/* Section header */}
        <div className="mb-6 flex items-center gap-3">
          <Newspaper className="size-5 text-muted-foreground" />
          <h2 className="text-xl font-bold text-foreground">Noticias</h2>
        </div>

        {/* News grid */}
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {articles.map((article) => (
            <Card
              key={article.id}
              className="group cursor-pointer border-white/5 bg-surface transition-all hover:border-white/10 hover:bg-surface-alt"
            >
              {/* Color accent bar */}
              <div className="h-0.5 rounded-t-xl bg-gradient-to-r from-primary/50 to-transparent" />

              <CardHeader className="gap-3 pb-2">
                <div className="flex items-center justify-between">
                  <Badge
                    variant="secondary"
                    className={`border-0 text-[10px] font-semibold ${categoryColors[article.category] ?? ""}`}
                  >
                    {article.category}
                  </Badge>
                  <span className="flex items-center gap-1 text-[10px] text-muted-foreground">
                    <Clock className="size-3" />
                    {formatRelativeDate(article.publishedAt)}
                  </span>
                </div>
                <CardTitle className="text-sm font-bold leading-snug text-foreground transition-colors group-hover:text-primary">
                  {article.title}
                </CardTitle>
              </CardHeader>

              <CardContent>
                <CardDescription className="line-clamp-2 text-xs leading-relaxed text-muted-foreground">
                  {article.excerpt}
                </CardDescription>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
}
