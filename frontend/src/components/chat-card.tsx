"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { api } from "@/lib/api";
import { Card, Field, Spinner, useStatus } from "./ui";

type Message = { id: number; room: string; username: string; body: string; sent_at: number };

const ROOMS = ["general", "random", "dev"] as const;

function MessageList({ messages }: { messages: Message[] }) {
  const boxRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    boxRef.current?.scrollTo({ top: boxRef.current.scrollHeight });
  }, [messages]);

  return (
    <div
      ref={boxRef}
      className="mb-4 flex h-72 flex-col gap-2.5 overflow-y-auto rounded-lg border border-edge bg-surface p-3"
    >
      {messages.length === 0 ? (
        <p className="my-auto text-center text-xs text-muted">
          No messages yet — be the first to say hello.
        </p>
      ) : (
        messages.map((m) => (
          <div key={m.id} className="max-w-[85%] self-start rounded-lg bg-panel-2 px-3 py-2">
            <p className="text-xs font-semibold text-accent">
              {m.username}
              <span className="ml-2 font-normal text-muted/70">
                {new Date(m.sent_at * 1000).toLocaleTimeString()}
              </span>
            </p>
            <p className="mt-0.5 text-sm leading-snug break-words text-ink">{m.body}</p>
          </div>
        ))
      )}
    </div>
  );
}

export function ChatCard({ className = "" }: { className?: string }) {
  const { status, show, node } = useStatus();
  const [room, setRoom] = useState<string>("general");
  const [messages, setMessages] = useState<Message[]>([]);
  const [sending, setSending] = useState(false);

  const refresh = useCallback(async (targetRoom: string) => {
    if (!targetRoom) return;
    try {
      setMessages(await api<Message[]>(`/api/chat/messages/${encodeURIComponent(targetRoom)}/`));
    } catch {
      setMessages([]);
    }
  }, []);

  useEffect(() => {
    refresh(room);
    const timer = setInterval(() => refresh(room), 3000);
    return () => clearInterval(timer);
  }, [room, refresh]);

  async function onSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = new FormData(e.currentTarget);
    const body = String(form.get("body") ?? "").trim();
    if (!body) return;
    setSending(true);
    try {
      await api("/api/chat/messages/", {
        method: "POST",
        body: JSON.stringify({ room, username: form.get("chat_username"), body }),
      });
      e.currentTarget.reset();
      await refresh(room);
    } catch (err) {
      show(err instanceof Error ? err.message : "Request failed");
    } finally {
      setSending(false);
    }
  }

  return (
    <Card step="03 · CHAT_SERVICE" title="Live chat" subtitle="In-memory rooms · polled every 3s" className={className}>
      <div className="mb-4 flex gap-1 rounded-lg border border-edge bg-surface p-1">
        {ROOMS.map((r) => (
          <button
            key={r}
            type="button"
            onClick={() => setRoom(r)}
            className={`flex-1 cursor-pointer rounded-md px-3 py-1.5 font-mono text-xs transition ${
              room === r ? "bg-panel-2 text-accent shadow-sm" : "text-muted hover:text-ink"
            }`}
          >
            #{r}
          </button>
        ))}
      </div>

      <MessageList messages={messages} />

      <form onSubmit={onSubmit}>
        <Field id="chat_username" label="Username" maxLength={32} />
        <div className="flex gap-2">
          <input
            id="body"
            name="body"
            required
            autoComplete="off"
            placeholder="Type a message…"
            className="min-w-0 flex-1 rounded-lg border border-edge bg-surface px-3 py-2 text-sm text-ink transition outline-none placeholder:text-muted/40 focus:border-accent focus:ring-2 focus:ring-accent/20"
          />
          <button
            type="submit"
            disabled={sending}
            className="shrink-0 cursor-pointer rounded-lg border border-edge-hi bg-panel-2 px-4 text-sm font-semibold text-accent transition hover:border-accent hover:bg-panel active:scale-[0.98] disabled:cursor-wait disabled:opacity-50"
          >
            {sending ? <Spinner /> : "Send"}
          </button>
        </div>
        {node}
      </form>
    </Card>
  );
}
