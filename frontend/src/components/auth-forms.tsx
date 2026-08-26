"use client";

import { useState } from "react";
import { api } from "@/lib/api";
import { Card, Field, SubmitButton, useStatus } from "./ui";

type LoginResponse = { user_id: number; token: string; username: string };

export function AuthCard() {
  const [tab, setTab] = useState<"register" | "login">("register");

  return (
    <Card
      step="01 · AUTH_SERVICE"
      title="Authentication"
      subtitle="gRPC → Postgres · bcrypt hashing · JWT sessions"
    >
      <div className="mb-5 grid grid-cols-2 gap-1 rounded-lg border border-edge bg-surface p-1">
        {(["register", "login"] as const).map((t) => (
          <button
            key={t}
            type="button"
            onClick={() => setTab(t)}
            className={`cursor-pointer rounded-md px-3 py-1.5 text-xs font-semibold capitalize transition ${
              tab === t
                ? "bg-panel-2 text-accent shadow-sm"
                : "text-muted hover:text-ink"
            }`}
          >
            {t}
          </button>
        ))}
      </div>

      <RegisterForm key={tab === "register" ? "reg" : "hidden"} hidden={tab === "login"} />
      <LoginForm key={tab === "login" ? "login" : "hidden"} hidden={tab === "register"} />
    </Card>
  );
}

function RegisterForm({ hidden }: { hidden: boolean }) {
  const { status, show, pending, node } = useStatus();

  if (hidden) return null;

  async function onSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = new FormData(e.currentTarget);
    show("Sending…");
    try {
      await api("/api/auth/register/", {
        method: "POST",
        body: JSON.stringify({
          username: form.get("username"),
          email: form.get("email"),
          password: form.get("password"),
          password_confirm: form.get("password_confirm"),
        }),
      });
      show("Registered successfully.", true);
      e.currentTarget.reset();
    } catch (err) {
      show(err instanceof Error ? err.message : "Request failed");
    }
  }

  return (
    <form onSubmit={onSubmit}>
      <Field id="username" label="Username" maxLength={32} />
      <Field id="email" label="Email" type="email" maxLength={32} />
      <div className="grid grid-cols-2 gap-3">
        <Field id="password" label="Password" type="password" maxLength={32} />
        <Field id="password_confirm" label="Confirm password" type="password" maxLength={32} />
      </div>
      <SubmitButton pending={pending} label="Create account" />
      {node}
    </form>
  );
}

function LoginForm({ hidden }: { hidden: boolean }) {
  const { status, show, pending, node } = useStatus();

  if (hidden) return null;

  async function onSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = new FormData(e.currentTarget);
    show("Sending…");
    try {
      const data = await api<LoginResponse>("/api/auth/login/", {
        method: "POST",
        body: JSON.stringify({
          username: form.get("login_username"),
          password: form.get("login_password"),
        }),
      });
      localStorage.setItem("token", data.token);
      show(`Logged in as ${data.username} — JWT stored (user #${data.user_id}).`, true);
    } catch (err) {
      show(err instanceof Error ? err.message : "Request failed");
    }
  }

  return (
    <form onSubmit={onSubmit}>
      <Field id="login_username" label="Username" maxLength={32} />
      <Field id="login_password" label="Password" type="password" maxLength={32} />
      <SubmitButton pending={pending} label="Sign in" />
      {node}
    </form>
  );
}
