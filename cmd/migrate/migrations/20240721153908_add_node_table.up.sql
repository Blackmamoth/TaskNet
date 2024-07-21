CREATE TYPE IF NOT EXISTS health_status AS ENUM ('healthy', 'unhealthy', 'offline');

CREATE TABLE IF NOT EXISTS nodes (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    status health_status DEFAULT 'healthy',
    previous_heartbeat_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now() ON UPDATE now(),

    CONSTRAINT "primary" PRIMARY KEY (id ASC)
);