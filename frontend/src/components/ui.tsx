import { useState } from "react";

export function Card({
  step,
  title,
  subtitle,
  children,
  className = "",
}: {
  step: string;
  title: string;
  subtitle: string;
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <section
      className={`group relative rounded-2xl border border-edge bg-panel/80 p-6 shadow-xl shadow-black/30 backdrop-blur transition-colors hover:border-edge-hi ${className}`}
    >
      <div className="mb-5">
        <span className="font-mono text-[0.65rem] font-semibold tracking-widest text-accent/70">
          {step}
        </span>
        <h2 className="mt-1 text-base font-semibold tracking-tight text-ink">{title}</h2>
        <p className="mt-0.5 text-xs text-muted">{subtitle}</p>
      </div>
      {children}
    </section>
  );
}

export function Field({
  id,
  label,
  type = "text",
  required = true,
  placeholder,
  maxLength,
}: {
  id: string;
  label: string;
  type?: string;
  required?: boolean;
  placeholder?: string;
  maxLength?: number;
}) {
  return (
    <div className="mb-3">
      <label htmlFor={id} className="mb-1 block text-[0.7rem] font-medium tracking-wide text-muted uppercase">
        {label}
      </label>
      <input
        id={id}
        name={id}
        type={type}
        required={required}
        placeholder={placeholder}
        maxLength={maxLength}
        className="w-full rounded-lg border border-edge bg-surface px-3 py-2 text-sm text-ink transition outline-none placeholder:text-muted/40 focus:border-accent focus:ring-2 focus:ring-accent/20"
      />
    </div>
  );
}

export function SubmitButton({ pending, label }: { pending: boolean; label: string }) {
  return (
    <button
      type="submit"
      disabled={pending}
      className="mt-1 w-full cursor-pointer rounded-lg bg-linear-to-r from-accent to-accent-2 px-4 py-2.5 text-sm font-semibold text-surface transition hover:brightness-110 active:scale-[0.99] disabled:cursor-wait disabled:opacity-50"
    >
      {pending ? (
        <span className="inline-flex items-center gap-2">
          <Spinner /> Sending…
        </span>
      ) : (
        label
      )}
    </button>
  );
}

export function Spinner() {
  return (
    <span className="inline-block size-3 animate-spin rounded-full border-2 border-current border-t-transparent" />
  );
}

export function Status({ status }: { status: { message: string; ok: boolean } | null }) {
  if (!status) return <p className="min-h-5" />;
  return (
    <p
      role="status"
      className={`mt-3 flex items-start gap-1.5 text-xs break-words ${
        status.ok ? "text-ok" : "text-err"
      }`}
    >
      <span aria-hidden>{status.ok ? "✓" : "✕"}</span>
      {status.message}
    </p>
  );
}

export function useStatus() {
  const [status, setStatus] = useState<{ message: string; ok: boolean } | null>(null);
  const show = (message: string, ok = false) => setStatus({ message, ok });
  const pending = Boolean(status && !status.ok && status.message === "Sending…");
  return { status, show, pending, node: <Status status={status} /> };
}
