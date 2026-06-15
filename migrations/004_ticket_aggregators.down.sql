DROP INDEX IF EXISTS idx_events_ticket_source_external;

ALTER TABLE site_settings
    DROP COLUMN IF EXISTS event_import_keywords,
    DROP COLUMN IF EXISTS ticketscloud_api_key,
    DROP COLUMN IF EXISTS ticketscloud_org_id,
    DROP COLUMN IF EXISTS timepad_api_key,
    DROP COLUMN IF EXISTS timepad_org_id;

ALTER TABLE events
    DROP COLUMN IF EXISTS external_id,
    DROP COLUMN IF EXISTS ticket_source;
