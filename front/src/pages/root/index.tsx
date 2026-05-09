import { useEffect, useMemo, useRef, useState } from "react";
import { useAtom } from "jotai";
import { tv } from "tailwind-variants";
import { sessionAtom } from "@/stores/auth";
import { useSignInWithEmail, useSignUpWithEmail, useSignOut } from "@/hooks/useAuth";
import { useCreateProfile } from "@/hooks/useProfile";
import { useThemes } from "@/hooks/useThemes";
import { useAvailableCard, useConfirmCard, useCreateCard, useGameCards } from "@/hooks/useCards";
import { usePlay } from "@/hooks/useGame";
import type { GameCard } from "@/types/schemas/card";
import type { PostPlayRecordResponse } from "@/types/schemas/playRecord";

type View = "top" | "auth" | "game";
type AuthTab = "login" | "signup";
type GameStep = "config" | "word" | "play" | "result";

// ---------------------------------------------------------------------------
// Variants
// ---------------------------------------------------------------------------

const btn = tv({
  base: "w-full py-4 font-medium disabled:opacity-40",
  variants: {
    variant: {
      primary: "bg-black text-white",
      outline: "border border-black bg-white text-black",
    },
  },
  defaultVariants: { variant: "primary" },
});

const textInput = tv({
  base: "w-full border border-black bg-white p-3",
});

const tab = tv({
  base: "flex-1 border-b-2 py-3 text-sm font-medium transition-opacity",
  variants: {
    active: {
      true: "border-black",
      false: "border-transparent opacity-40",
    },
  },
});

const slidePanel = tv({
  base: "absolute inset-0 transition-transform duration-300 ease-in-out",
});

const panelHeader = tv({
  base: "flex items-center gap-4 border-b border-black p-4 transition-opacity duration-500",
  variants: {
    hidden: { true: "pointer-events-none opacity-0" },
  },
});

// ---------------------------------------------------------------------------
// RootPage
// ---------------------------------------------------------------------------

const RootPage = () => {
  const [view, setView] = useState<View>("top");
  const [session] = useAtom(sessionAtom);
  const { signOut } = useSignOut();

  return (
    <div className="relative w-full overflow-hidden" style={{ height: "100svh" }}>
      <div
        className={slidePanel()}
        style={{
          transform:
            view === "auth" ? "translateX(-100%)" : view === "game" ? "translateY(-100%)" : "none",
        }}
      >
        <TopPanel
          isLoggedIn={!!session}
          onPlay={() => setView("game")}
          onAuth={() => setView("auth")}
          onSignOut={signOut}
        />
      </div>

      <div
        className={slidePanel()}
        style={{ transform: view === "auth" ? "none" : "translateX(100%)" }}
      >
        <AuthPanel onBack={() => setView("top")} onSuccess={() => setView("top")} />
      </div>

      <div
        className={slidePanel()}
        style={{ transform: view === "game" ? "none" : "translateY(100%)" }}
      >
        <GamePanel onBack={() => setView("top")} />
      </div>
    </div>
  );
};

// ---------------------------------------------------------------------------
// TopPanel
// ---------------------------------------------------------------------------

const TopPanel = ({
  isLoggedIn,
  onPlay,
  onAuth,
  onSignOut,
}: {
  isLoggedIn: boolean;
  onPlay: () => void;
  onAuth: () => void;
  onSignOut: () => void;
}) => (
  <div className="relative flex h-full flex-col items-center justify-center gap-10 bg-white px-6">
    <h1 className="text-5xl font-bold tracking-tight">h「ito」ri</h1>
    <div className="flex w-full max-w-xs flex-col gap-4">
      <button type="button" onClick={onPlay} className={btn()}>
        遊ぶ
      </button>
      {!isLoggedIn && (
        <button type="button" onClick={onAuth} className={btn({ variant: "outline" })}>
          登録・ログイン
        </button>
      )}
    </div>
    {isLoggedIn && (
      <button
        type="button"
        onClick={onSignOut}
        className="absolute right-4 bottom-4 text-sm text-gray-400 underline"
      >
        ログアウト
      </button>
    )}
  </div>
);

