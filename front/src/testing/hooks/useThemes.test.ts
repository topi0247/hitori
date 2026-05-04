import { renderHook, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { apiFetch } from "@/lib/api";
import { useThemes } from "@/hooks/useThemes";
import { createWrapper } from "@/testing/wrapper";

vi.mock("@/lib/api");
const mockApiFetch = vi.mocked(apiFetch);

describe("useThemes", () => {
  afterEach(() => vi.clearAllMocks());

  it("テーマ一覧を取得できる", async () => {
    mockApiFetch.mockResolvedValue({ themes: [{ id: 1, title: "大きさ" }] });
    const { result } = renderHook(() => useThemes(), { wrapper: createWrapper() });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual({ themes: [{ id: 1, title: "大きさ" }] });
  });

  it("テーマが0件でも取得できる", async () => {
    mockApiFetch.mockResolvedValue({ themes: [] });
    const { result } = renderHook(() => useThemes(), { wrapper: createWrapper() });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual({ themes: [] });
  });

  it("取得エラー時は isError が true になる", async () => {
    mockApiFetch.mockRejectedValue(new Error("Unexpected"));
    const { result } = renderHook(() => useThemes(), { wrapper: createWrapper() });
    await waitFor(() => expect(result.current.isError).toBe(true));
  });
});
