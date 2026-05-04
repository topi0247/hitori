import { renderHook, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { apiFetch } from "@/lib/api";
import { usePlay } from "@/hooks/useGame";
import { createWrapper } from "@/testing/wrapper";

vi.mock("@/lib/api");
const mockApiFetch = vi.mocked(apiFetch);

describe("usePlay", () => {
  afterEach(() => vi.clearAllMocks());

  const validRequest = {
    theme_id: 1,
    card_amount: 6,
    answers: [
      { uuid: "uuid-xxxx", order: 1 },
      { uuid: "uuid-yyyy", order: 2 },
      { uuid: "uuid-zzzz", order: 3 },
      { uuid: "uuid-aaaa", order: 4 },
    ],
  };

  const validResponse = {
    correct_rate: 83.33,
    cards: [
      { uuid: "uuid-xxxx", card_number: 15, word: "アリ", is_correct: true },
      { uuid: "uuid-yyyy", card_number: 32, word: "クジラ", is_correct: false },
    ],
  };

  it("プレイ結果を送信できる", async () => {
    mockApiFetch.mockResolvedValue(validResponse);
    const { result } = renderHook(() => usePlay("token"), { wrapper: createWrapper() });
    result.current.mutate(validRequest);
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data).toEqual(validResponse);
  });

  it("送信エラー時は isError が true になる", async () => {
    mockApiFetch.mockRejectedValue(new Error("Unauthorized"));
    const { result } = renderHook(() => usePlay("token"), { wrapper: createWrapper() });
    result.current.mutate(validRequest);
    await waitFor(() => expect(result.current.isError).toBe(true));
  });
});