// ---------------------------------------------------------------------------
// AuthPanel
// ---------------------------------------------------------------------------

const AuthPanel = ({ onBack, onSuccess }: { onBack: () => void; onSuccess: () => void }) => {
  const [activeTab, setActiveTab] = useState<AuthTab>("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [name, setName] = useState("");
  const [error, setError] = useState<string | null>(null);

  const { signIn } = useSignInWithEmail();
  const { signUp } = useSignUpWithEmail();
  const { mutateAsync: createProfile } = useCreateProfile();

  const reset = () => {
    setEmail("");
    setPassword("");
    setConfirmPassword("");
    setName("");
    setError(null);
  };

  const handleBack = () => {
    reset();
    onBack();
  };
  const handleTabChange = (next: AuthTab) => {
    reset();
    setActiveTab(next);
  };

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    try {
      await signIn(email, password);
      reset();
      onSuccess();
    } catch {
      setError("メールアドレスまたはパスワードが正しくありません");
    }
  };

  const handleSignup = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    if (password !== confirmPassword) {
      setError("パスワードが一致しません");
      return;
    }
    try {
      const session = await signUp(email, password);
      if (!session) throw new Error();
      await createProfile({ token: session.access_token, user_name: name });
      reset();
      onSuccess();
    } catch {
      setError("登録に失敗しました");
    }
  };

  return (
    <div className="flex h-full flex-col bg-white">
      <div className={panelHeader()}>
        <button type="button" onClick={handleBack} className="font-medium">
          ←
        </button>
        <span className="font-medium">登録・ログイン</span>
      </div>

      <div className="flex border-b border-black">
        {(["login", "signup"] as AuthTab[]).map((t) => (
          <button
            key={t}
            type="button"
            onClick={() => handleTabChange(t)}
            className={tab({ active: activeTab === t })}
          >
            {t === "login" ? "ログイン" : "新規登録"}
          </button>
        ))}
      </div>

      <div className="flex-1 overflow-y-auto p-6">
        {error && <p className="mb-4 text-sm font-medium">{error}</p>}
        {activeTab === "login" ? (
          <form onSubmit={handleLogin} className="flex flex-col gap-4">
            <input
              type="email"
              placeholder="メールアドレス"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className={textInput()}
            />
            <input
              type="password"
              placeholder="パスワード"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className={textInput()}
            />
            <button type="submit" className={btn({ class: "mt-2" })}>
              ログイン
            </button>
          </form>
        ) : (
          <form onSubmit={handleSignup} className="flex flex-col gap-4">
            <input
              type="text"
              placeholder="名前（10文字以内）"
              value={name}
              onChange={(e) => setName(e.target.value)}
              maxLength={10}
              required
              className={textInput()}
            />
            <input
              type="email"
              placeholder="メールアドレス"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className={textInput()}
            />
            <input
              type="password"
              placeholder="パスワード"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className={textInput()}
            />
            <input
              type="password"
              placeholder="パスワード（確認）"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
              className={textInput()}
            />
            <button type="submit" className={btn({ class: "mt-2" })}>
              登録
            </button>
          </form>
        )}
      </div>
    </div>
  );
};

// ---------------------------------------------------------------------------
// PlayingCard — word step 用（数字表示）
// ---------------------------------------------------------------------------

const PlayingCard = ({ number }: { number: number | null }) => {
  const label = number === null ? "" : String(number);
  return (
    <div
      className="relative flex items-center justify-center border-2 border-black bg-white select-none"
      style={{
        width: "9rem",
        height: "13rem",
        borderRadius: "0.5rem",
        boxShadow: "4px 4px 0 #000",
      }}
    >
      {number === null ? (
        <span className="text-4xl font-bold opacity-20">…</span>
      ) : (
        <>
          <span className="absolute top-2 left-3 text-xl leading-none font-bold">{label}</span>
          <span className="text-6xl leading-none font-bold">{label}</span>
          <span
            className="absolute right-3 bottom-2 text-xl leading-none font-bold"
            style={{ transform: "rotate(180deg)" }}
          >
            {label}
          </span>
        </>
      )}
    </div>
  );
};

