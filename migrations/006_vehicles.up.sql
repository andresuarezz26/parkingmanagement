CREATE TABLE IF NOT EXISTS vehicles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    make TEXT NOT NULL DEFAULT '',
    model TEXT NOT NULL DEFAULT '',
    year TEXT DEFAULT '',
    plate_number TEXT DEFAULT '',
    vehicle_type TEXT NOT NULL DEFAULT 'other' CHECK (vehicle_type IN ('semi_truck', 'tanker', 'flatbed', 'construction', 'mining', 'other')),
    description TEXT DEFAULT '',
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'removed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_vehicles_account ON vehicles(account_id);
