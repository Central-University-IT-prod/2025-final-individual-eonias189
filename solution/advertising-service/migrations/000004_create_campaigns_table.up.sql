CREATE TABLE IF NOT EXISTS campaigns (
    id UUID DEFAULT (gen_random_uuid()) PRIMARY KEY,
    advertiser_id UUID NOT NULL REFERENCES advertisers(id),
    impressions_limit INTEGER NOT NULL,
    clicks_limit INTEGER NOT NULL,
    cost_per_impression DOUBLE PRECISION NOT NULL,
    cost_per_click DOUBLE PRECISION NOT NULL,
    ad_title TEXT NOT NULL,
    ad_text TEXT NOT NULL,
    ad_image_url VARCHAR(255),
    start_date INTEGER NOT NULL,
    end_date INTEGER NOT NULL,
    gender VARCHAR(31),
    age_from INTEGER,
    age_to INTEGER,
    location TEXT,
    created_at TIMESTAMP DEFAULT (now())
);