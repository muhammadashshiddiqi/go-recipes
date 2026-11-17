CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(100) NOT NULL,
    email         VARCHAR(150) NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS recipes (
    id                BIGSERIAL PRIMARY KEY,
    user_id           BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title             VARCHAR(150) NOT NULL,
    description       TEXT         NOT NULL DEFAULT '',
    ingredients       TEXT[]       NOT NULL DEFAULT '{}',
    steps             TEXT[]       NOT NULL DEFAULT '{}',
    cook_time_minutes INT          NOT NULL DEFAULT 0 CHECK (cook_time_minutes >= 0),
    servings          INT          NOT NULL DEFAULT 1 CHECK (servings > 0),
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_recipes_user_id ON recipes (user_id);
CREATE INDEX IF NOT EXISTS idx_recipes_title ON recipes (title);