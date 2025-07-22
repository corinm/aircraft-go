-- +goose Up
-- +goose StatementBegin
CREATE TYPE cmpg AS ENUM ('Civilian', 'Military', 'Government', 'Police', 'Unknown');

CREATE TABLE aircraft (
    icao_hex_code VARCHAR(6) PRIMARY KEY,

    -- Apparently format and length vary by country. Allowing 10 chars to be on safe side
    registration VARCHAR(10) NULL,
    manufacturer VARCHAR(50) NULL,
    icao_type_code VARCHAR(4) NULL,
    aircraft_type VARCHAR(50) NULL,
    registered_owners VARCHAR(100) NULL,
    icao_airline_code VARCHAR(3) NULL,

    cmpg cmpg NOT NULL DEFAULT 'Unknown',

    plane_alert_db_category VARCHAR(50) NULL,
    plane_alert_db_tags TEXT[] CHECK (
        cardinality(plane_alert_db_tags) <= 3
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS aircraft;

DROP TYPE IF EXISTS cmpg;
-- +goose StatementEnd
