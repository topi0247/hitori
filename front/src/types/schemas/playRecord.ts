import { z } from "zod";

const Answer = z.object({
  uuid: z.string(),
  order: z.number().int().positive(),
});

const PlayResultCard = z.object({
  uuid: z.string(),
  card_number: z.number().int().min(1).max(100),
  word: z.string(),
  is_correct: z.boolean(),
});

export const PostPlayRecordRequest = z.object({
  theme_id: z.number(),
  answers: z.array(Answer).min(1),
});

export const PostPlayRecordResponse = z.object({
  correct_rate: z.number().min(0).max(100),
  cards: z.array(PlayResultCard),
});

export type PostPlayRecordRequest = z.infer<typeof PostPlayRecordRequest>;
export type PostPlayRecordResponse = z.infer<typeof PostPlayRecordResponse>;
