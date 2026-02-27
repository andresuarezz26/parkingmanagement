-- ============================================
-- HeavyPark — Full Schema Migration
-- Run in Supabase Dashboard → SQL Editor
-- ============================================

-- ========== 001_accounts.up.sql ==========
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type TEXT NOT NULL DEFAULT 'individual' CHECK (type IN ('individual', 'company')),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT DEFAULT '',
    billing_address TEXT DEFAULT '',
    tax_id TEXT DEFAULT '',
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'cancelled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_accounts_email ON accounts(email);


-- ========== 002_users.up.sql ==========
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID REFERENCES accounts(id) ON DELETE SET NULL,
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL,
    phone TEXT DEFAULT '',
    password_hash TEXT DEFAULT '',
    role TEXT NOT NULL DEFAULT 'driver' CHECK (role IN ('account_holder', 'driver', 'operator', 'attendant', 'super_admin')),
    last_login TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_account ON users(account_id);


-- ========== 003_parking_lots.up.sql ==========
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


-- ========== 004_parking_zones.up.sql ==========
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


-- ========== 005_gates.up.sql ==========
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


-- ========== 006_vehicles.up.sql ==========
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


-- ========== 007_qr_codes.up.sql ==========
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


-- ========== 008_membership_plans.up.sql ==========
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


-- ========== 009_subscriptions.up.sql ==========
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    plan_id UUID NOT NULL REFERENCES membership_plans(id),
    lot_id UUID NOT NULL REFERENCES parking_lots(id),
    start_date DATE NOT NULL DEFAULT CURRENT_DATE,
    end_date DATE,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'active', 'suspended', 'expired', 'cancelled')),
    payment_method_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_subs_account ON subscriptions(account_id);
CREATE INDEX IF NOT EXISTS idx_subs_lot ON subscriptions(lot_id);


-- ========== 010_parking_sessions.up.sql ==========
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


-- ========== 011_payments.up.sql ==========
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id),
    session_id UUID REFERENCES parking_sessions(id),
    invoice_id UUID,
    amount DECIMAL(10,2) NOT NULL,
    currency TEXT NOT NULL DEFAULT 'COP',
    payment_type TEXT NOT NULL CHECK (payment_type IN ('subscription', 'session', 'overstay', 'day_pass')),
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'succeeded', 'failed', 'refunded')),
    gateway_transaction_id TEXT,
    processed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_payments_account ON payments(account_id);
CREATE INDEX IF NOT EXISTS idx_payments_session ON payments(session_id);


-- ========== 012_invoices.up.sql ==========
CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id),
    subscription_id UUID REFERENCES subscriptions(id),
    invoice_number TEXT NOT NULL UNIQUE,
    billing_period_start DATE NOT NULL,
    billing_period_end DATE NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL DEFAULT 0,
    tax DECIMAL(10,2) NOT NULL DEFAULT 0,
    total DECIMAL(10,2) NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'issued', 'paid', 'overdue', 'void')),
    pdf_url TEXT,
    issued_at TIMESTAMPTZ,
    due_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_invoices_account ON invoices(account_id);
-- Add FK from payments to invoices now that invoices table exists
ALTER TABLE payments ADD CONSTRAINT fk_payments_invoice FOREIGN KEY (invoice_id) REFERENCES invoices(id);


-- ========== 013_day_passes.up.sql ==========
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


-- ========== 014_incidents.up.sql ==========
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


-- ========== 015_device_tokens.up.sql ==========
CREATE TABLE IF NOT EXISTS device_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    platform TEXT NOT NULL DEFAULT 'android' CHECK (platform IN ('ios', 'android')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_device_tokens_unique ON device_tokens(user_id, token);


-- ========== 016_indexes_and_seed.up.sql ==========
-- Additional composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_sessions_vehicle_active ON parking_sessions(vehicle_id) WHERE status = 'active';
CREATE INDEX IF NOT EXISTS idx_qr_active ON qr_codes(status) WHERE status = 'active';
CREATE INDEX IF NOT EXISTS idx_subs_active ON subscriptions(account_id) WHERE status = 'active';

-- Enable RLS on all tables
ALTER TABLE accounts ENABLE ROW LEVEL SECURITY;
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE vehicles ENABLE ROW LEVEL SECURITY;
ALTER TABLE qr_codes ENABLE ROW LEVEL SECURITY;
ALTER TABLE subscriptions ENABLE ROW LEVEL SECURITY;
ALTER TABLE parking_sessions ENABLE ROW LEVEL SECURITY;
ALTER TABLE payments ENABLE ROW LEVEL SECURITY;
ALTER TABLE invoices ENABLE ROW LEVEL SECURITY;
ALTER TABLE day_passes ENABLE ROW LEVEL SECURITY;
ALTER TABLE device_tokens ENABLE ROW LEVEL SECURITY;

-- RLS Policies: users can only see their own account data
CREATE POLICY IF NOT EXISTS accounts_own ON accounts FOR ALL USING (
    id IN (SELECT account_id FROM users WHERE id = auth.uid())
);

CREATE POLICY IF NOT EXISTS vehicles_own ON vehicles FOR ALL USING (
    account_id IN (SELECT account_id FROM users WHERE id = auth.uid())
);

-- Operators/admins can see everything (applied via role check in middleware, not RLS)
-- Service role bypasses RLS by default in Supabase


