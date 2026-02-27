CREATE TABLE IF NOT EXISTS gates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lot_id UUID NOT NULL REFERENCES parking_lots(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    hardware_id TEXT DEFAULT '',
    gate_type TEXT NOT NULL DEFAULT 'entry' CHECK (gate_type IN ('entry', 'exit', 'both')),
    status TEXT NOT NULL DEFAULT 'offline' CHECK (status IN ('online', 'offline', 'fault', 'maintenance')),
    ip_address TEXT DEFAULT '',
    last_heartbeat TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_gates_lot ON gates(lot_id);
