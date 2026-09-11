// Datos mock para la landing page LVBP
// Estos datos serán reemplazados por llamadas REST al backend en producción

export interface Team {
  id: string;
  name: string;
  shortName: string;
  city: string;
  primaryColor: string;
  logoInitials: string;
  logoUrl?: string;
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
  stadium?: string;
  homeHits?: number;
  awayHits?: number;
  homeErrors?: number;
  awayErrors?: number;
  pitcherWin?: string;
  pitcherLoss?: string;
  pitcherSave?: string;
  pitcherProbableHome?: string;
  pitcherProbableAway?: string;
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
    primaryColor: "#0B1A32",
    logoInitials: "LC",
    logoUrl: "/logo-leones.svg",
  },
  {
    id: "ndm",
    name: "Navegantes del Magallanes",
    shortName: "MAG",
    city: "Valencia",
    primaryColor: "#FFFFFF",
    logoUrl: "/logo-magallanes.svg",
    logoInitials: "NM",
  },
  {
    id: "tdlg",
    name: "Tiburones de La Guaira",
    shortName: "LAG",
    city: "La Guaira",
    primaryColor: "#0c2037ff",
    logoUrl: "/logo-tiburones.svg",
    logoInitials: "TG",
  },
  {
    id: "cdl",
    name: "Cardenales de Lara",
    shortName: "LAR",
    city: "Barquisimeto",
    primaryColor: "#af2005ff",
    logoUrl: "/logo-cardenales.svg",
    logoInitials: "CL",
  },
  {
    id: "adz",
    name: "Águilas del Zulia",
    shortName: "ZUL",
    city: "Maracaibo",
    primaryColor: "#cc6001ff",
    logoInitials: "AZ",
    logoUrl: "/logo-aguilas.svg",
  },
  {
    id: "cda",
    name: "Caribes de Anzoátegui",
    shortName: "ANZ",
    city: "Puerto La Cruz",
    primaryColor: "#f4c794ff",
    logoInitials: "CA",
    logoUrl: "/logo-caribes.svg",
  },
  {
    id: "bdm",
    name: "Bravos de Margarita",
    shortName: "MAR",
    city: "Porlamar",
    primaryColor: "#2c626cff",
    logoUrl: "/logo-bravos.svg",
    logoInitials: "BM",
  },
  {
    id: "tda",
    name: "Tigres de Aragua",
    shortName: "ARA",
    city: "Maracay",
    primaryColor: "#00205bff",
    logoInitials: "TA",
    logoUrl: "/logo-tigres.svg",
  },
];

export function getTeam(id: string): Team {
  return teams.find((t) => t.id === id) || teams[0];
}

