import { atom } from "jotai";
import type { Session } from "@supabase/supabase-js";

const sessionAtom = atom<Session | null>(null);

export { sessionAtom };
