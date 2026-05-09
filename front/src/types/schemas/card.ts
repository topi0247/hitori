import { z } from "zod";

const cardNumber = z.number().int().min(1).max(100);
const word = z.string().min(1).max(25);

export const AvailableCardResponse = z.object({
  card_number: cardNumber,
});

export const GameCard = z.object({
  uuid: z.string(),
  word: z.string(),
});

export const GameCardsResponse = z.object({
  cards: z.array(GameCard),
});

export const PostCardRequest = z.object({
  card_number: cardNumber,
  word,
  guest_name: z.string().min(1).max(10).optional(),
});

export const PostCardResponse = z.object({
  id: z.number(),
  uuid: z.string(),
  card_number: cardNumber,
  word: z.string(),
});

export const PatchCardRequest = z.object({
  word,
  guest_name: z.string().min(1).max(10).optional(),
});

export const PatchCardResponse = z.object({
  id: z.number(),
  card_number: cardNumber,
  word: z.string(),
});

export type AvailableCardResponse = z.infer<typeof AvailableCardResponse>;
export type GameCard = z.infer<typeof GameCard>;
export type GameCardsResponse = z.infer<typeof GameCardsResponse>;
export type PostCardRequest = z.infer<typeof PostCardRequest>;
export type PostCardResponse = z.infer<typeof PostCardResponse>;
export type PatchCardRequest = z.infer<typeof PatchCardRequest>;
export type PatchCardResponse = z.infer<typeof PatchCardResponse>;
