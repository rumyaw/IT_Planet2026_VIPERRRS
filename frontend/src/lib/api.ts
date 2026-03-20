export function apiBaseUrl() {
  return process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";
}

export async function apiGet<T>(path: string, params?: Record<string, string | number | boolean | undefined>) {
  const base = apiBaseUrl();
  const url = new URL(base + path, window.location.origin);

  if (params) {
    for (const [k, v] of Object.entries(params)) {
      if (v === undefined) continue;
      url.searchParams.set(k, String(v));
    }
  }

  const res = await fetch(url.toString(), {
    method: "GET",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
  });

  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(`API GET ${path} failed: ${res.status} ${text}`);
  }

  return (await res.json()) as T;
}

export async function apiPost<T>(
  path: string,
  body: unknown,
  opts?: { signal?: AbortSignal }
): Promise<T> {
  const base = apiBaseUrl();
  const url = new URL(base + path, window.location.origin);

  const res = await fetch(url.toString(), {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
    signal: opts?.signal,
  });

  if (!res.ok) {
    const text = await res.text().catch(() => '');
    throw new Error(`API POST ${path} failed: ${res.status} ${text}`);
  }

  return (await res.json()) as T;
}

export async function apiPatch<T>(path: string, body: unknown): Promise<T> {
  const base = apiBaseUrl();
  const url = new URL(base + path, window.location.origin);

  const res = await fetch(url.toString(), {
    method: 'PATCH',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });

  if (!res.ok) {
    const text = await res.text().catch(() => '');
    throw new Error(`API PATCH ${path} failed: ${res.status} ${text}`);
  }

  return (await res.json()) as T;
}

