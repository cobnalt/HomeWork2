CREATE TABLE IF NOT EXISTS currencyPair (
    id BIGSERIAL PRIMARY KEY,
    currency1 TEXT NOT NULL,
    currency2 TEXT NOT NULL,
    rate double precision,
    lastupdate time with time zone
);