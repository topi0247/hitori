# バックエンド仕様

## 技術スタック

| ツール | 用途 |
|---|---|
| Go 1.26 | 言語 |
| Echo v5 | Webフレームワーク |
| pgx/v5 | PostgreSQLドライバ |
| sqlc | SQLからGoコードを生成 |
| golang-jwt/jwt v5 | JWT検証 |
| go.uber.org/mock | テスト用モック生成 |

## アーキテクチャ

DDD + クリーンアーキテクチャ。依存方向は外から内のみ。

```
handler → usecase → domain ← infra
```

| レイヤー | パッケージ | 役割 |
|---|---|---|
| handler | `handler/` | HTTPリクエスト/レスポンス変換 |
| middleware | `handler/middleware/` | JWT認証 |
| usecase | `usecase/` | ユースケース・ビジネスロジック |
| repository interface | `usecase/repository/` | データアクセスの抽象 |
| domain | `domain/` | エンティティ・ドメインロジック |
| infra | `infra/store/` | リポジトリ実装 |
| sqlc生成コード | `infra/db/sqlcgen/` | 自動生成（編集禁止） |

## ディレクトリ構成

```
back/
├── main.go
├── domain/
│   ├── card/
│   ├── game/
│   ├── profile/
│   └── theme/
├── usecase/
│   ├── repository/
│   │   └── mock/        # go generate で生成
│   └── *.usecase.go
├── handler/
│   ├── middleware/
│   └── *.handler.go
└── infra/
    ├── db/
    │   ├── query/       # SQLクエリ（sqlc入力）
    │   ├── schema/      # スキーマ（sqlc入力）
    │   └── sqlcgen/     # sqlc生成コード（編集禁止）
    └── store/           # リポジトリ実装
```

## 認証

Supabase AuthのJWTをそのまま検証する。SupabaseのJWT Secretを使用。

| ミドルウェア | 用途 |
|---|---|
| `middleware.JWT` | 認証必須。Authorizationヘッダーがない or 無効なら401 |
| `middleware.OptionalJWT` | 認証任意。有効なJWTがあればuser_idをコンテキストにセット |

JWTの `sub` クレームをauth_user_idとして使用する。

## 環境変数

| 変数名 | 説明 |
|---|---|
| `DATABASE_URL` | PostgreSQL接続URL |
| `JWT_SECRET` | SupabaseのJWT Secret |
| `ALLOWED_ORIGIN` | CORSで許可するオリジン（例：`http://localhost:5173`） |

## sqlc

SQLからタイプセーフなGoコードを生成するCLIツール（Goライブラリではない）。

- クエリ定義：`infra/db/query/*.sql`
- スキーマ定義：`infra/db/schema/*.sql`
- 生成先：`infra/db/sqlcgen/`（編集禁止）
- 設定ファイル：`back/sqlc.yml`

生成コマンド：
```bash
make generate/sqlc
```

### 型マッピング（overrides）

| DB型 | Go型 |
|---|---|
| `uuid` | `string` |
| `timestamptz` | `time.Time` |
| `pg_catalog.numeric` | `float64` |

NULL許容カラムは `pgtype.Int8` / `pgtype.Text` / `pgtype.Timestamptz` を使用し、`infra/store/convert.go` でGoのポインタ型に変換する。

## テスト方針

TDD（Red → Green）で実装。

- `domain/` 層：ユニットテスト
- `usecase/` 層：モックを使ったユニットテスト（`go generate` でモック生成）
- `infra/store/` 層：テストなし（DB接続が必要なため）
- `handler/` 層：テストなし

モック生成：
```bash
make generate
```

## ローカル起動

```bash
make up        # Supabase + フロント + バックを一括起動
make dev/back  # バックのみ起動
```

`.env` ファイルを `back/` に作成（`.env.example` を参考）。
