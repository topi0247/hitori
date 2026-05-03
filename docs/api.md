# API仕様

## エラー種別

| エラー | 用途 |
|---|---|
| `ValidationError` | バリデーション失敗 |
| `NotFoundError` | 該当データなし |
| `ConflictError` | 数字の衝突など |
| `UnauthorizedError` | 認証切れ |
| `ForbiddenError` | 権限なし |
| `UnexpectedError` | 予期しないエラー（まとめて処理） |

エラーハンドリングは `@praha/byethrow` を使ったResult型で実装。

## プロフィール

### GET /profile
- 認証：必須
- レスポンス 200
```json
{ "user_name": "たろう" }
```
- エラー：401 `UnauthorizedError`

### PATCH /profile
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

### DELETE /profile
- 認証：必須
- レスポンス 200（bodyなし）
- エラー：401 `UnauthorizedError`

## テーマ

### GET /themes
- 認証：不要
- レスポンス 200
```json
{
  "themes": [
    { "id": 1, "title": "大きさ" }
  ]
}
```

## カード・ゲーム

### GET /themes/:id/cards/available
空きcard_numberをランダムに1件返す。
- 認証：不要
- レスポンス 200
```json
{ "card_number": 42 }
```
- エラー：404 `NotFoundError` / 409 `ConflictError`（空きなし）

### GET /themes/:id/cards/game
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

### POST /themes/:id/cards
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
  - `guest_name`：ゲスト時のみ必須
- バリデーション：`card_number` 1〜100 / `word` 25文字以内 / `guest_name` 10文字以内
- レスポンス 201
```json
{ "id": 1, "card_number": 42, "word": "アリ" }
```
- エラー：400 `ValidationError` / 404 `NotFoundError` / 409 `ConflictError`

### PATCH /cards/:id
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

### DELETE /cards/:id
仮登録カードの削除。
- 認証：不要（所有者確認）
- レスポンス 200（bodyなし）
- エラー：403 `ForbiddenError` / 404 `NotFoundError`

### POST /play_records
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
