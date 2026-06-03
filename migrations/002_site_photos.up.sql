CREATE TABLE IF NOT EXISTS site_settings (
    id INTEGER PRIMARY KEY DEFAULT 1,
    hero_image_url VARCHAR(512),
    portrait_image_url VARCHAR(512),
    hero_tagline VARCHAR(255),
    about_text TEXT,
    about_extra TEXT,
    youtube_url VARCHAR(512),
    telegram_url VARCHAR(512),
    instagram_url VARCHAR(512),
    youtube_handle VARCHAR(128),
    telegram_handle VARCHAR(128),
    instagram_handle VARCHAR(128),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT site_settings_singleton CHECK (id = 1)
);

CREATE TABLE IF NOT EXISTS photos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255),
    image_url VARCHAR(512) NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_photos_sort_order ON photos(sort_order);
