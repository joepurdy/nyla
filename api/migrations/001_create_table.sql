CREATE TABLE IF NOT EXISTS events (
  site_id VARCHAR NOT NULL,
  created_at INT NOT NULL,
  type VARCHAR NOT NULL,
  event VARCHAR NOT NULL,
  referrer VARCHAR NOT NULL,
  is_touch BOOLEAN NOT NULL,
  browser_name VARCHAR NOT NULL,
  os_name VARCHAR NOT NULL,
  device_type VARCHAR NOT NULL,
  country VARCHAR NOT NULL,
  region VARCHAR NOT NULL,
  timestamp TIMESTAMPTZ DEFAULT now()
);