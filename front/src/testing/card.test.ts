import { describe, it, expect } from "vitest";
import {
  AvailableCardResponse,
  GameCard,
  GameCardsResponse,
  PostCardRequest,
  PostCardResponse,
  PatchCardRequest,
  PatchCardResponse,
} from "../types/card";

describe("AvailableCardResponse", () => {
  it("有効なcard_numberを受け入れる", () => {
    const result = AvailableCardResponse.safeParse({ card_number: 42 });
    expect(result.success).toBe(true);
  });

  it("1は有効", () => {
    const result = AvailableCardResponse.safeParse({ card_number: 1 });
    expect(result.success).toBe(true);
  });

  it("100は有効", () => {
    const result = AvailableCardResponse.safeParse({ card_number: 100 });
    expect(result.success).toBe(true);
  });

  it("0は失敗する", () => {
    const result = AvailableCardResponse.safeParse({ card_number: 0 });
    expect(result.success).toBe(false);
  });

  it("101は失敗する", () => {
    const result = AvailableCardResponse.safeParse({ card_number: 101 });
    expect(result.success).toBe(false);
  });
});

describe("GameCard", () => {
  it("有効なカードを受け入れる", () => {
    const result = GameCard.safeParse({ uuid: "uuid-xxxx", word: "アリ" });
    expect(result.success).toBe(true);
  });

  it("uuidがないと失敗する", () => {
    const result = GameCard.safeParse({ word: "アリ" });
    expect(result.success).toBe(false);
  });

  it("wordがないと失敗する", () => {
    const result = GameCard.safeParse({ uuid: "uuid-xxxx" });
    expect(result.success).toBe(false);
  });
});

describe("GameCardsResponse", () => {
  it("有効なcards配列を受け入れる", () => {
    const result = GameCardsResponse.safeParse({
      cards: [{ uuid: "uuid-xxxx", word: "アリ" }],
    });
    expect(result.success).toBe(true);
  });

  it("cardsがないと失敗する", () => {
    const result = GameCardsResponse.safeParse({});
    expect(result.success).toBe(false);
  });
});

describe("PostCardRequest", () => {
  it("会員ユーザーは guest_name なしで有効", () => {
    const result = PostCardRequest.safeParse({ card_number: 42, word: "アリ" });
    expect(result.success).toBe(true);
  });

  it("ゲストは guest_name ありで有効", () => {
    const result = PostCardRequest.safeParse({
      card_number: 42,
      word: "アリ",
      guest_name: "ゲスト太郎",
    });
    expect(result.success).toBe(true);
  });

  it("card_number が 1 は有効", () => {
    const result = PostCardRequest.safeParse({ card_number: 1, word: "アリ" });
    expect(result.success).toBe(true);
  });

  it("card_number が 100 は有効", () => {
    const result = PostCardRequest.safeParse({ card_number: 100, word: "アリ" });
    expect(result.success).toBe(true);
  });

  it("card_number が 0 は失敗する", () => {
    const result = PostCardRequest.safeParse({ card_number: 0, word: "アリ" });
    expect(result.success).toBe(false);
  });

  it("card_number が 101 は失敗する", () => {
    const result = PostCardRequest.safeParse({ card_number: 101, word: "アリ" });
    expect(result.success).toBe(false);
  });

  it("word が 25 文字は有効", () => {
    const result = PostCardRequest.safeParse({
      card_number: 42,
      word: "あ".repeat(25),
    });
    expect(result.success).toBe(true);
  });

  it("word が 26 文字以上は失敗する", () => {
    const result = PostCardRequest.safeParse({
      card_number: 42,
      word: "あ".repeat(26),
    });
    expect(result.success).toBe(false);
  });

  it("word が空文字は失敗する", () => {
    const result = PostCardRequest.safeParse({ card_number: 42, word: "" });
    expect(result.success).toBe(false);
  });

  it("guest_name が 10 文字は有効", () => {
    const result = PostCardRequest.safeParse({
      card_number: 42,
      word: "アリ",
      guest_name: "あ".repeat(10),
    });
    expect(result.success).toBe(true);
  });

  it("guest_name が 11 文字以上は失敗する", () => {
    const result = PostCardRequest.safeParse({
      card_number: 42,
      word: "アリ",
      guest_name: "あ".repeat(11),
    });
    expect(result.success).toBe(false);
  });
});

describe("PostCardResponse", () => {
  it("有効なレスポンスを受け入れる", () => {
    const result = PostCardResponse.safeParse({ id: 1, card_number: 42, word: "アリ" });
    expect(result.success).toBe(true);
  });
});

describe("PatchCardRequest", () => {
  it("有効な word を受け入れる", () => {
    const result = PatchCardRequest.safeParse({ word: "クジラ" });
    expect(result.success).toBe(true);
  });

  it("word が 25 文字は有効", () => {
    const result = PatchCardRequest.safeParse({ word: "あ".repeat(25) });
    expect(result.success).toBe(true);
  });

  it("word が 26 文字以上は失敗する", () => {
    const result = PatchCardRequest.safeParse({ word: "あ".repeat(26) });
    expect(result.success).toBe(false);
  });

  it("word が空文字は失敗する", () => {
    const result = PatchCardRequest.safeParse({ word: "" });
    expect(result.success).toBe(false);
  });

  it("word がないと失敗する", () => {
    const result = PatchCardRequest.safeParse({});
    expect(result.success).toBe(false);
  });
});

describe("PatchCardResponse", () => {
  it("有効なレスポンスを受け入れる", () => {
    const result = PatchCardResponse.safeParse({ id: 1, card_number: 42, word: "クジラ" });
    expect(result.success).toBe(true);
  });
});
