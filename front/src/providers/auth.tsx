import { useEffect } from "react";
import { useSetAtom } from "jotai";
import { supabase } from "@/lib/supabase";
import { sessionAtom } from "@/stores/auth";
import type { ReactNode } from "react";

const AuthProvider = ({ children }: { children: ReactNode }) => {
  const setSession = useSetAtom(sessionAtom);

  useEffect(() => {
    supabase.auth.getSession().then(({ data: { session } }) => {
      setSession(session);
    });

    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange((_event, session) => {
      setSession(session);
    });

    return () => subscription.unsubscribe();
  }, [setSession]);

  return <>{children}</>;
};

export { AuthProvider };