// ---------------------------------------------------------------------------
// WordCard — play step 用（ワード表示）
// ---------------------------------------------------------------------------

const WordCard = ({
  word,
  cardNumber,
  isOwn = false,
}: {
  word: string;
  cardNumber?: number;
  isOwn?: boolean;
}) => (
  <div
    className="flex flex-col items-center justify-center border-2 border-black bg-white select-none"
    style={{
      width: isOwn ? "9rem" : "7rem",
      height: isOwn ? "12rem" : "9.5rem",
      borderRadius: "0.5rem",
      boxShadow: "3px 3px 0 #000",
      padding: "0.75rem",
    }}
  >
    <span
      className="text-center leading-snug font-bold"
      style={{ fontSize: isOwn ? "1rem" : "0.85rem" }}
    >
      {word}
    </span>
    {cardNumber !== undefined && <span className="mt-2 text-xs opacity-40">{cardNumber}</span>}
  </div>
);

// ---------------------------------------------------------------------------
// PlayStep — drag & drop ゲーム画面
// ---------------------------------------------------------------------------

type PlayCardEntry = {
  uuid: string;
  word: string;
  rotate: number;
  isOwn: boolean;
  cardNumber?: number;
  inZone: boolean;
  zoneX: number;
  zoneY: number;
};

const OWN_UUID = "__own__";

type PlayAnswer = { uuid: string; order: number };

