CREATE TABLE stock_exchanges (
    id       BIGSERIAL    PRIMARY KEY,
    name     VARCHAR(255) NOT NULL,
    acronym  VARCHAR(20)  NOT NULL,
    mic_code VARCHAR(10)  NOT NULL UNIQUE,
    polity   VARCHAR(100) NOT NULL,
    currency VARCHAR(10)  NOT NULL,
    timezone VARCHAR(50)  NOT NULL
);
