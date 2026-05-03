# Feature Module: `analytics`

## Overview
Package analytics is a generated GoMock package.

## Types and Interfaces

### `AdminDashboardResponseDTO`
**Kind**: Struct

### `DashboardResponseDTO`
**Kind**: Struct

### `DemographicStat`
DemographicStat represents a pure business entity for analytics.

**Kind**: Struct

**Constructors/Factory Functions:**
- `MapDemographicStatToDomain`
- `MapDemographicStatsToDomain`

### `DemographicStatDB`
DemographicStatDB represents the database model for statistics.

**Kind**: Struct

### `DemographicStatDTO`
**Kind**: Struct

### `EducationalBackground`
**Kind**: Struct

### `EducationalBackgroundDB`
**Kind**: Struct

### `FamilyBackground`
**Kind**: Struct

### `FamilyBackgroundDB`
**Kind**: Struct

### `Handler`
**Kind**: Struct

**Methods:**
- `GetAdminDashboard`
- `GetDashboard`

**Constructors/Factory Functions:**
- `NewHandler`

### `MockRepositoryInterface`
MockRepositoryInterface is a mock of RepositoryInterface interface.

**Kind**: Struct

**Methods:**
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetAgeStats`
  - *GetAgeStats mocks base method.*
- `GetCityAddressStats`
  - *GetCityAddressStats mocks base method.*
- `GetCivilStatusStats`
  - *GetCivilStatusStats mocks base method.*
- `GetElementaryStats`
  - *GetElementaryStats mocks base method.*
- `GetFatherEducationStats`
  - *GetFatherEducationStats mocks base method.*
- `GetGenderStats`
  - *GetGenderStats mocks base method.*
- `GetHSGWAStats`
  - *GetHSGWAStats mocks base method.*
- `GetJuniorHighStats`
  - *GetJuniorHighStats mocks base method.*
- `GetMonthlyAppointmentStats`
  - *GetMonthlyAppointmentStats mocks base method.*
- `GetMonthlyIncomeStats`
  - *GetMonthlyIncomeStats mocks base method.*
- `GetMonthlyVisitorStats`
  - *GetMonthlyVisitorStats mocks base method.*
- `GetMotherEducationStats`
  - *GetMotherEducationStats mocks base method.*
- `GetNatureOfSchoolingStats`
  - *GetNatureOfSchoolingStats mocks base method.*
- `GetOrdinalPositionStats`
  - *GetOrdinalPositionStats mocks base method.*
- `GetParentsMaritalStatusStats`
  - *GetParentsMaritalStatusStats mocks base method.*
- `GetQuietStudyPlaceStats`
  - *GetQuietStudyPlaceStats mocks base method.*
- `GetReligionStats`
  - *GetReligionStats mocks base method.*
- `GetSeniorHighStats`
  - *GetSeniorHighStats mocks base method.*
- `GetTotalAppointments`
  - *GetTotalAppointments mocks base method.*
- `GetTotalReports`
  - *GetTotalReports mocks base method.*
- `GetTotalSlips`
  - *GetTotalSlips mocks base method.*
- `GetTotalStudents`
  - *GetTotalStudents mocks base method.*

**Constructors/Factory Functions:**
- `NewMockRepositoryInterface`

### `MockRepositoryInterfaceMockRecorder`
MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.

**Kind**: Struct

**Methods:**
- `GetAgeStats`
  - *GetAgeStats indicates an expected call of GetAgeStats.*
- `GetCityAddressStats`
  - *GetCityAddressStats indicates an expected call of GetCityAddressStats.*
- `GetCivilStatusStats`
  - *GetCivilStatusStats indicates an expected call of GetCivilStatusStats.*
- `GetElementaryStats`
  - *GetElementaryStats indicates an expected call of GetElementaryStats.*
- `GetFatherEducationStats`
  - *GetFatherEducationStats indicates an expected call of GetFatherEducationStats.*
- `GetGenderStats`
  - *GetGenderStats indicates an expected call of GetGenderStats.*
- `GetHSGWAStats`
  - *GetHSGWAStats indicates an expected call of GetHSGWAStats.*
- `GetJuniorHighStats`
  - *GetJuniorHighStats indicates an expected call of GetJuniorHighStats.*
- `GetMonthlyAppointmentStats`
  - *GetMonthlyAppointmentStats indicates an expected call of GetMonthlyAppointmentStats.*
- `GetMonthlyIncomeStats`
  - *GetMonthlyIncomeStats indicates an expected call of GetMonthlyIncomeStats.*
- `GetMonthlyVisitorStats`
  - *GetMonthlyVisitorStats indicates an expected call of GetMonthlyVisitorStats.*
- `GetMotherEducationStats`
  - *GetMotherEducationStats indicates an expected call of GetMotherEducationStats.*
- `GetNatureOfSchoolingStats`
  - *GetNatureOfSchoolingStats indicates an expected call of GetNatureOfSchoolingStats.*
- `GetOrdinalPositionStats`
  - *GetOrdinalPositionStats indicates an expected call of GetOrdinalPositionStats.*
- `GetParentsMaritalStatusStats`
  - *GetParentsMaritalStatusStats indicates an expected call of GetParentsMaritalStatusStats.*
- `GetQuietStudyPlaceStats`
  - *GetQuietStudyPlaceStats indicates an expected call of GetQuietStudyPlaceStats.*
- `GetReligionStats`
  - *GetReligionStats indicates an expected call of GetReligionStats.*
- `GetSeniorHighStats`
  - *GetSeniorHighStats indicates an expected call of GetSeniorHighStats.*
- `GetTotalAppointments`
  - *GetTotalAppointments indicates an expected call of GetTotalAppointments.*
- `GetTotalReports`
  - *GetTotalReports indicates an expected call of GetTotalReports.*
- `GetTotalSlips`
  - *GetTotalSlips indicates an expected call of GetTotalSlips.*
- `GetTotalStudents`
  - *GetTotalStudents indicates an expected call of GetTotalStudents.*

### `MockServiceInterface`
MockServiceInterface is a mock of ServiceInterface interface.

**Kind**: Struct

**Methods:**
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetAdminDashboard`
  - *GetAdminDashboard mocks base method.*
