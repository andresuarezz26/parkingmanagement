CREATE TABLE IF NOT EXISTS incidents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lot_id UUID NOT NULL REFERENCES parking_lots(id),
    gate_id UUID REFERENCES gates(id),
    vehicle_id UUID REFERENCES vehicles(id),
    reported_by UUID REFERENCES users(id),
    type TEXT NOT NULL CHECK (type IN ('unauthorized_entry', 'tailgating', 'gate_fault', 'scan_failure', 'overstay')),
    description TEXT DEFAULT '',
    image_url TEXT,
    status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'investigating', 'resolved')),
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_incidents_lot ON incidents(lot_id);
CREATE INDEX IF NOT EXISTS idx_incidents_status ON incidents(status);
