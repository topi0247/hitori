# テーブル設計

## スキーマ変更手順

スキーマ定義は以下の2箇所を**必ず両方**更新する。

| ファイル | 用途 |
|---|---|
| `supabase/migrations/` | 実際のDB変更（Supabase が適用） |
| `back/db/schema/schema.sql` | sqlc の型推論用（コピー） |

```
1. supabase/migrations/ に新しい .sql ファイルを作成
2. back/db/schema/schema.sql を同じ内容に更新
3. back/db/query/ に必要なクエリを追加・変更
4. make generate/sqlc を実行（Go コード再生成）
```

> `back/db/schema/schema.sql` は `supabase/migrations/` の写し。  
> 将来的には sqlc の schema を直接 `supabase/migrations/` に向けることで一本化できる。

---

## profiles

```sql
id            BIGSERIAL     PK
auth_user_id  UUID          NOT NULL, UNIQUE, FK → auth.users.id
user_name     VARCHAR(10)   NOT NULL
created_at    TIMESTAMPTZ   NOT NULL, DEFAULT now()
updated_at    TIMESTAMPTZ   NOT NULL, DEFAULT now()
```

## themes

```sql
id            BIGSERIAL     PK
title         VARCHAR(100)  NOT NULL, UNIQUE
created_by    BIGINT        nullable, FK → profiles.id ON DELETE SET NULL
created_at    TIMESTAMPTZ   NOT NULL, DEFAULT now()
```

## cards

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
- 仮登録（`is_confirmed = false`）は有効期限24時間、pg_cronで1時間おきに自動削除

## play_records

```sql
id            BIGSERIAL     PK
theme_id      BIGINT        NOT NULL, FK → themes.id
profile_id    BIGINT        NOT NULL, FK → profiles.id
card_amount   SMALLINT      NOT NULL, CHECK(card_amount >= 4 AND card_amount <= 10)
correct_rate  NUMERIC(5,2)  NOT NULL, CHECK(correct_rate >= 0 AND correct_rate <= 100)
created_at    TIMESTAMPTZ   NOT NULL, DEFAULT now()
```

## マッチポイント計算

| 条件 | 点数 |
|---|---|
| 位置ぴったり | 3pt |
| 隣 | 2pt |
| ユーザー投票「納得！」（MVP第2弾以降） | 1pt |

プレイ時に即時計算し `cards.match_points` に加算。並べ替え結果は保持しない。
