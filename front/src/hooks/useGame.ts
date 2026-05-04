import { useMutation } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { PostPlayRecordRequest, PostPlayRecordResponse } from "@/types/schemas/playRecord";

const usePlay = (token: string) =>
  useMutation({
    mutationFn: (data: PostPlayRecordRequest) =>
      apiFetch<PostPlayRecordResponse>("/play_records", {
        method: "POST",
        body: JSON.stringify(data),
        token,
      }),
  });

export { usePlay };