const PlayStep = ({
  ownCard,
  ownCardUuid,
  gameCards,
  scatterAngles,
  onSubmit,
}: {
  ownCard: { word: string; cardNumber: number };
  ownCardUuid: string;
  gameCards: GameCard[];
  scatterAngles: number[];
  onSubmit?: (answers: PlayAnswer[]) => void;
}) => {
  const zoneRef = useRef<HTMLDivElement>(null);
  const [isAligned, setIsAligned] = useState(false);

  const [cards, setCards] = useState<PlayCardEntry[]>(() => [
    ...gameCards.map((c, i) => ({
      uuid: c.uuid,
      word: c.word,
      rotate: scatterAngles[i] ?? 0,
      isOwn: false,
      inZone: false,
      zoneX: 0,
      zoneY: 0,
    })),
    {
      uuid: OWN_UUID,
      word: ownCard.word,
      rotate: 0,
      isOwn: true,
      cardNumber: ownCard.cardNumber,
      inZone: false,
      zoneX: 0,
      zoneY: 0,
    },
  ]);

  const [dragState, setDragState] = useState<{
    uuid: string;
    fixedX: number;
    fixedY: number;
  } | null>(null);

  // ドラッグ開始情報は ref で保持（pointermove での stale closure を回避）
  const dragStartRef = useRef<{
    uuid: string;
    fromZone: boolean;
    startPointerX: number;
    startPointerY: number;
    startCardFixedX: number;
    startCardFixedY: number;
  } | null>(null);

  useEffect(() => {
    const onMove = (e: PointerEvent) => {
      const ds = dragStartRef.current;
      if (!ds) return;
      setDragState({
        uuid: ds.uuid,
        fixedX: ds.startCardFixedX + (e.clientX - ds.startPointerX),
        fixedY: ds.startCardFixedY + (e.clientY - ds.startPointerY),
      });
    };

    const onUp = (e: PointerEvent) => {
      const ds = dragStartRef.current;
      if (!ds) return;

      const zoneRect = zoneRef.current?.getBoundingClientRect();
      const finalX = ds.startCardFixedX + (e.clientX - ds.startPointerX);
      const finalY = ds.startCardFixedY + (e.clientY - ds.startPointerY);

      if (zoneRect) {
        const inZone =
          e.clientX >= zoneRect.left &&
          e.clientX <= zoneRect.right &&
          e.clientY >= zoneRect.top &&
          e.clientY <= zoneRect.bottom;

        if (inZone) {
          setCards((prev) =>
            prev.map((c) =>
              c.uuid === ds.uuid
                ? {
                    ...c,
                    inZone: true,
                    zoneX: finalX - zoneRect.left,
                    zoneY: finalY - zoneRect.top,
                  }
                : c,
            ),
          );
        } else if (ds.fromZone) {
          // ゾーン外に出したら元の位置（カード列 or 下）へ戻す
          setCards((prev) => prev.map((c) => (c.uuid === ds.uuid ? { ...c, inZone: false } : c)));
        }
      }

      dragStartRef.current = null;
      setDragState(null);
    };

    window.addEventListener("pointermove", onMove);
    window.addEventListener("pointerup", onUp);
    return () => {
      window.removeEventListener("pointermove", onMove);
      window.removeEventListener("pointerup", onUp);
    };
  }, []);

  const handleCardPointerDown = (
    uuid: string,
    fromZone: boolean,
    e: React.PointerEvent<HTMLDivElement>,
  ) => {
    e.preventDefault();
    const rect = e.currentTarget.getBoundingClientRect();
    dragStartRef.current = {
      uuid,
      fromZone,
      startPointerX: e.clientX,
      startPointerY: e.clientY,
      startCardFixedX: rect.left,
      startCardFixedY: rect.top,
    };
    setDragState({ uuid, fixedX: rect.left, fixedY: rect.top });
  };

  const draggingCard = dragState ? cards.find((c) => c.uuid === dragState.uuid) : null;
  const columnCards = cards.filter((c) => !c.inZone && !c.isOwn);
  const zoneCards = cards.filter((c) => c.inZone);
  const ownEntry = cards.find((c) => c.isOwn)!;
  const showOwnAtBottom = !ownEntry.inZone;

  const isDragging = (uuid: string) => dragState?.uuid === uuid;
  const canSeeResult = zoneCards.length === cards.length;

  const handleAlign = () => {
    const zone = zoneRef.current;
    if (!zone) return;
    const zoneW = zone.offsetWidth;
    const zoneH = zone.offsetHeight;
    const gap = 12;

    const cardW = (c: PlayCardEntry) => (c.isOwn ? 144 : 112);
    const cardH = (c: PlayCardEntry) => (c.isOwn ? 192 : 152);

    // 現在の x 位置でソートして順序を保持
    const sorted = [...zoneCards].sort((a, b) => a.zoneX - b.zoneX);
    const rowW = sorted.reduce((s, c) => s + cardW(c), 0) + (sorted.length - 1) * gap;
    let x = Math.max(0, (zoneW - rowW) / 2);

    const positions = new Map<string, { x: number; y: number }>();
    for (const c of sorted) {
      positions.set(c.uuid, { x, y: (zoneH - cardH(c)) / 2 });
      x += cardW(c) + gap;
    }

    setCards((prev) =>
      prev.map((c) => {
        const p = positions.get(c.uuid);
        return p ? { ...c, zoneX: p.x, zoneY: p.y } : c;
      }),
    );
    setIsAligned(true);
  };

  return (
    <div className="absolute inset-0 flex flex-col overflow-hidden bg-white">
      {/* 確定するボタン */}
      <div className="flex-shrink-0 px-4 pt-4">
        <button type="button" disabled={!canSeeResult} onClick={handleAlign} className={btn()}>
          確定する
        </button>
      </div>

      {/* 配置ゾーン（点線枠） */}
      <div
        ref={zoneRef}
        className="relative mx-4 mt-3 flex-shrink-0 border-2 border-dashed border-black"
        style={{ height: "36%" }}
      >
        {zoneCards.map((c) => (
          <div
            key={c.uuid}
            className="absolute cursor-grab touch-none"
            style={{
              left: c.zoneX,
              top: c.zoneY,
              transform: `rotate(${c.rotate}deg)`,
              visibility: isDragging(c.uuid) ? "hidden" : "visible",
            }}
            onPointerDown={(e) => handleCardPointerDown(c.uuid, true, e)}
          >
            <WordCard word={c.word} cardNumber={c.cardNumber} isOwn={c.isOwn} />
          </div>
        ))}
      </div>

      {/* カード列 or 確定後の送信エリア */}
      {isAligned ? (
        <div className="flex min-h-0 flex-1 flex-col items-center justify-center gap-4 px-4">
          <p className="text-base font-medium">回答しますか？</p>
          <button
            type="button"
            onClick={() => {
              const sorted = [...zoneCards].sort((a, b) => a.zoneX - b.zoneX);
              const answers: PlayAnswer[] = sorted.map((c, i) => ({
                uuid: c.uuid === OWN_UUID ? ownCardUuid : c.uuid,
                order: i + 1,
              }));
              onSubmit?.(answers);
            }}
            className={btn()}
          >
            送信
          </button>
        </div>
      ) : (
        <div className="flex min-h-0 flex-1 items-center overflow-x-auto overflow-y-hidden">
          <div className="mx-auto flex flex-row gap-4 px-4">
            {columnCards.map((c) => (
              <div
                key={c.uuid}
                className="flex-shrink-0 cursor-grab touch-none"
                style={{
                  transform: `rotate(${c.rotate}deg)`,
                  visibility: isDragging(c.uuid) ? "hidden" : "visible",
                }}
                onPointerDown={(e) => handleCardPointerDown(c.uuid, false, e)}
              >
                <WordCard word={c.word} />
              </div>
            ))}
          </div>
        </div>
      )}

      {/* 自分のカード（下固定） */}
      {showOwnAtBottom && (
        <div className="flex flex-shrink-0 justify-center pt-2 pb-8">
          <div
            className="cursor-grab touch-none"
            style={{ visibility: isDragging(OWN_UUID) ? "hidden" : "visible" }}
            onPointerDown={(e) => handleCardPointerDown(OWN_UUID, false, e)}
          >
            <WordCard word={ownEntry.word} cardNumber={ownEntry.cardNumber} isOwn />
          </div>
        </div>
      )}

      {/* ドラッグ中のゴーストカード */}
      {dragState && draggingCard && (
        <div
          className="pointer-events-none fixed z-50 opacity-90"
          style={{
            left: dragState.fixedX,
            top: dragState.fixedY,
            transform: `rotate(${draggingCard.rotate}deg)`,
          }}
        >
          <WordCard
            word={draggingCard.word}
            cardNumber={draggingCard.cardNumber}
            isOwn={draggingCard.isOwn}
          />
        </div>
      )}
    </div>
  );
};

