import { z } from "zod";

export const ProfileResponse = z.object({
  user_name: z.string(),
});

export const PatchProfileRequest = z.object({
  user_name: z.string().min(1).max(10),
});

export const PatchProfileResponse = ProfileResponse;

export type ProfileResponse = z.infer<typeof ProfileResponse>;
export type PatchProfileRequest = z.infer<typeof PatchProfileRequest>;
export type PatchProfileResponse = z.infer<typeof PatchProfileResponse>;
