# h「ito」ri 仕様まとめ

## アプリ概要

- **アプリ名：** h「ito」ri（読み：ひとり）
- **コンセプト：** ひとりでできるito
- **概要：** 意思疎通カードゲーム「ito」を1人で遊べるWebアプリケーション

---

## デザイン方針

- シンプル・ミニマム
- 白ベース、アクセントカラーは青系統
- 余計なデザインやアイコンは入れない
- モバイルファースト

---

## フェーズ

### MVP第1弾
- トップ画面
- 登録・ログイン
- ゲーム画面（1画面・ステップ形式）
  - 枚数指定（4〜10枚）
  - 自分の言葉設定
  - 並び替え（ドラッグ＆ドロップ、スマホ対応）
  - 結果
  - 自分の言葉確定
- アカウント設定画面
  - ユーザー名変更
  - 退会（退会してもカードは削除しない）
- マッチポイント機能（裏でロジックのみ・表示なし）

### MVP第2弾
- マッチポイント表示
- 再プレイによるカード編集機能
- マイページ（投稿したカード一覧）

### 最終リリース
- お題投稿機能
- 難易度拡張（MVP様子見で検討）
- マッチポイントユーザー投票（MVP様子見で検討）
- マイページにお題一覧

### 実装しない
- カード・お題の削除（開発者モデレーションのみ）

### いつかやりたい
- 言葉確定時のAI不適切判定
- 禁止お題・禁止言葉設定

---

## ゲームルール・仕様

### カード
- 数字範囲：1〜100（お題ごとに独立）
- 上限：100枚（100人）
- 初期投稿：開発者が設定（MVP時点）

### ゲームフロー
1. 空き数字がランダムに割り当てられる
2. その数字に対して言葉を設定（仮登録）
3. DB（3〜9枚）＋自分（1枚）＝4〜10枚で並べ替え
4. 答え確認
5. 言葉を編集 → 最終確定 → DB保存（本登録）

### マッチポイント
| 条件 | 点数 |
|---|---|
| 位置ぴったり | 3pt |
| 隣 | 2pt |
| ユーザー投票「納得！」（MVP第2弾以降） | 1pt |

マッチポイントはプレイ時に即時計算し `cards.match_points` に加算。並べ替え結果は保持しない。

### ユーザー
| | ログイン | ゲスト |
|---|---|---|
| 名前設定 | ✅ | ✅（自由入力） |
| カード投稿 | ✅ | ✅ |
| カード再編集 | ✅（再プレイ必須） | ❌ |
| 自分のカード一覧 | ✅ | ❌ |
| マッチポイント確認 | ✅ | ❌ |
| 後からログインで紐づけ | - | ✅（言葉確定前のみ） |

ゲスト投稿はDBに保存される。言葉確定後は誰であっても変更不可。

### 仮登録ルール
- 仮登録（`is_confirmed = false`）は数字衝突防止のためDBに保存
- 有効期限：仮登録から24時間（`expires_at`で管理）
- 本登録しなかった場合はフロントから `DELETE /cards/:id` で削除
- pg_cronで1時間おきに期限切れ仮登録を自動削除

---

## 画面構成

```
トップ画面
├─ 「遊ぶ」→ そのままトップ内でゲーム開始（画面遷移なし）
│     └─ 言葉確定前にログイン誘導モーダル
└─ 「登録・ログイン」→ 登録・ログイン画面
          └─ 完了後 → ゲーム画面（遷移）
```

| # | 画面 | 備考 |
|---|---|---|
| 1 | トップ＋ゲーム | シームレス・1画面 |
| 2 | 登録・ログイン | 別画面 |
| 3 | アカウント設定 | 別画面 |

---

## 技術スタック

| カテゴリ | 技術 |
|---|---|
| フロントエンド | React (TypeScript) / Vite |
| バックエンド | Go / Echo |
| 認証 | Supabase Auth |
| DB | Supabase (PostgreSQL) |
| コンテナ | Docker |
| CI/CD | GitHub Actions |
| ドメイン・DNS | Cloudflare |
| CDN・フロント配信 | CloudFront + S3 |
| APIサーバー | ECS Fargate + ECR |
| ロードバランサー | ALB |
| ネットワーク | VPC（パブリックサブネットのみ） |
| SSL証明書 | ACM |

### インフラ構成
```
ユーザー
  ↓
CloudFront → S3（React）
  ↓
ALB
  ↓
ECS Fargate（Go）
  ↓ インターネット経由
Supabase（Auth + PostgreSQL）
```

- 独自ドメインはCloudflareで取得・DNS管理
- サブドメインをCNAMEでCloudFront / ALBに向ける
- CloudflareのプロキシはOFF（グレーの雲マーク）

---

## DB設計

### profiles
```sql
id            BIGSERIAL     PK
auth_user_id  UUID          NOT NULL, UNIQUE, FK → auth.users.id
user_name     VARCHAR(10)   NOT NULL
created_at    TIMESTAMPTZ   NOT NULL, DEFAULT now()
updated_at    TIMESTAMPTZ   NOT NULL, DEFAULT now()
```

### themes
```sql
id            BIGSERIAL     PK
title         VARCHAR(100)  NOT NULL, UNIQUE
created_by    BIGINT        nullable, FK → profiles.id ON DELETE SET NULL
created_at    TIMESTAMPTZ   NOT NULL, DEFAULT now()
```

