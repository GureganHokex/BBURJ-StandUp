ALTER TABLE site_settings
    DROP COLUMN IF EXISTS contact_email,
    DROP COLUMN IF EXISTS contact_phone,
    DROP COLUMN IF EXISTS contact_telegram;
