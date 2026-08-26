import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Microservice Demo",
  description: "Auth, greet and chat microservices behind a Go gateway",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body className="min-h-screen px-4 py-8 font-sans antialiased">
        {children}
      </body>
    </html>
  );
}
