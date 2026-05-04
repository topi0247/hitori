import { useQuery } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { ThemesResponse } from "@/types/schemas/theme";

const useThemes = () =>
  useQuery({
    queryKey: ["themes"],
    queryFn: () => apiFetch<ThemesResponse>("/themes"),
  });

export { useThemes };
