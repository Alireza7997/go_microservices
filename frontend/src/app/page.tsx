import { AuthCard } from "@/components/auth-forms";
import { ChatCard } from "@/components/chat-card";
import { GreetCard } from "@/components/greet-card";
import { Architecture } from "@/components/architecture";

const GATEWAY = process.env.NEXT_PUBLIC_GATEWAY_URL ?? "http://localhost:8000";

export default function Home() {
  return (
    <main className="mx-auto max-w-6xl px-1">
      <header className="py-14 text-center">
        <span className="inline-block rounded-full border border-edge bg-panel px-4 py-1 font-mono text-[0.7rem] tracking-widest text-accent uppercase">
          Go · gRPC · Docker
        </span>
        <h1 className="mt-5 text-4xl font-bold tracking-tight md:text-5xl">
          <span className="bg-linear-to-r from-accent to-accent-2 bg-clip-text text-transparent">
            Microservice Demo
          </span>
        </h1>
        <p className="mx-auto mt-3 max-w-xl text-sm leading-relaxed text-muted">
          Three independent services communicating over gRPC, fronted by a single Go
          gateway. Try each one below — no setup required.
        </p>
        <p className="mt-5 inline-flex items-center gap-2 rounded-full border border-edge bg-panel px-4 py-1.5 font-mono text-xs text-muted">
          <span className="size-1.5 animate-pulse rounded-full bg-ok" />
          gateway live at {GATEWAY}
        </p>
      </header>

      <Architecture />

      <div className="grid items-start gap-6 py-10 md:grid-cols-2 lg:grid-cols-3">
        <AuthCard />
        <GreetCard />
        <ChatCard className="lg:row-span-2" />
      </div>

      <footer className="border-t border-edge/60 py-8 text-center text-xs text-muted">
        register → login (JWT) → greet → chat · every request travels through the gateway
      </footer>
    </main>
  );
}
