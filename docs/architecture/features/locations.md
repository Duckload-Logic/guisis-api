# Feature Module: `locations`

## Overview
Package locations is a generated GoMock package.

## Types and Interfaces

### `Address`
**Kind**: Struct

**Methods:**
- `ToPersistence`

### `AddressDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Address mappers*

### `AddressDTO`
**Kind**: Struct

### `Barangay`
**Kind**: Struct

**Methods:**
- `ToPersistence`

### `BarangayDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Barangay mappers*

### `City`
**Kind**: Struct

**Methods:**
- `ToPersistence`

### `CityDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *City mappers*

### `Handler`
**Kind**: Struct

**Methods:**
- `GetBarangaysByCity`
- `GetCitiesByProvince`
- `GetCitiesByRegion`
- `GetProvincesByRegion`
- `GetRegions`

**Constructors/Factory Functions:**
- `NewHandler`

### `MockRepositoryInterface`
MockRepositoryInterface is a mock of RepositoryInterface interface.

**Kind**: Struct

**Methods:**
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetAddressByID`
  - *GetAddressByID mocks base method.*
- `GetBarangayByCode`
  - *GetBarangayByCode mocks base method.*
- `GetBarangaysByCity`
  - *GetBarangaysByCity mocks base method.*
- `GetCitiesByProvince`
  - *GetCitiesByProvince mocks base method.*
- `GetCitiesByRegion`
  - *GetCitiesByRegion mocks base method.*
- `GetCityByCode`
  - *GetCityByCode mocks base method.*
- `GetDB`
  - *GetDB mocks base method.*
- `GetProvinceByCode`
  - *GetProvinceByCode mocks base method.*
- `GetProvincesByRegion`
  - *GetProvincesByRegion mocks base method.*
- `GetRegionByCode`
  - *GetRegionByCode mocks base method.*
- `GetRegions`
  - *GetRegions mocks base method.*
- `UpsertAddress`
  - *UpsertAddress mocks base method.*
- `WithTransaction`
  - *WithTransaction mocks base method.*

**Constructors/Factory Functions:**
- `NewMockRepositoryInterface`

### `MockRepositoryInterfaceMockRecorder`
MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.

**Kind**: Struct

**Methods:**
- `GetAddressByID`
  - *GetAddressByID indicates an expected call of GetAddressByID.*
- `GetBarangayByCode`
  - *GetBarangayByCode indicates an expected call of GetBarangayByCode.*
- `GetBarangaysByCity`
  - *GetBarangaysByCity indicates an expected call of GetBarangaysByCity.*
- `GetCitiesByProvince`
  - *GetCitiesByProvince indicates an expected call of GetCitiesByProvince.*
- `GetCitiesByRegion`
  - *GetCitiesByRegion indicates an expected call of GetCitiesByRegion.*
- `GetCityByCode`
  - *GetCityByCode indicates an expected call of GetCityByCode.*
- `GetDB`
  - *GetDB indicates an expected call of GetDB.*
- `GetProvinceByCode`
  - *GetProvinceByCode indicates an expected call of GetProvinceByCode.*
- `GetProvincesByRegion`
  - *GetProvincesByRegion indicates an expected call of GetProvincesByRegion.*
- `GetRegionByCode`
  - *GetRegionByCode indicates an expected call of GetRegionByCode.*
- `GetRegions`
  - *GetRegions indicates an expected call of GetRegions.*
- `UpsertAddress`
  - *UpsertAddress indicates an expected call of UpsertAddress.*
- `WithTransaction`
  - *WithTransaction indicates an expected call of WithTransaction.*

### `MockServiceInterface`
MockServiceInterface is a mock of ServiceInterface interface.

**Kind**: Struct

**Methods:**
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetAddressByID`
  - *GetAddressByID mocks base method.*
- `GetBarangayByCode`
  - *GetBarangayByCode mocks base method.*
