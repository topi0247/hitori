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

const useAvailableCard = (themeId: number, options?: { enabled?: boolean }) =>
  useQuery({
    queryKey: ["cards", "available", themeId],
    queryFn: () => apiFetch<AvailableCardResponse>(`/themes/${themeId}/cards/available`),
    enabled: options?.enabled ?? true,
  });

const useGameCards = (themeId: number, cardAmount: number, options?: { enabled?: boolean }) =>
  useQuery({
    queryKey: ["cards", "game", themeId, cardAmount],
    queryFn: () =>
      apiFetch<GameCardsResponse>(`/themes/${themeId}/cards/game?card_amount=${cardAmount}`),
    enabled: options?.enabled ?? true,
  });

const useCreateCard = () =>
  useMutation({
    mutationFn: ({
      themeId,
      token,
      ...data
    }: { themeId: number; token?: string } & PostCardRequest) =>
      apiFetch<PostCardResponse>(`/themes/${themeId}/cards`, {
        method: "POST",
        body: JSON.stringify(data),
        token,
      }),
  });

const useConfirmCard = () =>
  useMutation({
    mutationFn: ({ id, token, ...data }: { id: number; token?: string } & PatchCardRequest) =>
      apiFetch<PatchCardResponse>(`/cards/${id}`, {
        method: "PATCH",
        body: JSON.stringify(data),
        token,
      }),
  });

const useDeleteCard = () =>
  useMutation({
    mutationFn: (id: number) => apiFetch<void>(`/cards/${id}`, { method: "DELETE" }),
  });

export { useAvailableCard, useGameCards, useCreateCard, useConfirmCard, useDeleteCard };
