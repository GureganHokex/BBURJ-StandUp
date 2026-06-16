ALTER TABLE site_settings
    DROP COLUMN IF EXISTS show_events,
    DROP COLUMN IF EXISTS show_videos,
    DROP COLUMN IF EXISTS show_photos,
    DROP COLUMN IF EXISTS show_merch,
    DROP COLUMN IF EXISTS show_about;
