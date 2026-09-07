import { HeroSection } from "@/components/landing/hero-section";
import { LiveGamesSection } from "@/components/landing/live-games-section";
import { StandingsSection } from "@/components/landing/standings-section";
import { NewsFeedSection } from "@/components/landing/news-feed-section";

export default function HomePage() {
  return (
    <>
      <HeroSection />
      <LiveGamesSection />
      <StandingsSection />
      <NewsFeedSection />
    </>
  );
}
