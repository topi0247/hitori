import { renderHook, act } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { supabase } from "@/lib/supabase";
import { sessionAtom } from "@/stores/auth";
import {
  useSession,
  useSignInWithGoogle,
  useSignInWithEmail,
  useSignUpWithEmail,
  useSignOut,
} from "@/hooks/useAuth";
import { createWrapper } from "@/testing/wrapper";
import type { Session } from "@supabase/supabase-js";

vi.mock("@/lib/supabase", () => ({
  supabase: {
    auth: {
      signInWithOAuth: vi.fn().mockResolvedValue({ data: {}, error: null }),
      signInWithPassword: vi.fn().mockResolvedValue({
        data: { session: { access_token: "signin-token" } },
        error: null,
      }),
      signUp: vi.fn().mockResolvedValue({
        data: { session: { access_token: "test-token" } },
        error: null,
      }),
      signOut: vi.fn().mockResolvedValue({ error: null }),
    },
  },
}));

const mockSignInWithOAuth = vi.mocked(supabase.auth.signInWithOAuth);
const mockSignInWithPassword = vi.mocked(supabase.auth.signInWithPassword);
const mockSignUp = vi.mocked(supabase.auth.signUp);
const mockSignOut = vi.mocked(supabase.auth.signOut);

const mockSession = { access_token: "test-token", user: { id: "user-1" } } as unknown as Session;

describe("useSession", () => {
  afterEach(() => vi.clearAllMocks());

  it("セッションがない場合はnullを返す", () => {
    const { result } = renderHook(() => useSession(), { wrapper: createWrapper() });
    expect(result.current).toBeNull();
  });

  it("セッションがある場合はSessionを返す", () => {
    const { result } = renderHook(() => useSession(), {
      wrapper: createWrapper([[sessionAtom, mockSession]]),
    });
    expect(result.current).toEqual(mockSession);
  });
});

describe("useSignInWithGoogle", () => {
  afterEach(() => vi.clearAllMocks());

  it("GoogleのOAuthサインインを呼び出す", async () => {
    const { result } = renderHook(() => useSignInWithGoogle(), { wrapper: createWrapper() });
    await act(async () => {
      await result.current.signIn();
    });
    expect(mockSignInWithOAuth).toHaveBeenCalledWith({
      provider: "google",
      options: { redirectTo: expect.any(String) },
    });
  });
});

describe("useSignInWithEmail", () => {
  afterEach(() => vi.clearAllMocks());

  it("メール・パスワードでサインインを呼び出す", async () => {
    const { result } = renderHook(() => useSignInWithEmail(), { wrapper: createWrapper() });
    await act(async () => {
      await result.current.signIn("test@example.com", "password");
    });
    expect(mockSignInWithPassword).toHaveBeenCalledWith({
      email: "test@example.com",
      password: "password",
    });
  });
});

describe("useSignUpWithEmail", () => {
  afterEach(() => vi.clearAllMocks());

  it("サインアップしセッションを返す", async () => {
    const { result } = renderHook(() => useSignUpWithEmail(), { wrapper: createWrapper() });
    let session;
    await act(async () => {
      session = await result.current.signUp("test@example.com", "password");
    });
    expect(mockSignUp).toHaveBeenCalledWith({ email: "test@example.com", password: "password" });
    expect(session).toEqual({ access_token: "test-token" });
  });

  it("サインアップでsessionがnullの場合はサインインしてセッションを返す", async () => {
    mockSignUp.mockResolvedValueOnce({ data: { session: null }, error: null } as never);
    const { result } = renderHook(() => useSignUpWithEmail(), { wrapper: createWrapper() });
    let session;
    await act(async () => {
      session = await result.current.signUp("test@example.com", "password");
    });
    expect(mockSignInWithPassword).toHaveBeenCalledWith({
      email: "test@example.com",
      password: "password",
    });
    expect(session).toEqual({ access_token: "signin-token" });
  });
});

describe("useSignOut", () => {
  afterEach(() => vi.clearAllMocks());

  it("サインアウトを呼び出す", async () => {
    const { result } = renderHook(() => useSignOut(), { wrapper: createWrapper() });
    await act(async () => {
      await result.current.signOut();
    });
    expect(mockSignOut).toHaveBeenCalled();
  });
});
