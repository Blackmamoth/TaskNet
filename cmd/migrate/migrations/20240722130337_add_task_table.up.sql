CREATE TYPE IF NOT EXISTS task_type AS ENUM ('script', 'command');

CREATE TYPE IF NOT EXISTS task_status AS ENUM ('unscheduled', 'scheduled', 'active', 'inactive');

CREATE TYPE IF NOT EXISTS task_execution_mode AS ENUM ('once', 'recurring');

CREATE TABLE IF NOT EXISTS task (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    type task_type NOT NULL,
    script_id UUID DEFAULT NULL,
    command STRING,
    status task_status DEFAULT 'unscheduled',
    execution_mode task_execution_mode DEFAULT 'once',
    previous_run TIMESTAMPTZ DEFAULT NULL,
    next_run TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now() ON UPDATE now(),

    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    UNIQUE (name, script_id)
);