- `GetDashboard`
  - *GetDashboard mocks base method.*

**Constructors/Factory Functions:**
- `NewMockServiceInterface`

### `MockServiceInterfaceMockRecorder`
MockServiceInterfaceMockRecorder is the mock recorder for MockServiceInterface.

**Kind**: Struct

**Methods:**
- `GetAdminDashboard`
  - *GetAdminDashboard indicates an expected call of GetAdminDashboard.*
- `GetDashboard`
  - *GetDashboard indicates an expected call of GetDashboard.*

### `MonthlyVisitorStatDTO`
**Kind**: Struct

### `Repository`
**Kind**: Struct

**Methods:**
- `GetAgeStats`
- `GetCityAddressStats`
- `GetCivilStatusStats`
- `GetElementaryStats`
- `GetFatherEducationStats`
- `GetGenderStats`
- `GetHSGWAStats`
- `GetJuniorHighStats`
- `GetMonthlyAppointmentStats`
- `GetMonthlyIncomeStats`
- `GetMonthlyVisitorStats`
- `GetMotherEducationStats`
- `GetNatureOfSchoolingStats`
- `GetOrdinalPositionStats`
- `GetParentsMaritalStatusStats`
- `GetQuietStudyPlaceStats`
- `GetReligionStats`
- `GetSeniorHighStats`
- `GetTotalAppointments`
- `GetTotalReports`
- `GetTotalSlips`
- `GetTotalStudents`

**Constructors/Factory Functions:**
- `NewRepository`

### `RepositoryInterface`
**Kind**: Interface

### `Service`
**Kind**: Struct

**Methods:**
- `GetAdminDashboard`
- `GetDashboard`

**Constructors/Factory Functions:**
- `NewService`

### `ServiceInterface`
**Kind**: Interface

### `StudentFinances`
**Kind**: Struct

### `StudentFinancesDB`
**Kind**: Struct

### `StudentPersonalInfo`
StudentPersonalInfo represents pure student analytics data.

**Kind**: Struct

### `StudentPersonalInfoDB`
StudentPersonalInfoDB represents the database model for personal info.

**Kind**: Struct

## Package Level Functions

### `RegisterRoutes`
### `TestService_GetAdminDashboard`
### `TestService_GetDashboard`
