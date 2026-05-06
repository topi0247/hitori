import { useAtom, useSetAtom } from "jotai";
import { sessionAtom } from "@/stores/auth";
import { supabase } from "@/lib/supabase";

const useSession = () => {
  const [session] = useAtom(sessionAtom);
  return session;
};

const useSignInWithGoogle = () => {
  const signIn = async () => {
    await supabase.auth.signInWithOAuth({
      provider: "google",
      options: { redirectTo: window.location.origin },
    });
  };
  return { signIn };
};

const useSignInWithEmail = () => {
  const signIn = async (email: string, password: string) => {
    await supabase.auth.signInWithPassword({ email, password });
  };
  return { signIn };
};

const useSignUpWithEmail = () => {
  const signUp = async (email: string, password: string) => {
    const { data } = await supabase.auth.signUp({ email, password });
    if (data.session) return data.session;
    const { data: signInData } = await supabase.auth.signInWithPassword({ email, password });
    return signInData.session;
  };
  return { signUp };
};

const useSignOut = () => {
  const setSession = useSetAtom(sessionAtom);
  const signOut = async () => {
    await supabase.auth.signOut();
    setSession(null);
  };
  return { signOut };
};

export { useSession, useSignInWithGoogle, useSignInWithEmail, useSignUpWithEmail, useSignOut };
