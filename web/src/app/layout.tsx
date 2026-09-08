import type { Metadata } from "next";
import { Inter } from "next/font/google";
import { Geist_Mono } from "next/font/google";
import { Header } from "@/components/header";
import { Footer } from "@/components/footer";
import "./globals.css";

const inter = Inter({
  variable: "--font-sans",
  subsets: ["latin"],
  display: "swap",
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "LVBP GameCast — Béisbol Profesional Venezolano en Vivo",
  description:
    "Sigue en vivo los juegos de la Liga Venezolana de Béisbol Profesional. Marcadores, posiciones, noticias y GameCast en tiempo real.",
  keywords: [
    "LVBP",
    "béisbol venezolano",
    "GameCast",
    "marcador en vivo",
    "posiciones",
  ],
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html
      lang="es"
      className={`${inter.variable} ${geistMono.variable} h-full antialiased`}
    >
      <body className="flex min-h-full flex-col bg-background">
        <Header />
        <main className="flex-1">{children}</main>
        <Footer />
      </body>
    </html>
  );
}
