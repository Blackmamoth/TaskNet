CREATE TYPE IF NOT EXISTS task_status AS ENUM ('unscheduled', 'scheduled', 'active', 'inactive');

CREATE TYPE IF NOT EXISTS task_execution_mode AS ENUM ('non-recurring', 'recurring');

CREATE TYPE IF NOT EXISTS task_priority AS ENUM ('low', 'medium', 'high');

CREATE TABLE IF NOT EXISTS task (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    status task_status DEFAULT 'unscheduled',
    execution_mode task_execution_mode DEFAULT 'non-recurring',
    priority task_priority NOT NULL,
    previous_run TIMESTAMPTZ DEFAULT NULL,
    next_run TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now() ON UPDATE now(),

    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    UNIQUE (name)
);