// --- Juegos LVBP (Jornada y Calendario) ---
export const games: Game[] = [
  // 7 de Septiembre (Hoy - Juegos en vivo y finalizados)
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
    stadium: "Estadio Universitario, Caracas",
    homeHits: 8,
    awayHits: 6,
    homeErrors: 0,
    awayErrors: 1,
    pitcherProbableHome: "E. Rodríguez (3-1, 2.84 ERA)",
    pitcherProbableAway: "J. Hernández (2-3, 4.12 ERA)",
  },
  {
    id: "g2",
    homeTeam: getTeam("cdl"),
    awayTeam: getTeam("adz"),
    homeScore: 2,
    awayScore: 5,
    inning: 5,
    isTopInning: true,
    status: "live",
    scheduledAt: "2026-09-07T19:00:00-04:00",
    stadium: "Estadio Antonio Herrera Gutiérrez, Barquisimeto",
    homeHits: 4,
    awayHits: 9,
    homeErrors: 2,
    awayErrors: 0,
    pitcherProbableHome: "R. Rivero (4-2, 3.10 ERA)",
    pitcherProbableAway: "M. Castillo (5-0, 1.95 ERA)",
  },
  {
    id: "g-sep7-3",
    homeTeam: getTeam("tdlg"),
    awayTeam: getTeam("tda"),
    homeScore: 6,
    awayScore: 2,
    inning: 9,
    isTopInning: false,
    status: "final",
    scheduledAt: "2026-09-07T17:00:00-04:00",
    stadium: "Estadio Jorge Luis García Carneiro, Macuto",
    homeHits: 10,
    awayHits: 5,
    homeErrors: 1,
    awayErrors: 2,
    pitcherWin: "A. Idrogo (4-1)",
    pitcherLoss: "G. Moscoso (2-4)",
    pitcherSave: "A. Cavanerio (8)",
  },
  {
    id: "g-sep7-4",
    homeTeam: getTeam("cda"),
    awayTeam: getTeam("bdm"),
    homeScore: 8,
    awayScore: 7,
    inning: 9,
    isTopInning: false,
    status: "final",
    scheduledAt: "2026-09-07T19:00:00-04:00",
    stadium: "Estadio Alfonso 'Chico' Carrasquel, Puerto La Cruz",
    homeHits: 12,
    awayHits: 11,
    homeErrors: 0,
    awayErrors: 1,
    pitcherWin: "L. Chirinos (2-0)",
    pitcherLoss: "C. Navas (1-3)",
    pitcherSave: "R. Rodríguez (11)",
  },

  // 8 de Septiembre (Mañana - Programados)
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
    stadium: "Estadio José Pérez Colmenares, Maracay",
    pitcherProbableHome: "Y. Pino (2-2, 3.45 ERA)",
    pitcherProbableAway: "R. Pinto (3-1, 2.90 ERA)",
  },
  {
    id: "g4",
    homeTeam: getTeam("bdm"),
    awayTeam: getTeam("tda"),
    homeScore: 0,
    awayScore: 0,
    inning: 0,
    isTopInning: true,
    status: "scheduled",
    scheduledAt: "2026-09-08T19:00:00-04:00",
    stadium: "Estadio Nueva Esparta, Porlamar",
    pitcherProbableHome: "F. Morales (1-4, 4.80 ERA)",
    pitcherProbableAway: "J. Martínez (3-3, 3.75 ERA)",
  },
  {
    id: "g-sep8-3",
    homeTeam: getTeam("ndm"),
    awayTeam: getTeam("ldc"),
    homeScore: 0,
    awayScore: 0,
    inning: 0,
    isTopInning: true,
    status: "scheduled",
    scheduledAt: "2026-09-08T19:00:00-04:00",
    stadium: "Estadio José Bernardo Pérez, Valencia",
    pitcherProbableHome: "E. Leal (4-1, 2.65 ERA)",
    pitcherProbableAway: "A. Rondón (5-2, 3.05 ERA)",
  },
  {
    id: "g-sep8-4",
    homeTeam: getTeam("adz"),
    awayTeam: getTeam("cda"),
    homeScore: 0,
    awayScore: 0,
    inning: 0,
    isTopInning: true,
    status: "scheduled",
    scheduledAt: "2026-09-08T19:00:00-04:00",
    stadium: "Estadio Luis Aparicio 'El Grande', Maracaibo",
    pitcherProbableHome: "S. Tomshaw (3-0, 2.15 ERA)",
    pitcherProbableAway: "D. Tomalin (1-2, 4.30 ERA)",
  },

  // 6 de Septiembre (Ayer - Finalizados)
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
    stadium: "Estadio Luis Aparicio 'El Grande', Maracaibo",
    homeHits: 9,
    awayHits: 13,
    homeErrors: 2,
    awayErrors: 0,
    pitcherWin: "R. Suárez (3-1)",
    pitcherLoss: "E. Paredes (0-2)",
    pitcherSave: "A. Vizcaíno (5)",
  },
  {
    id: "g-sep6-2",
    homeTeam: getTeam("ndm"),
    awayTeam: getTeam("cdl"),
    homeScore: 3,
    awayScore: 1,
    inning: 9,
    isTopInning: false,
    status: "final",
    scheduledAt: "2026-09-06T17:30:00-04:00",
    stadium: "Estadio José Bernardo Pérez, Valencia",
    homeHits: 7,
    awayHits: 4,
    homeErrors: 0,
    awayErrors: 1,
    pitcherWin: "Y. Méndez (5-1)",
    pitcherLoss: "N. Molina (3-2)",
    pitcherSave: "A. Machado (9)",
  },

  // 9 de Septiembre (Próximos)
  {
    id: "g-sep9-1",
    homeTeam: getTeam("ldc"),
    awayTeam: getTeam("tdlg"),
    homeScore: 0,
    awayScore: 0,
    inning: 0,
    isTopInning: true,
    status: "scheduled",
    scheduledAt: "2026-09-09T19:00:00-04:00",
    stadium: "Estadio Universitario, Caracas",
    pitcherProbableHome: "J. Mujica (2-1, 3.20 ERA)",
    pitcherProbableAway: "E. Barboza (1-0, 2.50 ERA)",
  },
  {
    id: "g-sep9-2",
    homeTeam: getTeam("cdl"),
    awayTeam: getTeam("tda"),
    homeScore: 0,
    awayScore: 0,
    inning: 0,
    isTopInning: true,
    status: "scheduled",
    scheduledAt: "2026-09-09T19:00:00-04:00",
    stadium: "Estadio Antonio Herrera Gutiérrez, Barquisimeto",
    pitcherProbableHome: "M. Socolovich (3-3, 3.88 ERA)",
    pitcherProbableAway: "A. Benítez (2-4, 4.50 ERA)",
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
