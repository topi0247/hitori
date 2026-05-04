import { useMutation, useQuery } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type {
  PatchProfileRequest,
  PatchProfileResponse,
  ProfileResponse,
} from "@/types/schemas/profile";

const useProfile = (token: string) =>
  useQuery({
    queryKey: ["profile"],
    queryFn: () => apiFetch<ProfileResponse>("/profile", { token }),
  });

const usePatchProfile = (token: string) =>
  useMutation({
    mutationFn: (data: PatchProfileRequest) =>
      apiFetch<PatchProfileResponse>("/profile", {
        method: "PATCH",
        body: JSON.stringify(data),
        token,
      }),
  });

const useDeleteProfile = (token: string) =>
  useMutation({
    mutationFn: () => apiFetch<void>("/profile", { method: "DELETE", token }),
  });

export { useProfile, usePatchProfile, useDeleteProfile };
