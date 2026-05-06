import { z } from "zod";

export const ProfileResponse = z.object({
  user_name: z.string(),
});

export const CreateProfileRequest = z.object({
  user_name: z.string().min(1).max(10),
});

export const CreateProfileResponse = ProfileResponse;

export const PatchProfileRequest = z.object({
  user_name: z.string().min(1).max(10),
});

export const PatchProfileResponse = ProfileResponse;

export type ProfileResponse = z.infer<typeof ProfileResponse>;
export type CreateProfileRequest = z.infer<typeof CreateProfileRequest>;
export type CreateProfileResponse = z.infer<typeof CreateProfileResponse>;
export type PatchProfileRequest = z.infer<typeof PatchProfileRequest>;
export type PatchProfileResponse = z.infer<typeof PatchProfileResponse>;
