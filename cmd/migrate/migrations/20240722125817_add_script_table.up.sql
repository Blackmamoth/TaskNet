CREATE TABLE IF NOT EXISTS script (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    original_name STRING NOT NULL,
    parameters STRING[] NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now() ON UPDATE now(),

    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    UNIQUE (name)
);