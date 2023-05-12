CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS currency_rates (
    id          UUID            DEFAULT uuid_generate_v4(),
    title       TEXT                                        NOT NULL,
    code        TEXT                                        NOT NULL,
    rate        NUMERIC(18, 2)                              NOT NULL,
    quant       INTEGER         DEFAULT 1,
    change      TEXT                                        NOT NULL, 
    valid_at    DATE                                        NOT NULL,
    PRIMARY KEY(id),
    UNIQUE (valid_at, code)
);