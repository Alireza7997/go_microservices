const NODES = [
  { name: "Next.js", role: "frontend", accent: "text-accent-2" },
  { name: "Gateway", role: "Go · HTTP :8000", accent: "text-accent" },
  { name: "auth", role: "gRPC :6001 · Postgres", accent: "text-ok" },
  { name: "chat", role: "gRPC :6002", accent: "text-ok" },
  { name: "greet", role: "gRPC :6003", accent: "text-ok" },
] as const;

function Node({ name, role, accent }: (typeof NODES)[number]) {
  return (
    <div className="flex min-w-24 flex-col items-center gap-1 rounded-lg border border-edge bg-panel px-3 py-2.5">
      <span className={`font-mono text-xs font-semibold ${accent}`}>{name}</span>
      <span className="text-center text-[0.65rem] text-muted">{role}</span>
    </div>
  );
}

function Arrow() {
  return <span aria-hidden className="hidden text-edge-hi md:inline">→</span>;
}

export function Architecture() {
  return (
    <div className="flex flex-wrap items-center justify-center gap-x-3 gap-y-4 rounded-2xl border border-edge bg-panel/50 px-6 py-5">
      <Node {...NODES[0]} />
      <Arrow />
      <Node {...NODES[1]} />
      <Arrow />
      <div className="flex flex-wrap items-center justify-center gap-2 md:flex-col">
        <div className="flex flex-col items-center gap-2 md:flex-row">
          <Node {...NODES[2]} />
          <span aria-hidden className="text-[0.65rem] text-muted md:hidden">· Postgres ·</span>
          <Node {...NODES[3]} />
        </div>
      </div>
      <Arrow />
      <Node {...NODES[4]} />
    </div>
  );
}