// ---------------------------------------------------------------------------
// GamePanel
// ---------------------------------------------------------------------------

const MAX_CARD_AMOUNT = 10;

const GamePanel = ({ onBack }: { onBack: () => void }) => {
  const [step, setStep] = useState<GameStep>("config");
  const [cardAmount, setCardAmount] = useState(6);
  const [word, setWord] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [playData, setPlayData] = useState<{
    word: string;
    cardNumber: number;
    ownCardUuid: string;
    ownCardId: number;
  } | null>(null);
  const [playResult, setPlayResult] = useState<PostPlayRecordResponse | null>(null);
  const [editWord, setEditWord] = useState("");
  const [isEditingWord, setIsEditingWord] = useState(false);
  const [guestName, setGuestName] = useState("");

  const [session] = useAtom(sessionAtom);

  const { data: themes } = useThemes();
  const themeId = themes?.themes[0]?.id ?? null;

  const { data: probeCards, isLoading: probeLoading } = useGameCards(
    themeId ?? 0,
    MAX_CARD_AMOUNT,
    { enabled: themeId !== null },
  );

  const dbCardCount = probeCards?.cards.length ?? 0;
  const maxAmount = Math.min(MAX_CARD_AMOUNT, dbCardCount + 1);
  const canPlay = maxAmount >= 4;
  const clampedAmount = Math.max(4, Math.min(cardAmount, maxAmount));

  useEffect(() => {
    if (!probeLoading && canPlay && maxAmount === 4 && step === "config") {
      setStep("word");
    }
  }, [probeLoading, canPlay, maxAmount, step]);

  const { data: availableCard, isLoading: cardLoading } = useAvailableCard(themeId ?? 0, {
    enabled: step === "word" && themeId !== null,
  });

  // play フェーズ用: n-1 枚の他プレイヤーカード
  // backend は card_amount >= 4 を要求するため max(4, n-1) でリクエストし先頭 n-1 枚を使う
  const playFetchAmount = Math.max(4, clampedAmount - 1);
  const { data: playCardsData } = useGameCards(themeId ?? 0, playFetchAmount, {
    enabled: step === "play" && themeId !== null,
  });
  const otherCards = playCardsData?.cards.slice(0, clampedAmount - 1) ?? [];

  // カード枚数が確定した時点で傾き角度を生成（play 中は再生成しない）
  const scatterAngles = useMemo(
    () => Array.from({ length: clampedAmount - 1 }, () => Math.random() * 40 - 20),
    [clampedAmount],
  );

  const { mutateAsync: createCard, isPending: submitting } = useCreateCard();
  const { mutateAsync: play } = usePlay(session?.access_token ?? "");
  const { mutateAsync: confirmCard, isPending: confirmSubmitting } = useConfirmCard();

  const handleBack = () => {
    if (step === "word") {
      if (maxAmount === 4) {
        setWord("");
        setError(null);
        setStep("config");
        onBack();
      } else {
        setStep("config");
        setWord("");
        setError(null);
      }
    } else {
      setCardAmount(6);
      onBack();
    }
  };

  const handlePlaySubmit = async (answers: PlayAnswer[]) => {
    if (!themeId) return;
    try {
      const result = await play({ theme_id: themeId, answers });
      setPlayResult(result);
      setEditWord(playData?.word ?? "");
      setIsEditingWord(false);
      setStep("result");
    } catch {
      // エラーは結果なしで続行しない
    }
  };

  const handleConfirmWord = async () => {
    if (!playData) return;
    try {
      await confirmCard({
        id: playData.ownCardId,
        word: editWord,
        guest_name: session ? undefined : guestName.trim() || undefined,
        token: session?.access_token,
      });
      setStep("config");
      setWord("");
      setPlayData(null);
      setPlayResult(null);
      setEditWord("");
      setIsEditingWord(false);
      onBack();
    } catch {
      // エラー表示は今後の課題
    }
  };

  const handleSubmitWord = async () => {
    if (!themeId || !availableCard || !word.trim()) return;
    setError(null);
    try {
      const result = await createCard({
        themeId,
        card_number: availableCard.card_number,
        word: word.trim(),
        guest_name: session ? undefined : guestName.trim() || undefined,
        token: session?.access_token,
      });
      setPlayData({
        word: word.trim(),
        cardNumber: availableCard.card_number,
        ownCardUuid: result.uuid,
        ownCardId: result.id,
      });
      setStep("play");
    } catch {
      setError("登録に失敗しました");
    }
  };

  const stepLabel = {
    config: "ゲーム設定",
    word: "あなたの言葉",
    play: "プレイ中",
    result: "結果",
  }[step];
  const isPlay = step === "play";

  return (
    <div className="relative flex h-full flex-col bg-white">
      <div className={panelHeader({ hidden: isPlay })}>
        {!isPlay && (
          <button type="button" onClick={handleBack} className="font-medium">
            ←
          </button>
        )}
        <span className="font-medium">{stepLabel}</span>
      </div>

      {step === "config" && (
        <div className="flex flex-1 flex-col gap-8 overflow-y-auto p-6">
          <div className="flex flex-col gap-2">
            {probeLoading ? (
              <p className="text-sm opacity-40">読み込み中...</p>
            ) : !canPlay ? (
              <p className="text-sm">まだ遊べるカードが足りません</p>
            ) : (
              <>
                <label className="text-sm font-medium">枚数：{clampedAmount}枚</label>
                <input
                  type="range"
                  min={4}
                  max={maxAmount}
                  value={clampedAmount}
                  onChange={(e) => setCardAmount(Number(e.target.value))}
                  className="w-full accent-black"
                />
                <div className="flex justify-between text-xs opacity-40">
                  <span>4</span>
                  <span>{maxAmount}</span>
                </div>
              </>
            )}
          </div>
          <button
            type="button"
            onClick={() => setStep("word")}
            disabled={!canPlay || probeLoading}
            className={btn({ class: "mt-auto" })}
          >
            はじめる
          </button>
        </div>
      )}

      {step === "word" && (
        <div className="flex flex-1 flex-col gap-8 overflow-y-auto p-6">
          {themes?.themes[0]?.title && (
            <div className="self-start">
              <p className="text-xs opacity-40">お題</p>
              <p className="text-lg font-bold">{themes.themes[0].title}</p>
            </div>
          )}
          <div className="flex flex-col items-center gap-2">
            <p className="self-start text-sm font-medium">あなたの数字</p>
            <PlayingCard number={cardLoading ? null : (availableCard?.card_number ?? null)} />
          </div>
          {!session && (
            <div className="flex flex-col gap-2">
              <label className="text-sm font-medium">ゲスト名（10文字以内）</label>
              <input
                type="text"
                placeholder="名前を入力"
                value={guestName}
                onChange={(e) => setGuestName(e.target.value)}
                maxLength={10}
                className={textInput()}
              />
            </div>
          )}
          <div className="flex flex-col gap-2">
            <label className="text-sm font-medium">あなたの言葉（25文字以内）</label>
            <input
              type="text"
              placeholder="この数字を表す言葉を入力"
              value={word}
              onChange={(e) => setWord(e.target.value)}
              maxLength={25}
              className={textInput()}
            />
            {error && <p className="text-sm">{error}</p>}
          </div>
          <button
            type="button"
            onClick={handleSubmitWord}
            disabled={!word.trim() || cardLoading || submitting || (!session && !guestName.trim())}
            className={btn({ class: "mt-auto" })}
          >
            {submitting ? "登録中..." : "次へ"}
          </button>
        </div>
      )}

      {isPlay &&
        playData &&
        (otherCards.length > 0 ? (
          <PlayStep
            ownCard={playData}
            ownCardUuid={playData.ownCardUuid}
            gameCards={otherCards}
            scatterAngles={scatterAngles}
            onSubmit={handlePlaySubmit}
          />
        ) : (
          <div className="absolute inset-0 flex items-center justify-center">
            <p className="text-sm opacity-40">読み込み中...</p>
          </div>
        ))}

      {step === "result" && playResult && playData && (
        <div className="flex flex-1 flex-col gap-8 overflow-y-auto p-6">
          {/* 正答率 */}
          <div className="py-2 text-center">
            <p className="text-6xl font-bold">{Math.round(playResult.correct_rate)}%</p>
            <p className="mt-2 text-sm opacity-40">正答率</p>
          </div>

          {/* 正しい順番 */}
          <div className="flex flex-col gap-2">
            <p className="text-sm font-medium">正しい順番</p>
            <div className="overflow-x-auto pb-2">
              <div className="mx-auto flex w-fit flex-row gap-3 px-1">
                {[...playResult.cards]
                  .sort((a, b) => a.card_number - b.card_number)
                  .map((c) => (
                    <div key={c.uuid} className="flex-shrink-0">
                      <WordCard
                        word={c.word}
                        cardNumber={c.card_number}
                        isOwn={c.uuid === playData.ownCardUuid}
                      />
                    </div>
                  ))}
              </div>
            </div>
          </div>

          {/* 自分の言葉を確定 */}
          <div className="flex flex-col gap-3 border border-black p-4">
            <p className="text-sm font-medium">あなたの言葉</p>
            {isEditingWord ? (
              <input
                type="text"
                value={editWord}
                onChange={(e) => setEditWord(e.target.value)}
                maxLength={25}
                autoFocus
                className={textInput()}
              />
            ) : (
              <p className="text-lg font-bold">{editWord}</p>
            )}
            <p className="text-xs opacity-40">一度確定したら編集できません</p>
            <div className="flex gap-2">
              <button
                type="button"
                onClick={() => setIsEditingWord((v) => !v)}
                className={btn({ variant: "outline" })}
              >
                {isEditingWord ? "キャンセル" : "編集する"}
              </button>
              <button
                type="button"
                onClick={handleConfirmWord}
                disabled={confirmSubmitting || !editWord.trim()}
                className={btn()}
              >
                {confirmSubmitting ? "送信中..." : "確定する"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export { RootPage };
