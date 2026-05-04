import { describe, it, expect } from "vitest";
import { Theme, ThemesResponse } from "@/types/schemas/theme";

describe("Theme", () => {
  it("有効なthemeを受け入れる", () => {
    const result = Theme.safeParse({ id: 1, title: "大きさ" });
    expect(result.success).toBe(true);
  });

  it("idがないと失敗する", () => {
    const result = Theme.safeParse({ title: "大きさ" });
    expect(result.success).toBe(false);
  });

  it("titleがないと失敗する", () => {
    const result = Theme.safeParse({ id: 1 });
    expect(result.success).toBe(false);
  });
});

describe("ThemesResponse", () => {
  it("有効なthemes配列を受け入れる", () => {
    const result = ThemesResponse.safeParse({
      themes: [{ id: 1, title: "大きさ" }],
    });
    expect(result.success).toBe(true);
  });

  it("空配列は有効", () => {
    const result = ThemesResponse.safeParse({ themes: [] });
    expect(result.success).toBe(true);
  });

  it("themesがないと失敗する", () => {
    const result = ThemesResponse.safeParse({});
    expect(result.success).toBe(false);
  });
});
