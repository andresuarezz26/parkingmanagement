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
DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE policyname = 'accounts_own') THEN
    CREATE POLICY accounts_own ON accounts FOR ALL USING (
      id IN (SELECT account_id FROM users WHERE id = auth.uid())
    );
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE policyname = 'vehicles_own') THEN
    CREATE POLICY vehicles_own ON vehicles FOR ALL USING (
      account_id IN (SELECT account_id FROM users WHERE id = auth.uid())
    );
  END IF;
END $$;
