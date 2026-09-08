// Datos mock para la landing page LVBP
// Estos datos serán reemplazados por llamadas REST al backend en producción

export interface Team {
  id: string;
  name: string;
  shortName: string;
  city: string;
  primaryColor: string;
  logoInitials: string;
}

export interface Game {
  id: string;
  homeTeam: Team;
  awayTeam: Team;
  homeScore: number;
  awayScore: number;
  inning: number;
  isTopInning: boolean;
  status: "live" | "scheduled" | "final";
  scheduledAt: string;
}

export interface NewsArticle {
  id: string;
  title: string;
  excerpt: string;
  category: "Resultados" | "Fichajes" | "Equipos" | "Temporada" | "Destacado";
  imageUrl: string;
  publishedAt: string;
  isFeatured: boolean;
}

export interface StandingsEntry {
  team: Team;
  gamesPlayed: number;
  wins: number;
  losses: number;
  pct: string;
  diff: string;
}

// --- Equipos LVBP ---
export const teams: Team[] = [
  {
    id: "ldc",
    name: "Leones del Caracas",
    shortName: "LEO",
    city: "Caracas",
    primaryColor: "#C8102E",
    logoInitials: "LC",
  },
  {
    id: "ndm",
    name: "Navegantes del Magallanes",
    shortName: "MAG",
    city: "Valencia",
    primaryColor: "#F57C00",
    logoInitials: "NM",
  },
  {
    id: "cda",
    name: "Caribes de Anzoátegui",
    shortName: "CAR",
    city: "Puerto La Cruz",
    primaryColor: "#D32F2F",
    logoInitials: "CA",
  },
  {
    id: "adz",
    name: "Águilas del Zulia",
    shortName: "AGU",
    city: "Maracaibo",
    primaryColor: "#1565C0",
    logoInitials: "AZ",
  },
  {
    id: "tda",
    name: "Tigres de Aragua",
    shortName: "TIG",
    city: "Maracay",
    primaryColor: "#FF8F00",
    logoInitials: "TA",
  },
  {
    id: "tdlg",
    name: "Tiburones de La Guaira",
    shortName: "TIB",
    city: "La Guaira",
    primaryColor: "#0D47A1",
    logoInitials: "TG",
  },
  {
    id: "bdm",
    name: "Bravos de Margarita",
    shortName: "BRA",
    city: "Porlamar",
    primaryColor: "#2E7D32",
    logoInitials: "BM",
  },
  {
    id: "cdl",
    name: "Cardenales de Lara",
    shortName: "CDL",
    city: "Barquisimeto",
    primaryColor: "#B71C1C",
    logoInitials: "CL",
  },
];

// --- Helper para buscar equipos ---
function getTeam(id: string): Team {
  return teams.find((t) => t.id === id)!;
}

// --- Juegos ---
export const games: Game[] = [
  {
    id: "g1",
    homeTeam: getTeam("ldc"),
    awayTeam: getTeam("ndm"),
    homeScore: 4,
    awayScore: 3,
    inning: 7,
    isTopInning: false,
    status: "live",
    scheduledAt: "2026-09-07T19:00:00-04:00",
  },
  {
    id: "g2",
    homeTeam: getTeam("cda"),
    awayTeam: getTeam("adz"),
    homeScore: 2,
    awayScore: 5,
    inning: 5,
    isTopInning: true,
    status: "live",
    scheduledAt: "2026-09-07T19:00:00-04:00",
  },
  {
    id: "g3",
    homeTeam: getTeam("tda"),
    awayTeam: getTeam("tdlg"),
    homeScore: 0,
    awayScore: 0,
    inning: 0,
    isTopInning: true,
    status: "scheduled",
    scheduledAt: "2026-09-08T18:30:00-04:00",
  },
  {
    id: "g4",
    homeTeam: getTeam("bdm"),
    awayTeam: getTeam("cdl"),
    homeScore: 0,
    awayScore: 0,
    inning: 0,
    isTopInning: true,
    status: "scheduled",
    scheduledAt: "2026-09-08T19:00:00-04:00",
  },
  {
    id: "g5",
    homeTeam: getTeam("adz"),
    awayTeam: getTeam("ldc"),
    homeScore: 6,
    awayScore: 8,
    inning: 9,
    isTopInning: false,
    status: "final",
    scheduledAt: "2026-09-06T19:00:00-04:00",
  },
];

// --- Standings ---
export const standings: StandingsEntry[] = [
  { team: getTeam("ldc"), gamesPlayed: 32, wins: 21, losses: 11, pct: ".656", diff: "--" },
  { team: getTeam("cda"), gamesPlayed: 32, wins: 19, losses: 13, pct: ".594", diff: "2.0" },
  { team: getTeam("ndm"), gamesPlayed: 31, wins: 18, losses: 13, pct: ".581", diff: "2.5" },
  { team: getTeam("adz"), gamesPlayed: 32, wins: 17, losses: 15, pct: ".531", diff: "4.0" },
  { team: getTeam("tdlg"), gamesPlayed: 31, wins: 15, losses: 16, pct: ".484", diff: "5.5" },
  { team: getTeam("tda"), gamesPlayed: 32, wins: 14, losses: 18, pct: ".438", diff: "7.0" },
  { team: getTeam("cdl"), gamesPlayed: 31, wins: 12, losses: 19, pct: ".387", diff: "8.5" },
  { team: getTeam("bdm"), gamesPlayed: 32, wins: 10, losses: 22, pct: ".313", diff: "11.0" },
];

// --- Noticias ---
export const newsArticles: NewsArticle[] = [
  {
    id: "n1",
    title: "Leones aseguran el liderato con victoria sobre Magallanes en clásico capitalino",
    excerpt:
      "Con jonrón de tres carreras en la séptima entrada, los melenudos sellaron una victoria crucial que los consolida en la cima de la tabla.",
    category: "Resultados",
    imageUrl: "/placeholder-news-1.jpg",
    publishedAt: "2026-09-07T22:00:00-04:00",
    isFeatured: true,
  },
  {
    id: "n2",
    title: "Caribes anuncian la incorporación del lanzador dominicano Ramírez para el round robin",
    excerpt:
      "La directiva de los aborígenes confirmó el refuerzo importado que buscará fortalecer la rotación para la fase eliminatoria.",
    category: "Fichajes",
    imageUrl: "/placeholder-news-2.jpg",
    publishedAt: "2026-09-07T14:30:00-04:00",
    isFeatured: false,
  },
  {
    id: "n3",
    title: "Águilas del Zulia regresan a Maracaibo tras gira exitosa por el oriente",
    excerpt:
      "Los rapaces consiguieron 4 victorias en 5 juegos durante su recorrido por Puerto La Cruz y Porlamar.",
    category: "Equipos",
    imageUrl: "/placeholder-news-3.jpg",
    publishedAt: "2026-09-06T18:00:00-04:00",
    isFeatured: false,
  },
  {
    id: "n4",
    title: "La LVBP confirma calendario del Round Robin 2026-2027",
    excerpt:
      "La liga dio a conocer las fechas y sedes para la fase semifinal, que arrancará el próximo 28 de diciembre.",
    category: "Temporada",
    imageUrl: "/placeholder-news-4.jpg",
    publishedAt: "2026-09-05T10:00:00-04:00",
    isFeatured: false,
  },
];
