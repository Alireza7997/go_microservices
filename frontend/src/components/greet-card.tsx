"use client";

import { useState } from "react";
import { api } from "@/lib/api";
import { Card, Field, SubmitButton, useStatus } from "./ui";

type PingResponse = { greeting: string; server_time: number };

export function GreetCard() {
  const { status, show, pending, node } = useStatus();
  const [result, setResult] = useState<PingResponse | null>(null);

  async function onSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const name = String(new FormData(e.currentTarget).get("name") ?? "").trim() || "world";
    show("Sending…");
    try {
      const data = await api<PingResponse>(`/api/greet/ping/?name=${encodeURIComponent(name)}`);
      setResult(data);
      show(`Round-trip via gateway in ${new Date(data.server_time * 1000).toLocaleTimeString()}.`, true);
    } catch (err) {
      setResult(null);
      show(err instanceof Error ? err.message : "Request failed");
    }
  }

  return (
    <Card step="02 · GREET_SERVICE" title="Greeting" subtitle="Stateless RPC · no database">
      <form onSubmit={onSubmit}>
        <Field id="name" label="Your name" required={false} placeholder="world" />
        <SubmitButton pending={pending} label="Ping the service" />
        {node}
        {result && (
          <div className="mt-4 rounded-lg border border-edge bg-surface px-4 py-3 text-center">
            <p className="bg-linear-to-r from-accent to-accent-2 bg-clip-text text-lg font-semibold text-transparent">
              {result.greeting}
            </p>
          </div>
        )}
      </form>
    </Card>
  );
}
