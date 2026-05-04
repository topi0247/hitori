-- sqlc 用の auth スキーマスタブ（Supabase が実際に管理する）
CREATE SCHEMA IF NOT EXISTS auth;
CREATE TABLE auth.users (
  id UUID PRIMARY KEY
);
