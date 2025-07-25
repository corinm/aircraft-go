// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: aircraft.sql

package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createAircraft = `-- name: CreateAircraft :exec
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
    plane_alert_db_tags = EXCLUDED.plane_alert_db_tags
`

type CreateAircraftParams struct {
	IcaoHexCode          string
	Registration         pgtype.Text
	Manufacturer         pgtype.Text
	IcaoTypeCode         pgtype.Text
	AircraftType         pgtype.Text
	RegisteredOwners     pgtype.Text
	IcaoAirlineCode      pgtype.Text
	Cmpg                 Cmpg
	PlaneAlertDbCategory pgtype.Text
	PlaneAlertDbTags     []string
}

func (q *Queries) CreateAircraft(ctx context.Context, arg CreateAircraftParams) error {
	_, err := q.db.Exec(ctx, createAircraft,
		arg.IcaoHexCode,
		arg.Registration,
		arg.Manufacturer,
		arg.IcaoTypeCode,
		arg.AircraftType,
		arg.RegisteredOwners,
		arg.IcaoAirlineCode,
		arg.Cmpg,
		arg.PlaneAlertDbCategory,
		arg.PlaneAlertDbTags,
	)
	return err
}

const getAircraft = `-- name: GetAircraft :one
SELECT icao_hex_code, registration, manufacturer, icao_type_code, aircraft_type, registered_owners, icao_airline_code, cmpg, plane_alert_db_category, plane_alert_db_tags FROM aircraft
WHERE icao_hex_code = $1 LIMIT 1
`

func (q *Queries) GetAircraft(ctx context.Context, icaoHexCode string) (Aircraft, error) {
	row := q.db.QueryRow(ctx, getAircraft, icaoHexCode)
	var i Aircraft
	err := row.Scan(
		&i.IcaoHexCode,
		&i.Registration,
		&i.Manufacturer,
		&i.IcaoTypeCode,
		&i.AircraftType,
		&i.RegisteredOwners,
		&i.IcaoAirlineCode,
		&i.Cmpg,
		&i.PlaneAlertDbCategory,
		&i.PlaneAlertDbTags,
	)
	return i, err
}

const listAircraft = `-- name: ListAircraft :many
SELECT icao_hex_code, registration, manufacturer, icao_type_code, aircraft_type, registered_owners, icao_airline_code, cmpg, plane_alert_db_category, plane_alert_db_tags FROM aircraft
ORDER BY icao_hex_code
`

func (q *Queries) ListAircraft(ctx context.Context) ([]Aircraft, error) {
	rows, err := q.db.Query(ctx, listAircraft)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Aircraft
	for rows.Next() {
		var i Aircraft
		if err := rows.Scan(
			&i.IcaoHexCode,
			&i.Registration,
			&i.Manufacturer,
			&i.IcaoTypeCode,
			&i.AircraftType,
			&i.RegisteredOwners,
			&i.IcaoAirlineCode,
			&i.Cmpg,
			&i.PlaneAlertDbCategory,
			&i.PlaneAlertDbTags,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
