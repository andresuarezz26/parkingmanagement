CREATE TABLE IF NOT EXISTS membership_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lot_id UUID NOT NULL REFERENCES parking_lots(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    billing_cycle TEXT NOT NULL CHECK (billing_cycle IN ('day', 'month', 'year', 'per_use')),
    price DECIMAL(10,2) NOT NULL DEFAULT 0,
    max_vehicles INTEGER NOT NULL DEFAULT 1,
    max_concurrent_entries INTEGER NOT NULL DEFAULT 1,
    grace_period_minutes INTEGER NOT NULL DEFAULT 30,
    overstay_rate_per_hour DECIMAL(10,2) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX IF NOT EXISTS idx_plans_lot ON membership_plans(lot_id);
