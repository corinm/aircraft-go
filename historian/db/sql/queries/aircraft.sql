-- name: GetAircraft :one
SELECT * FROM aircraft
WHERE icao_hex_code = $1 LIMIT 1;

-- name: ListAircraft :many
SELECT * FROM aircraft
ORDER BY icao_hex_code;

-- name: CreateAircraft :exec
INSERT INTO aircraft (icao_hex_code, registration, manufacturer, icao_type_code, aircraft_type, registered_owners, icao_airline_code, cmpg, plane_alert_db_category, plane_alert_db_tags)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (icao_hex_code) DO UPDATE SET
    registration = EXCLUDED.registration,
    manufacturer = EXCLUDED.manufacturer,
    icao_type_code = EXCLUDED.icao_type_code,
    aircraft_type = EXCLUDED.aircraft_type,
    registered_owners = EXCLUDED.registered_owners,
    icao_airline_code = EXCLUDED.icao_airline_code,
    cmpg = EXCLUDED.cmpg,
    plane_alert_db_category = EXCLUDED.plane_alert_db_category,
    plane_alert_db_tags = EXCLUDED.plane_alert_db_tags;
