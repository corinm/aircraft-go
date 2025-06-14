package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/corinm/aircraft/discovery/data"
)

type tar1090response struct {
	Now float64 `json:"now"`
	Messages int `json:"messages"`
	Aircraft []tar1090response_aircraft `json:"aircraft"`
}

type tar1090response_aircraft struct {
	Hex                string        `json:"hex"`
	Type               string        `json:"type"`
	Flight             *string       `json:"flight"`
	AltBaro            float64       `json:"alt_baro"`
	AltGeom            float64       `json:"alt_geom"`
	GS                 float64       `json:"gs"`
	IAS                float64       `json:"ias"`
	TAS                float64       `json:"tas"`
	Mach               float64       `json:"mach"`
	WD                 float64       `json:"wd"`
	WS                 float64       `json:"ws"`
	OAT                float64       `json:"oat"`
	TAT                float64       `json:"tat"`
	Track              float64       `json:"track"`
	TrackRate          float64       `json:"track_rate"`
	Roll               float64       `json:"roll"`
	MagHeading         float64       `json:"mag_heading"`
	TrueHeading        float64       `json:"true_heading"`
	BaroRate           float64       `json:"baro_rate"`
	GeomRate           float64       `json:"geom_rate"`
	Squawk             string        `json:"squawk"`
	Emergency          string        `json:"emergency"`
	Category           string        `json:"category"`
	NavQNH             float64       `json:"nav_qnh"`
	NavAltitudeMCP     float64       `json:"nav_altitude_mcp"`
	NavAltitudeFMS     float64       `json:"nav_altitude_fms"`
	NavHeading         float64       `json:"nav_heading"`
	Lat                float64       `json:"lat"`
	Lon                float64       `json:"lon"`
	NIC                float64       `json:"nic"`
	RC                 float64       `json:"rc"`
	SeenPos            float64       `json:"seen_pos"`
	RDst               float64       `json:"r_dst"`
	RDir               float64       `json:"r_dir"`
	Version            float64       `json:"version"`
	NICBaro            float64       `json:"nic_baro"`
	NACp               float64       `json:"nac_p"`
	NACv               float64       `json:"nac_v"`
	SIL                float64       `json:"sil"`
	SILType            string        `json:"sil_type"`
	GVA                float64       `json:"gva"`
	SDA                float64       `json:"sda"`
	Alert              float64       `json:"alert"`
	SPI                float64       `json:"spi"`
	MLAT               []interface{} `json:"mlat"`
	TISB               []interface{} `json:"tisb"`
	Messages           float64       `json:"messages"`
	Seen               float64       `json:"seen"`
	RSSI               float64       `json:"rssi"`
}

type Tar1090AdsbFetcher struct {
	URL string
}

func (f Tar1090AdsbFetcher) FetchAircraft() ([]data.Aircraft, error) {
	resp, err := http.Get(f.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("HTTP GET request to tar1090 API completed with status:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch tar1090 data, resp.Status: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	d := tar1090response{}

	if err := json.Unmarshal(body, &d); err != nil {
		return nil, err
	}

	aircraft := []data.Aircraft{}

	for _, a := range d.Aircraft {
		aircraft = append(aircraft, transformTar1090AircraftToAircraft(a))
	}

	fmt.Println("Processed", len(d.Aircraft), "aircraft from tar1090 response")

	return aircraft, nil
}

func transformTar1090AircraftToAircraft(a tar1090response_aircraft) data.Aircraft {
	return data.Aircraft{
		AiocHexCode: a.Hex,
	}
}
