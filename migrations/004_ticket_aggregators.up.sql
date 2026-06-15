ALTER TABLE events
    ADD COLUMN IF NOT EXISTS ticket_source VARCHAR(32) NOT NULL DEFAULT 'manual',
    ADD COLUMN IF NOT EXISTS external_id VARCHAR(128) NOT NULL DEFAULT '';

ALTER TABLE site_settings
    ADD COLUMN IF NOT EXISTS timepad_org_id VARCHAR(32) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS timepad_api_key VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS ticketscloud_org_id VARCHAR(128) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS ticketscloud_api_key VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS event_import_keywords VARCHAR(512) NOT NULL DEFAULT '';

CREATE UNIQUE INDEX IF NOT EXISTS idx_events_ticket_source_external
    ON events (ticket_source, external_id)
    WHERE external_id <> '' AND ticket_source <> 'manual';
