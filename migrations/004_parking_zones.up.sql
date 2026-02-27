CREATE TABLE IF NOT EXISTS parking_zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lot_id UUID NOT NULL REFERENCES parking_lots(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    capacity INTEGER NOT NULL DEFAULT 0,
    zone_type TEXT NOT NULL DEFAULT 'standard' CHECK (zone_type IN ('standard', 'oversized', 'hazmat', 'refrigerated')),
    current_occupancy INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_zones_lot ON parking_zones(lot_id);
