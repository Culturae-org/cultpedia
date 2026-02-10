package models

type Country struct {
	Slug         string                       `json:"slug"`
	ISOAlpha2    string                       `json:"iso_alpha2"`
	ISOAlpha3    string                       `json:"iso_alpha3"`
	ISONumerics  string                       `json:"iso_numeric"`
	Name         map[string]string            `json:"name"`
	OfficialName map[string]string            `json:"official_name"`
	Capital      map[string]string            `json:"capital"`
	Continent    string                       `json:"continent"`
	Region       string                       `json:"region"`
	Coordinates  Coordinates                  `json:"coordinates"`
	Flag         string                       `json:"flag"`
	Population   int64                        `json:"population"`
	AreaKm2      float64                      `json:"area_km2"`
	Currency     Currency                     `json:"currency"`
	Languages    []string                     `json:"languages"`
	Neighbors    []string                     `json:"neighbors"`
	TLD          string                       `json:"tld"`
	PhoneCode    string                       `json:"phone_code"`
	DrivingSide  string                       `json:"driving_side"`
	Independent  bool                         `json:"independent"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Currency struct {
	Code   string            `json:"code"`
	Name   map[string]string `json:"name"`
	Symbol string            `json:"symbol"`
}

type Continent struct {
	Slug       string            `json:"slug"`
	Name       map[string]string `json:"name"`
	Countries  []string          `json:"countries"`
	AreaKm2    float64           `json:"area_km2"`
	Population int64             `json:"population"`
}

type Region struct {
	Slug      string            `json:"slug"`
	Name      map[string]string `json:"name"`
	Continent string            `json:"continent"`
	Countries []string          `json:"countries"`
}

type GeographyManifest struct {
	SchemaVersion string              `json:"schema_version"`
	Dataset       string              `json:"dataset"`
	Version       string              `json:"version"`
	Type          string              `json:"type"`
	License       string              `json:"license"`
	CreatedAt     string              `json:"created_at"`
	UpdatedAt     string              `json:"updated_at"`
	Sources       []Source            `json:"sources"`
	Includes      []string            `json:"includes"`
	Assets        map[string]string   `json:"assets"`
	Counts        map[string]int      `json:"counts"`
	Checksums     map[string]string   `json:"checksums"`
}

type Source struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	License string `json:"license"`
}
