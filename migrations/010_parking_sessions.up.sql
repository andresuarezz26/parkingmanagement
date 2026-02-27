CREATE TABLE IF NOT EXISTS parking_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vehicle_id UUID NOT NULL REFERENCES vehicles(id),
    lot_id UUID NOT NULL REFERENCES parking_lots(id),
    zone_id UUID REFERENCES parking_zones(id),
    entry_gate_id UUID NOT NULL REFERENCES gates(id),
    exit_gate_id UUID REFERENCES gates(id),
    subscription_id UUID REFERENCES subscriptions(id),
    entry_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    exit_time TIMESTAMPTZ,
    duration_minutes INTEGER,
    base_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    overstay_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'overstay', 'completed', 'disputed')),
    entry_scan_image_url TEXT,
    exit_scan_image_url TEXT
);

CREATE INDEX IF NOT EXISTS idx_sessions_vehicle ON parking_sessions(vehicle_id);
CREATE INDEX IF NOT EXISTS idx_sessions_lot ON parking_sessions(lot_id);
CREATE INDEX IF NOT EXISTS idx_sessions_status ON parking_sessions(status);
CREATE INDEX IF NOT EXISTS idx_sessions_entry ON parking_sessions(entry_time);
