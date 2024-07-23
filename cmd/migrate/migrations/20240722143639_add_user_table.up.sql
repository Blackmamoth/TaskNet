CREATE TABLE IF NOT EXISTS "user" (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    username STRING NOT NULL,
    email STRING NOT NULL,
    password STRING NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now() ON UPDATE now(),

    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    UNIQUE (username, email)
);