- `GetBarangaysByCity`
  - *GetBarangaysByCity mocks base method.*
- `GetCitiesByProvince`
  - *GetCitiesByProvince mocks base method.*
- `GetCitiesByRegion`
  - *GetCitiesByRegion mocks base method.*
- `GetCityByCode`
  - *GetCityByCode mocks base method.*
- `GetProvinceByCode`
  - *GetProvinceByCode mocks base method.*
- `GetProvincesByRegion`
  - *GetProvincesByRegion mocks base method.*
- `GetRegionByCode`
  - *GetRegionByCode mocks base method.*
- `GetRegions`
  - *GetRegions mocks base method.*
- `SaveAddress`
  - *SaveAddress mocks base method.*
- `WithTransaction`
  - *WithTransaction mocks base method.*

**Constructors/Factory Functions:**
- `NewMockServiceInterface`

### `MockServiceInterfaceMockRecorder`
MockServiceInterfaceMockRecorder is the mock recorder for MockServiceInterface.

**Kind**: Struct

**Methods:**
- `GetAddressByID`
  - *GetAddressByID indicates an expected call of GetAddressByID.*
- `GetBarangayByCode`
  - *GetBarangayByCode indicates an expected call of GetBarangayByCode.*
- `GetBarangaysByCity`
  - *GetBarangaysByCity indicates an expected call of GetBarangaysByCity.*
- `GetCitiesByProvince`
  - *GetCitiesByProvince indicates an expected call of GetCitiesByProvince.*
- `GetCitiesByRegion`
  - *GetCitiesByRegion indicates an expected call of GetCitiesByRegion.*
- `GetCityByCode`
  - *GetCityByCode indicates an expected call of GetCityByCode.*
- `GetProvinceByCode`
  - *GetProvinceByCode indicates an expected call of GetProvinceByCode.*
- `GetProvincesByRegion`
  - *GetProvincesByRegion indicates an expected call of GetProvincesByRegion.*
- `GetRegionByCode`
  - *GetRegionByCode indicates an expected call of GetRegionByCode.*
- `GetRegions`
  - *GetRegions indicates an expected call of GetRegions.*
- `SaveAddress`
  - *SaveAddress indicates an expected call of SaveAddress.*
- `WithTransaction`
  - *WithTransaction indicates an expected call of WithTransaction.*

### `Province`
**Kind**: Struct

**Methods:**
- `ToPersistence`

### `ProvinceDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Province mappers*

### `ProvinceDTO`
**Kind**: Struct

### `Region`
**Kind**: Struct

**Methods:**
- `ToPersistence`

### `RegionDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Region mappers*

### `Repository`
**Kind**: Struct

**Methods:**
- `GetAddressByID`
- `GetBarangayByCode`
- `GetBarangaysByCity`
- `GetCitiesByProvince`
- `GetCitiesByRegion`
- `GetCityByCode`
- `GetDB`
- `GetProvinceByCode`
- `GetProvincesByRegion`
- `GetRegionByCode`
- `GetRegions`
- `UpsertAddress`
- `WithTransaction`

**Constructors/Factory Functions:**
- `NewRepository`

### `RepositoryInterface`
RepositoryInterface defines the data access layer for location-based
operations.

**Kind**: Interface

### `Service`
**Kind**: Struct

**Methods:**
- `GetAddressByID`
- `GetBarangayByCode`
- `GetBarangaysByCity`
- `GetCitiesByProvince`
- `GetCitiesByRegion`
- `GetCityByCode`
- `GetProvinceByCode`
- `GetProvincesByRegion`
- `GetRegionByCode`
- `GetRegions`
- `SaveAddress`
- `WithTransaction`

**Constructors/Factory Functions:**
- `NewService`

### `ServiceInterface`
ServiceInterface defines the business logic for location-based operations.

**Kind**: Interface

## Package Level Functions

### `RegisterRoutes`
### `TestService_GetAddressByID`
### `TestService_GetRegions`
