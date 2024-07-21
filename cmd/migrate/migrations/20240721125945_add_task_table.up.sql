CREATE TYPE IF NOT EXISTS task_status AS ENUM ('pending', 'scheduled', 'processing', 'failed', 'processed');


CREATE TABLE IF NOT EXISTS tasks (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    status task_status DEFAULT 'pending',
    task_code STRING NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now() ON UPDATE now(),

    CONSTRAINT "primary" PRIMARY KEY (id ASC)
);