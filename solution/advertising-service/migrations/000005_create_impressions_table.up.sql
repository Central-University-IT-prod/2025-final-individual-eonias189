CREATE TABLE IF NOT EXISTS impressions (
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    client_id UUID NOT NULL REFERENCES clients(id),
    date INTEGER NOT NULL,
    profit DOUBLE PRECISION NOT NULL,
    UNIQUE (campaign_id, client_id)
);