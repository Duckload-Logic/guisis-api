CREATE TABLE push_subscriptions (
  id          VARCHAR(36)   NOT NULL PRIMARY KEY,
  user_id     VARCHAR(36)   NOT NULL,
  endpoint    VARCHAR(512)  NOT NULL UNIQUE,
  p256dh_key  TEXT          NOT NULL,
  auth_key    TEXT          NOT NULL,
  created_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_push_sub_user (user_id)
);
