-- profiles
CREATE TABLE profiles (
  id            BIGSERIAL     PRIMARY KEY,
  auth_user_id  UUID          NOT NULL UNIQUE REFERENCES auth.users(id) ON DELETE CASCADE,
  user_name     VARCHAR(10)   NOT NULL,
  created_at    TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ   NOT NULL DEFAULT now()
);

-- themes
CREATE TABLE themes (
  id            BIGSERIAL     PRIMARY KEY,
  title         VARCHAR(100)  NOT NULL UNIQUE,
  created_by    BIGINT        REFERENCES profiles(id) ON DELETE SET NULL,
  created_at    TIMESTAMPTZ   NOT NULL DEFAULT now()
);

-- cards
CREATE TABLE cards (
  id            BIGSERIAL     PRIMARY KEY,
  uuid          UUID          NOT NULL UNIQUE DEFAULT gen_random_uuid(),
  theme_id      BIGINT        NOT NULL REFERENCES themes(id),
  profile_id    BIGINT        REFERENCES profiles(id),
  guest_name    VARCHAR(10),
  card_number   SMALLINT      NOT NULL CHECK (card_number >= 1 AND card_number <= 100),
  word          VARCHAR(25)   NOT NULL,
  is_confirmed  BOOLEAN       NOT NULL DEFAULT false,
  match_points  INTEGER       NOT NULL DEFAULT 0,
  expires_at    TIMESTAMPTZ,
  created_at    TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ   NOT NULL DEFAULT now(),
  UNIQUE (theme_id, card_number),
  UNIQUE (theme_id, guest_name),
  CHECK (profile_id IS NOT NULL OR guest_name IS NOT NULL)
);

-- play_records
CREATE TABLE play_records (
  id            BIGSERIAL     PRIMARY KEY,
  theme_id      BIGINT        NOT NULL REFERENCES themes(id),
  profile_id    BIGINT        NOT NULL REFERENCES profiles(id),
  card_amount   SMALLINT      NOT NULL CHECK (card_amount >= 4 AND card_amount <= 10),
  correct_rate  NUMERIC(5,2)  NOT NULL CHECK (correct_rate >= 0 AND correct_rate <= 100),
  created_at    TIMESTAMPTZ   NOT NULL DEFAULT now()
);

-- 期限切れ仮登録の自動削除（1時間おき）
SELECT cron.schedule(
  'delete_expired_cards',
  '0 * * * *',
  $$DELETE FROM cards WHERE is_confirmed = false AND expires_at < now()$$
);
