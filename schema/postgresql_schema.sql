CREATE TABLE login (
  id text PRIMARY KEY,
  email text UNIQUE NOT NULL,
  salted_hash text NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC'),
  updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (NOW() AT TIME ZONE 'UTC')
);

CREATE INDEX login_email_idx ON login (email);
CREATE INDEX login_created_at_idx ON login (created_at);
CREATE INDEX login_updated_at_idx ON login (updated_at);

CREATE OR REPLACE FUNCTION set_updated_at()
  RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER login_set_updated_at_trg
  BEFORE UPDATE ON login
  FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();
