CREATE TABLE IF NOT EXISTS qr_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    code_data TEXT NOT NULL UNIQUE,
    image_url TEXT DEFAULT '',
    status TEXT NOT NULL DEFAULT 'generated' CHECK (status IN ('generated', 'active', 'suspended', 'revoked')),
    issued_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ,
    last_scanned_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_qr_vehicle ON qr_codes(vehicle_id) WHERE status IN ('generated', 'active');
CREATE INDEX IF NOT EXISTS idx_qr_code_data ON qr_codes(code_data);
