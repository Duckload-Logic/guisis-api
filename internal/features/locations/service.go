package locations

import (
	"context"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetRegions(ctx context.Context) ([]Region, error) {
	return s.repo.GetRegions(ctx)
}

func (s *Service) GetProvincesByRegion(
	ctx context.Context,
	regionCode string,
) ([]Province, error) {
	return s.repo.GetProvincesByRegion(ctx, regionCode)
}

func (s *Service) GetCitiesByProvince(
	ctx context.Context,
	provinceCode string,
) ([]City, error) {
	return s.repo.GetCitiesByProvince(ctx, provinceCode)
}

func (s *Service) GetCitiesByRegion(
	ctx context.Context,
	regionCode string,
) ([]City, error) {
	return s.repo.GetCitiesByRegion(ctx, regionCode)
}

func (s *Service) GetBarangaysByCity(
	ctx context.Context,
	cityCode string,
) ([]Barangay, error) {
	return s.repo.GetBarangaysByCity(ctx, cityCode)
}

func (s *Service) GetAddressByID(
	ctx context.Context,
	addressID int,
) (AddressDTO, error) {
	addr, err := s.repo.GetAddressByID(ctx, addressID)
	if err != nil {
		return AddressDTO{}, err
	}

	barangay, err := s.repo.GetBarangayByCode(ctx, addr.BarangayCode)
	if err != nil {
		return AddressDTO{}, err
	}
	city, err := s.repo.GetCityByCode(ctx, addr.CityCode)
	if err != nil {
		return AddressDTO{}, err
	}
	region, err := s.repo.GetRegionByCode(ctx, addr.RegionCode)
	if err != nil {
		return AddressDTO{}, err
	}

	addrDTO := AddressDTO{
		ID:           addr.ID,
		StreetDetail: addr.StreetDetail.String,
		Barangay:     *barangay,
		City:         *city,
		Region:       *region,
	}

	// Province is optional (e.g. NCR has no provinces)
	if addr.ProvinceCode.Valid && addr.ProvinceCode.String != "" {
		province, err := s.repo.GetProvinceByCode(ctx, addr.ProvinceCode.String)
		if err != nil {
			return AddressDTO{}, err
		}
		addrDTO.Province = province
	}

	return addrDTO, nil
}

func (s *Service) GetCityByCode(
	ctx context.Context,
	cityCode string,
) (*City, error) {
	return s.repo.GetCityByCode(ctx, cityCode)
}

func (s *Service) GetRegionByCode(
	ctx context.Context,
	regionCode string,
) (*Region, error) {
	return s.repo.GetRegionByCode(ctx, regionCode)
}

func (s *Service) GetBarangayByCode(
	ctx context.Context,
	barangayCode string,
) (*Barangay, error) {
	return s.repo.GetBarangayByCode(ctx, barangayCode)
}

func (s *Service) GetProvinceByCode(
	ctx context.Context,
	provinceCode string,
) (*Province, error) {
	return s.repo.GetProvinceByCode(ctx, provinceCode)
}

func (s *Service) SaveAddress(
	ctx context.Context,
	tx datastore.DB,
	addr *Address,
) (int, error) {
	return s.repo.UpsertAddress(ctx, tx, addr)
}

func (s *Service) UpsertAddress(
	ctx context.Context,
	tx datastore.DB,
	dto AddressDTO,
) (int, error) {
	addr := s.mapDTOToAddress(dto)
	return s.repo.UpsertAddress(ctx, tx, &addr)
}

func (s *Service) mapDTOToAddress(dto AddressDTO) Address {
	addr := Address{
		ID:           dto.ID,
		RegionCode:   dto.Region.Code,
		CityCode:     dto.City.Code,
		BarangayCode: dto.Barangay.Code,
		StreetDetail: structs.StringToNullableString(dto.StreetDetail),
	}

	if dto.Province != nil {
		addr.ProvinceCode = structs.StringToNullableString(dto.Province.Code)
	}

	return addr
}
