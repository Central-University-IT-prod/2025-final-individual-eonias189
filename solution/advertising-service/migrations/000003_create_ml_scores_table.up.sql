CREATE TABLE IF NOT EXISTS ml_scores(
    client_id UUID NOT NULL REFERENCES clients(id),
    advertiser_id UUID NOT NULL REFERENCES advertisers(id),
    score INT NOT NULL,
    UNIQUE (client_id, advertiser_id)
);