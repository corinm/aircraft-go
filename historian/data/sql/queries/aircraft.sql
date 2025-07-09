-- name: GetAircraft :one
SELECT * FROM aircraft
WHERE icao_hex_code = $1 LIMIT 1;

-- name: ListAircraft :many
SELECT * FROM aircraft
ORDER BY icao_hex_code;