### cards
```sql
id            BIGSERIAL     PK
uuid          UUID          NOT NULL, UNIQUE, DEFAULT gen_random_uuid()
theme_id      BIGINT        NOT NULL, FK → themes.id
profile_id    BIGINT        nullable, FK → profiles.id
guest_name    VARCHAR(10)   nullable
card_number   SMALLINT      NOT NULL, CHECK(card_number >= 1 AND card_number <= 100)
word          VARCHAR(25)   NOT NULL
is_confirmed  BOOLEAN       NOT NULL, DEFAULT false
match_points  INTEGER       NOT NULL, DEFAULT 0
expires_at    TIMESTAMPTZ   nullable
created_at    TIMESTAMPTZ   NOT NULL, DEFAULT now()
updated_at    TIMESTAMPTZ   NOT NULL, DEFAULT now()

UNIQUE(theme_id, card_number)
UNIQUE(theme_id, guest_name)
CHECK(profile_id IS NOT NULL OR guest_name IS NOT NULL)
```

- `id`：管理用（外部非公開）
- `uuid`：識別用（フロントとのやり取りに使用）
- 仮登録時：`expires_at = created_at + 24時間`
- 本登録時：`is_confirmed = true`、`expires_at = NULL`

### play_records
```sql
id            BIGSERIAL     PK
theme_id      BIGINT        NOT NULL, FK → themes.id
profile_id    BIGINT        NOT NULL, FK → profiles.id
card_amount   SMALLINT      NOT NULL, CHECK(card_amount >= 4 AND card_amount <= 10)
correct_rate  NUMERIC(5,2)  NOT NULL, CHECK(correct_rate >= 0 AND correct_rate <= 100)
created_at    TIMESTAMPTZ   NOT NULL, DEFAULT now()
```

---

## API設計

### エラー種別
| エラー | 用途 |
|---|---|
| `ValidationError` | バリデーション失敗 |
| `NotFoundError` | 該当データなし |
| `ConflictError` | 数字の衝突など |
| `UnauthorizedError` | 認証切れ |
| `ForbiddenError` | 権限なし |
| `UnexpectedError` | 予期しないエラー（まとめて処理） |

エラーハンドリングは `@praha/byethrow` を使ったResult型で実装。

### プロフィール

#### GET /profile
- 認証：必須
- レスポンス 200
```json
{ "user_name": "たろう" }
```
- エラー：401 `UnauthorizedError`

#### PATCH /profile
- 認証：必須
- リクエスト
```json
{ "user_name": "じろう" }
```
- バリデーション：`user_name` 必須・10文字以内
- レスポンス 200
```json
{ "user_name": "じろう" }
```
- エラー：400 `ValidationError` / 401 `UnauthorizedError`

#### DELETE /profile
- 認証：必須
- レスポンス 200（bodyなし）
- エラー：401 `UnauthorizedError`

---

### テーマ

#### GET /themes
- 認証：不要
- レスポンス 200
```json
{
  "themes": [
    { "id": 1, "title": "大きさ" }
  ]
}
```

---

### カード・ゲーム

#### GET /themes/:id/cards/available
空きcard_numberをランダムに1件返す。
- 認証：不要
- レスポンス 200
```json
{ "card_number": 42 }
```
- エラー：404 `NotFoundError`（テーマなし）/ 409 `ConflictError`（空きなし）

#### GET /themes/:id/cards/game
ゲーム用カードをランダムに3〜9枚返す（`is_confirmed = true` のみ）。
- 認証：不要
- クエリパラメータ：`card_amount: number`（4〜10）
- レスポンス 200
```json
{
  "cards": [
    { "uuid": "uuid-xxxx", "word": "アリ" }
  ]
}
```
- エラー：400 `ValidationError` / 404 `NotFoundError`

#### POST /themes/:id/cards
仮登録。
- 認証：不要
- リクエスト
```json
{
  "card_number": 42,
  "word": "アリ",
  "guest_name": "ゲスト太郎"
}
```
  - `guest_name`：ゲスト時のみ必須（ログイン済みはJWTから取得）
- バリデーション：`card_number` 1〜100 / `word` 25文字以内 / `guest_name` 10文字以内
- レスポンス 201
```json
{ "id": 1, "card_number": 42, "word": "アリ" }
```
  - レスポンスの `id` は管理用ID（本登録・削除に使用）
- エラー：400 `ValidationError` / 404 `NotFoundError` / 409 `ConflictError`

#### PATCH /cards/:id
本登録（word確定）。
- 認証：不要（所有者確認：ゲストは `guest_name` をリクエストに含める）
- リクエスト
```json
{ "word": "クジラ" }
```
- バリデーション：`word` 必須・25文字以内
- レスポンス 200
```json
{ "id": 1, "card_number": 42, "word": "クジラ" }
```
- エラー：400 `ValidationError` / 403 `ForbiddenError` / 404 `NotFoundError`

#### DELETE /cards/:id
仮登録カードの削除。
- 認証：不要（PATCH同様、所有者確認）
- レスポンス 200（bodyなし）
- エラー：403 `ForbiddenError` / 404 `NotFoundError`

#### POST /play_records
プレイ結果記録 ＋ `correct_rate` ・ `match_points` 計算。
- 認証：必須
- リクエスト
```json
{
  "theme_id": 1,
  "card_amount": 6,
  "answers": [
    { "uuid": "uuid-xxxx", "order": 1 },
    { "uuid": "uuid-yyyy", "order": 2 }
  ]
}
```
- レスポンス 201
```json
{
  "correct_rate": 83.33,
  "cards": [
    { "uuid": "uuid-xxxx", "card_number": 15, "word": "アリ", "is_correct": true },
    { "uuid": "uuid-yyyy", "card_number": 32, "word": "クジラ", "is_correct": false }
  ]
}
```
- エラー：400 `ValidationError` / 401 `UnauthorizedError` / 404 `NotFoundError`