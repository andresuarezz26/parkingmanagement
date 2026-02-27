CREATE TABLE IF NOT EXISTS parking_lots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    operator_id UUID NOT NULL REFERENCES users(id),
    name TEXT NOT NULL,
    address TEXT DEFAULT '',
    city TEXT DEFAULT '',
    state TEXT DEFAULT '',
    latitude DOUBLE PRECISION DEFAULT 0,
    longitude DOUBLE PRECISION DEFAULT 0,
    total_capacity INTEGER NOT NULL DEFAULT 0,
    heavy_vehicle_spaces INTEGER NOT NULL DEFAULT 0,
    operating_hours JSONB DEFAULT '{}',
    status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'closed', 'maintenance')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_lots_operator ON parking_lots(operator_id);
