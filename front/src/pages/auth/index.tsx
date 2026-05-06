import { useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { useSignInWithEmail, useSignUpWithEmail } from "@/hooks/useAuth";
import { useCreateProfile } from "@/hooks/useProfile";

const AuthPage = () => {
  const navigate = useNavigate();
  const [loginEmail, setLoginEmail] = useState("");
  const [loginPassword, setLoginPassword] = useState("");
  const [signupEmail, setSignupEmail] = useState("");
  const [signupPassword, setSignupPassword] = useState("");
  const [signupName, setSignupName] = useState("");
  const [signupConfirmPassword, setSignupConfirmPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const { signIn } = useSignInWithEmail();
  const { signUp } = useSignUpWithEmail();
  const { mutateAsync: createProfile } = useCreateProfile();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    try {
      await signIn(loginEmail, loginPassword);
      await navigate({ to: "/" });
    } catch (e) {
      setError(e instanceof Error ? e.message : "ログインに失敗しました");
    }
  };

  const handleSignup = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    if (signupPassword !== signupConfirmPassword) {
      setError("パスワードが一致しません");
      return;
    }
    try {
      const session = await signUp(signupEmail, signupPassword);
      if (!session) throw new Error("セッションの取得に失敗しました");
      await createProfile({ token: session.access_token, user_name: signupName });
      await navigate({ to: "/" });
    } catch (e) {
      setError(e instanceof Error ? e.message : "登録に失敗しました");
    }
  };

  return (
    <div>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <section>
        <h2>ログイン</h2>
        <form onSubmit={handleLogin}>
          <input
            type="email"
            placeholder="メールアドレス"
            value={loginEmail}
            onChange={(e) => setLoginEmail(e.target.value)}
          />
          <input
            type="password"
            placeholder="パスワード"
            value={loginPassword}
            onChange={(e) => setLoginPassword(e.target.value)}
          />
          <button type="submit">ログイン</button>
        </form>
      </section>

      <hr />

      <section>
        <h2>新規登録</h2>
        <form onSubmit={handleSignup}>
          <input
            type="text"
            placeholder="名前（10文字以内）"
            value={signupName}
            onChange={(e) => setSignupName(e.target.value)}
          />
          <input
            type="email"
            placeholder="メールアドレス"
            value={signupEmail}
            onChange={(e) => setSignupEmail(e.target.value)}
          />
          <input
            type="password"
            placeholder="パスワード"
            value={signupPassword}
            onChange={(e) => setSignupPassword(e.target.value)}
          />
          <input
            type="password"
            placeholder="パスワード（確認）"
            value={signupConfirmPassword}
            onChange={(e) => setSignupConfirmPassword(e.target.value)}
          />
          <button type="submit">登録</button>
        </form>
      </section>
    </div>
  );
};

export { AuthPage };
