CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id varchar(64) NOT NULL UNIQUE,
    email varchar(320) NOT NULL,
    password_hash varchar(255) NOT NULL,
    role varchar(32) NOT NULL DEFAULT 'customer',
    is_confirmed boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT users_role_not_empty CHECK (length(trim(role)) > 0)
);

CREATE UNIQUE INDEX users_email_lower_unique ON users (lower(email));

CREATE TABLE auth_tokens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash varchar(255) NOT NULL UNIQUE,
    purpose varchar(64) NOT NULL,
    expires_at timestamptz NOT NULL,
    used_at timestamptz,
    CONSTRAINT auth_tokens_purpose_not_empty CHECK (length(trim(purpose)) > 0)
);

CREATE INDEX auth_tokens_user_id_idx ON auth_tokens (user_id);
CREATE INDEX auth_tokens_expires_at_idx ON auth_tokens (expires_at);

CREATE TABLE admin_actions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_user_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    action_type varchar(100) NOT NULL,
    target_type varchar(100) NOT NULL,
    target_business_id varchar(64),
    details jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT admin_actions_action_type_not_empty CHECK (length(trim(action_type)) > 0),
    CONSTRAINT admin_actions_target_type_not_empty CHECK (length(trim(target_type)) > 0)
);

CREATE INDEX admin_actions_admin_user_id_idx ON admin_actions (admin_user_id);
CREATE INDEX admin_actions_created_at_idx ON admin_actions (created_at);
