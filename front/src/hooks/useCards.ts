import { useMutation, useQuery } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type {
  AvailableCardResponse,
  GameCardsResponse,
  PatchCardRequest,
  PatchCardResponse,
  PostCardRequest,
  PostCardResponse,
} from "@/types/schemas/card";

const useAvailableCard = (themeId: number) =>
  useQuery({
    queryKey: ["cards", "available", themeId],
    queryFn: () => apiFetch<AvailableCardResponse>(`/themes/${themeId}/cards/available`),
  });

const useGameCards = (themeId: number, cardAmount: number) =>
  useQuery({
    queryKey: ["cards", "game", themeId, cardAmount],
    queryFn: () =>
      apiFetch<GameCardsResponse>(`/themes/${themeId}/cards/game?card_amount=${cardAmount}`),
  });

const useCreateCard = () =>
  useMutation({
    mutationFn: ({ themeId, ...data }: { themeId: number } & PostCardRequest) =>
      apiFetch<PostCardResponse>(`/themes/${themeId}/cards`, {
        method: "POST",
        body: JSON.stringify(data),
      }),
  });

const useConfirmCard = () =>
  useMutation({
    mutationFn: ({ id, ...data }: { id: number } & PatchCardRequest) =>
      apiFetch<PatchCardResponse>(`/cards/${id}`, {
        method: "PATCH",
        body: JSON.stringify(data),
      }),
  });

const useDeleteCard = () =>
  useMutation({
    mutationFn: (id: number) => apiFetch<void>(`/cards/${id}`, { method: "DELETE" }),
  });

export { useAvailableCard, useGameCards, useCreateCard, useConfirmCard, useDeleteCard };
