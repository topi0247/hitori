import { describe, it, expect } from "vitest";
import { PostPlayRecordRequest, PostPlayRecordResponse } from "../types/playRecord";

describe("PostPlayRecordRequest", () => {
  const validAnswers = [
    { uuid: "uuid-xxxx", order: 1 },
    { uuid: "uuid-yyyy", order: 2 },
    { uuid: "uuid-zzzz", order: 3 },
    { uuid: "uuid-aaaa", order: 4 },
  ];

  it("有効なリクエストを受け入れる", () => {
    const result = PostPlayRecordRequest.safeParse({
      theme_id: 1,
      card_amount: 6,
      answers: validAnswers,
    });
    expect(result.success).toBe(true);
  });

  it("card_amount が 4 は有効", () => {
    const result = PostPlayRecordRequest.safeParse({
      theme_id: 1,
      card_amount: 4,
      answers: validAnswers,
    });
    expect(result.success).toBe(true);
  });

  it("card_amount が 10 は有効", () => {
    const result = PostPlayRecordRequest.safeParse({
      theme_id: 1,
      card_amount: 10,
      answers: validAnswers,
    });
    expect(result.success).toBe(true);
  });

  it("card_amount が 3 は失敗する", () => {
    const result = PostPlayRecordRequest.safeParse({
      theme_id: 1,
      card_amount: 3,
      answers: validAnswers,
    });
    expect(result.success).toBe(false);
  });

  it("card_amount が 11 は失敗する", () => {
    const result = PostPlayRecordRequest.safeParse({
      theme_id: 1,
      card_amount: 11,
      answers: validAnswers,
    });
    expect(result.success).toBe(false);
  });

  it("answers が空配列は失敗する", () => {
    const result = PostPlayRecordRequest.safeParse({
      theme_id: 1,
      card_amount: 6,
      answers: [],
    });
    expect(result.success).toBe(false);
  });

  it("theme_id がないと失敗する", () => {
    const result = PostPlayRecordRequest.safeParse({
      card_amount: 6,
      answers: validAnswers,
    });
    expect(result.success).toBe(false);
  });
});

describe("PostPlayRecordResponse", () => {
  it("有効なレスポンスを受け入れる", () => {
    const result = PostPlayRecordResponse.safeParse({
      correct_rate: 83.33,
      cards: [
        { uuid: "uuid-xxxx", card_number: 15, word: "アリ", is_correct: true },
        { uuid: "uuid-yyyy", card_number: 32, word: "クジラ", is_correct: false },
      ],
    });
    expect(result.success).toBe(true);
  });

  it("correct_rate が 0 は有効", () => {
    const result = PostPlayRecordResponse.safeParse({
      correct_rate: 0,
      cards: [{ uuid: "uuid-xxxx", card_number: 15, word: "アリ", is_correct: false }],
    });
    expect(result.success).toBe(true);
  });

  it("correct_rate が 100 は有効", () => {
    const result = PostPlayRecordResponse.safeParse({
      correct_rate: 100,
      cards: [{ uuid: "uuid-xxxx", card_number: 15, word: "アリ", is_correct: true }],
    });
    expect(result.success).toBe(true);
  });

  it("correct_rate が 0 未満は失敗する", () => {
    const result = PostPlayRecordResponse.safeParse({
      correct_rate: -1,
      cards: [],
    });
    expect(result.success).toBe(false);
  });

  it("correct_rate が 100 超は失敗する", () => {
    const result = PostPlayRecordResponse.safeParse({
      correct_rate: 100.01,
      cards: [],
    });
    expect(result.success).toBe(false);
  });
});
