package locations

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// Region represents a geographical region.
type Region struct {
	ID   int    `db:"id"   json:"id"`
	Code string `db:"code" json:"code"`
	Name string `db:"name" json:"name"`
}

// Province represents a province within a region.
type Province struct {
	ID         int    `db:"id"          json:"id"`
	Code       string `db:"code"        json:"code"`
	Name       string `db:"name"        json:"name"`
	RegionCode string `db:"region_code" json:"regionCode"`
}

// City represents a city or municipality.
type City struct {
	ID           int                    `db:"id"            json:"id"`
	Code         string                 `db:"code"          json:"code"`
	Name         string                 `db:"name"          json:"name"`
	ProvinceCode structs.NullableString `db:"province_code" json:"provinceCode"`
	Type         structs.NullableString `db:"type"          json:"type"`
	ZipCode      structs.NullableString `db:"zip_code"      json:"zipCode"`
	District     structs.NullableString `db:"district"      json:"district"`
	RegionCode   structs.NullableString `db:"region_code"   json:"regionCode"`
}

// Barangay represents the smallest administrative division.
type Barangay struct {
	ID       int    `db:"id"        json:"id"`
	Code     string `db:"code"      json:"code"`
	Name     string `db:"name"      json:"name"`
	CityCode int    `db:"city_code" json:"cityCode"`
}

// Address represents a physical location linked to an entity.
type Address struct {
	ID           int                    `db:"id"            json:"id"`
	RegionCode   string                 `db:"region_code"   json:"regionCode"`
	ProvinceCode structs.NullableString `db:"province_code" json:"provinceCode"`
	CityCode     string                 `db:"city_code"     json:"cityCode"`
	BarangayCode string                 `db:"barangay_code" json:"barangayCode"`
	StreetDetail structs.NullableString `db:"street_detail" json:"streetDetail"`
	CreatedAt    time.Time              `db:"created_at"    json:"createdAt"`
	UpdatedAt    time.Time              `db:"updated_at"    json:"updatedAt"`
}
