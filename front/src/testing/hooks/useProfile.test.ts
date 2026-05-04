import { renderHook, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { apiFetch } from "@/lib/api";
import { useDeleteProfile, usePatchProfile, useProfile } from "@/hooks/useProfile";
import { createWrapper } from "@/testing/wrapper";

vi.mock("@/lib/api");
const mockApiFetch = vi.mocked(apiFetch);

describe("useProfile", () => {
  afterEach(() => vi.clearAllMocks());

  it("プロフィールを取得できる", async () => {
    mockApiFetch.mockResolvedValue({ user_name: "たろう" });
    const { result } = renderHook(() => useProfile("token"), { wrapper: createWrapper() });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual({ user_name: "たろう" });
  });

  it("取得エラー時は isError が true になる", async () => {
    mockApiFetch.mockRejectedValue(new Error("Unauthorized"));
    const { result } = renderHook(() => useProfile("token"), { wrapper: createWrapper() });
    await waitFor(() => expect(result.current.isError).toBe(true));
  });
});

describe("usePatchProfile", () => {
  afterEach(() => vi.clearAllMocks());

  it("プロフィールを更新できる", async () => {
    mockApiFetch.mockResolvedValue({ user_name: "じろう" });
    const { result } = renderHook(() => usePatchProfile("token"), { wrapper: createWrapper() });
    result.current.mutate({ user_name: "じろう" });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual({ user_name: "じろう" });
  });

  it("更新エラー時は isError が true になる", async () => {
    mockApiFetch.mockRejectedValue(new Error("ValidationError"));
    const { result } = renderHook(() => usePatchProfile("token"), { wrapper: createWrapper() });
    result.current.mutate({ user_name: "じろう" });
    await waitFor(() => expect(result.current.isError).toBe(true));
  });
});

describe("useDeleteProfile", () => {
  afterEach(() => vi.clearAllMocks());

  it("プロフィールを削除できる", async () => {
    mockApiFetch.mockResolvedValue(undefined);
    const { result } = renderHook(() => useDeleteProfile("token"), { wrapper: createWrapper() });
    result.current.mutate();
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });
});
