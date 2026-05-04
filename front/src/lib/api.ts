export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080";

type FetchOptions = RequestInit & { token?: string };

const apiFetch = async <T>(path: string, options: FetchOptions = {}): Promise<T> => {
  const { token, ...init } = options;
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...(token && { Authorization: `Bearer ${token}` }),
    ...(init.headers as Record<string, string>),
  };
  const res = await fetch(`${API_BASE_URL}${path}`, { ...init, headers });
  if (!res.ok) {
    const body = (await res.json().catch(() => ({}))) as { message?: string };
    throw new Error(body.message ?? res.statusText);
  }
  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
};

export { apiFetch };
