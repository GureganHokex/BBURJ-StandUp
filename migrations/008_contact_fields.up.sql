ALTER TABLE site_settings
    ADD COLUMN IF NOT EXISTS contact_email VARCHAR(255),
    ADD COLUMN IF NOT EXISTS contact_phone VARCHAR(64),
    ADD COLUMN IF NOT EXISTS contact_telegram VARCHAR(128);

UPDATE site_settings
SET
    contact_email = COALESCE(NULLIF(contact_email, ''), 'booking@bburj.ru'),
    contact_phone = COALESCE(NULLIF(contact_phone, ''), '+7 (999) 000-00-00'),
    contact_telegram = COALESCE(NULLIF(contact_telegram, ''), '@bburj')
WHERE id = 1;
