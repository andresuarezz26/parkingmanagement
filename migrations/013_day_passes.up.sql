CREATE TABLE IF NOT EXISTS day_passes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lot_id UUID NOT NULL REFERENCES parking_lots(id),
    visitor_name TEXT NOT NULL DEFAULT '',
    visitor_phone TEXT DEFAULT '',
    visitor_email TEXT DEFAULT '',
    qr_code_id UUID REFERENCES qr_codes(id),
    valid_from TIMESTAMPTZ NOT NULL DEFAULT now(),
    valid_until TIMESTAMPTZ NOT NULL DEFAULT (now() + interval '24 hours'),
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'used', 'expired')),
    amount_paid DECIMAL(10,2) NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_day_passes_lot ON day_passes(lot_id);
