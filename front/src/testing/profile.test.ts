import { describe, it, expect } from "vitest";
import { ProfileResponse, PatchProfileRequest } from "../types/profile";

describe("ProfileResponse", () => {
  it("有効なuser_nameを受け入れる", () => {
    const result = ProfileResponse.safeParse({ user_name: "たろう" });
    expect(result.success).toBe(true);
  });

  it("user_nameがないと失敗する", () => {
    const result = ProfileResponse.safeParse({});
    expect(result.success).toBe(false);
  });
});

describe("PatchProfileRequest", () => {
  it("有効なuser_nameを受け入れる", () => {
    const result = PatchProfileRequest.safeParse({ user_name: "じろう" });
    expect(result.success).toBe(true);
  });

  it("10文字ちょうどは有効", () => {
    const result = PatchProfileRequest.safeParse({ user_name: "あいうえおかきくけこ" });
    expect(result.success).toBe(true);
  });

  it("11文字以上は失敗する", () => {
    const result = PatchProfileRequest.safeParse({ user_name: "あいうえおかきくけこさ" });
    expect(result.success).toBe(false);
  });

  it("空文字は失敗する", () => {
    const result = PatchProfileRequest.safeParse({ user_name: "" });
    expect(result.success).toBe(false);
  });

  it("user_nameがないと失敗する", () => {
    const result = PatchProfileRequest.safeParse({});
    expect(result.success).toBe(false);
  });
});
