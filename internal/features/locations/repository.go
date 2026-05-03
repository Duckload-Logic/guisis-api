package locations

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetDB() *sqlx.DB {
	return r.db
}

func (r *Repository) WithTransaction(
	ctx context.Context,
	fn func(datastore.DB) error,
) error {
	return datastore.RunInTransaction(ctx, r.db, fn)
}

func (r *Repository) GetRegions(ctx context.Context) ([]Region, error) {
	query := "SELECT id, code, name FROM regions ORDER BY name"
	var regions []Region
	err := r.db.SelectContext(ctx, &regions, query)
	if err != nil {
		return nil, err
	}

	return regions, nil
}

func (r *Repository) GetProvincesByRegion(
	ctx context.Context,
	regionCode string,
) ([]Province, error) {
	query := `
		SELECT id, code, name, region_code
		FROM provinces
		WHERE region_code = ?
		ORDER BY name
	`
	var provinces []Province
	err := r.db.SelectContext(ctx, &provinces, query, regionCode)
	if err != nil {
		return nil, err
	}

	return provinces, nil
}

func (r *Repository) GetCitiesByProvince(
	ctx context.Context,
	provinceCode string,
) ([]City, error) {
	query := `
		SELECT id, code, name, province_code, type, zip_code, district, region_code
		FROM cities
		WHERE province_code = ?
		ORDER BY name
	`
	var cities []City
	err := r.db.SelectContext(ctx, &cities, query, provinceCode)
	if err != nil {
		return nil, err
	}

	return cities, nil
}

func (r *Repository) GetCitiesByRegion(
	ctx context.Context,
	regionCode string,
) ([]City, error) {
	query := `
		SELECT id, code, name, province_code, type, zip_code, district, region_code
		FROM cities
		WHERE region_code = ?
		ORDER BY name
	`
	var cities []City
	err := r.db.SelectContext(ctx, &cities, query, regionCode)
	if err != nil {
		return nil, err
	}

	return cities, nil
}

func (r *Repository) GetBarangaysByCity(
	ctx context.Context,
	cityCode string,
) ([]Barangay, error) {
	query := `
		SELECT id, code, name, city_code
		FROM barangays
		WHERE city_code = ?
		ORDER BY name
	`
	var barangays []Barangay
	err := r.db.SelectContext(ctx, &barangays, query, cityCode)
	if err != nil {
		return nil, err
	}

	return barangays, nil
}

func (r *Repository) GetAddressByID(
	ctx context.Context,
	addressID int,
) (*Address, error) {
	query := `
		SELECT id, region_code, province_code, city_code, barangay_code, street_detail, created_at, updated_at
		FROM addresses
		WHERE id = ?
	`
	var addr Address
	err := r.db.GetContext(ctx, &addr, query, addressID)
	if err != nil {
		return nil, err
	}
	return &addr, nil
}

func (r *Repository) GetCityByCode(
	ctx context.Context,
	cityCode string,
) (*City, error) {
	query := `
		SELECT id, code, name, province_code, type, zip_code, district, region_code
		FROM cities
		WHERE code = ?
	`
	var city City
	err := r.db.GetContext(ctx, &city, query, cityCode)
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func (r *Repository) GetRegionByCode(
	ctx context.Context,
	regionCode string,
) (*Region, error) {
	query := "SELECT id, code, name FROM regions WHERE code = ?"
	var region Region
	err := r.db.GetContext(ctx, &region, query, regionCode)
	if err != nil {
		return nil, err
	}
	return &region, nil
}

func (r *Repository) GetBarangayByCode(
	ctx context.Context,
	barangayCode string,
) (*Barangay, error) {
	query := "SELECT id, code, name, city_code FROM barangays WHERE code = ?"
	var brgy Barangay
	err := r.db.GetContext(ctx, &brgy, query, barangayCode)
	if err != nil {
		return nil, err
	}
	return &brgy, nil
}

func (r *Repository) GetProvinceByCode(
	ctx context.Context,
	provinceCode string,
) (*Province, error) {
	query := "SELECT id, code, name, region_code FROM provinces WHERE code = ?"
	var prov Province
	err := r.db.GetContext(ctx, &prov, query, provinceCode)
	if err != nil {
		return nil, err
	}
	return &prov, nil
}

func (r *Repository) UpsertAddress(
	ctx context.Context,
	tx datastore.DB,
	addr *Address,
) (int, error) {
	query := `
		INSERT INTO addresses (
			region_code, province_code, city_code, barangay_code, street_detail
		) VALUES (
			:region_code, :province_code, :city_code, :barangay_code, :street_detail
		) ON DUPLICATE KEY UPDATE
			region_code = VALUES(region_code),
			province_code = VALUES(province_code),
			city_code = VALUES(city_code),
			barangay_code = VALUES(barangay_code),
			street_detail = VALUES(street_detail),
			updated_at = NOW()
	`
	result, err := tx.NamedExecContext(ctx, query, addr)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert address: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get last insert ID for address: %w",
			err,
		)
	}

	return int(lastID), nil
}

