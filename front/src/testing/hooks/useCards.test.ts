import { renderHook, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { apiFetch } from "@/lib/api";
import {
  useAvailableCard,
  useConfirmCard,
  useCreateCard,
  useDeleteCard,
  useGameCards,
} from "@/hooks/useCards";
import { createWrapper } from "@/testing/wrapper";

vi.mock("@/lib/api");
const mockApiFetch = vi.mocked(apiFetch);

describe("useAvailableCard", () => {
  afterEach(() => vi.clearAllMocks());

  it("空きcard_numberを取得できる", async () => {
    mockApiFetch.mockResolvedValue({ card_number: 42 });
    const { result } = renderHook(() => useAvailableCard(1), { wrapper: createWrapper() });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual({ card_number: 42 });
  });

  it("空きなし時は isError が true になる", async () => {
    mockApiFetch.mockRejectedValue(new Error("ConflictError"));
    const { result } = renderHook(() => useAvailableCard(1), { wrapper: createWrapper() });
    await waitFor(() => expect(result.current.isError).toBe(true));
  });
});

describe("useGameCards", () => {
  afterEach(() => vi.clearAllMocks());

  it("ゲーム用カードを取得できる", async () => {
    const cards = [
      { uuid: "uuid-xxxx", word: "アリ" },
      { uuid: "uuid-yyyy", word: "クジラ" },
    ];
    mockApiFetch.mockResolvedValue({ cards });
    const { result } = renderHook(() => useGameCards(1, 6), { wrapper: createWrapper() });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual({ cards });
  });

  it("取得エラー時は isError が true になる", async () => {
    mockApiFetch.mockRejectedValue(new Error("NotFound"));
    const { result } = renderHook(() => useGameCards(1, 6), { wrapper: createWrapper() });
    await waitFor(() => expect(result.current.isError).toBe(true));
  });
});

describe("useCreateCard", () => {
  afterEach(() => vi.clearAllMocks());

  it("カードを仮登録できる", async () => {
    mockApiFetch.mockResolvedValue({ id: 1, card_number: 42, word: "アリ" });
    const { result } = renderHook(() => useCreateCard(), { wrapper: createWrapper() });
    result.current.mutate({ themeId: 1, card_number: 42, word: "アリ" });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual({ id: 1, card_number: 42, word: "アリ" });
  });

  it("ゲストとしてカードを仮登録できる", async () => {
    mockApiFetch.mockResolvedValue({ id: 1, card_number: 42, word: "アリ" });
    const { result } = renderHook(() => useCreateCard(), { wrapper: createWrapper() });
    result.current.mutate({ themeId: 1, card_number: 42, word: "アリ", guest_name: "ゲスト太郎" });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });

  it("登録エラー時は isError が true になる", async () => {
    mockApiFetch.mockRejectedValue(new Error("ConflictError"));
    const { result } = renderHook(() => useCreateCard(), { wrapper: createWrapper() });
    result.current.mutate({ themeId: 1, card_number: 42, word: "アリ" });
    await waitFor(() => expect(result.current.isError).toBe(true));
  });
});

describe("useConfirmCard", () => {
  afterEach(() => vi.clearAllMocks());

  it("カードを本登録できる", async () => {
    mockApiFetch.mockResolvedValue({ id: 1, card_number: 42, word: "クジラ" });
    const { result } = renderHook(() => useConfirmCard(), { wrapper: createWrapper() });
    result.current.mutate({ id: 1, word: "クジラ" });
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual({ id: 1, card_number: 42, word: "クジラ" });
  });

  it("本登録エラー時は isError が true になる", async () => {
    mockApiFetch.mockRejectedValue(new Error("ForbiddenError"));
    const { result } = renderHook(() => useConfirmCard(), { wrapper: createWrapper() });
    result.current.mutate({ id: 1, word: "クジラ" });
    await waitFor(() => expect(result.current.isError).toBe(true));
  });
});

describe("useDeleteCard", () => {
  afterEach(() => vi.clearAllMocks());

  it("カードを削除できる", async () => {
    mockApiFetch.mockResolvedValue(undefined);
    const { result } = renderHook(() => useDeleteCard(), { wrapper: createWrapper() });
    result.current.mutate(1);
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
  });

  it("削除エラー時は isError が true になる", async () => {
    mockApiFetch.mockRejectedValue(new Error("ForbiddenError"));
    const { result } = renderHook(() => useDeleteCard(), { wrapper: createWrapper() });
    result.current.mutate(1);
    await waitFor(() => expect(result.current.isError).toBe(true));
  });
});
