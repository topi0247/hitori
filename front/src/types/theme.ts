import { z } from "zod";

export const Theme = z.object({
  id: z.number(),
  title: z.string(),
});

export const ThemesResponse = z.object({
  themes: z.array(Theme),
});

export type Theme = z.infer<typeof Theme>;
export type ThemesResponse = z.infer<typeof ThemesResponse>;
