# Feature Module: `appointments`

## Overview
Package appointments is a generated GoMock package.

Package appointments is a generated GoMock package.

## Types and Interfaces

### `Appointment`
**Kind**: Struct

**Methods:**
- `ToPersistence`

### `AppointmentCategory`
**Kind**: Struct

### `AppointmentCategoryDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Category mappers*

### `AppointmentDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Appointment mappers*

### `AppointmentDTO`
**Kind**: Struct

### `AppointmentStatus`
**Kind**: Struct

### `AppointmentStatusDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Status mappers*

### `AppointmentWithDetailsView`
**Kind**: Struct

### `AppointmentWithDetailsViewDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Details View mapper*

### `AvailableTimeSlotView`
**Kind**: Struct

### `AvailableTimeSlotViewDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *AvailableTimeSlot mappers*

### `CancelAppointmentRequest`
PatchAppointment godoc
@Summary      Update appointment
@Description  Updates appointment details (reschedule).
@Tags         Appointments
@Accept       json
@Produce      json
@Param        id   path      int            true  "Appointment ID"
@Param        body body      AppointmentDTO true  "Updated appointment"
@Success      200  {object} map[string]string
@Failure      400  {object} map[string]string
@Failure      404  {object} map[string]string
@Failure      500  {object} map[string]string
@Router       /appointments/id/{id} [patch]

**Kind**: Struct

### `DailyStatusCount`
**Kind**: Struct

### `DailyStatusCountDB`
**Kind**: Struct

**Methods:**
- `ToDomain`

### `Handler`
**Kind**: Struct

**Methods:**
- `GetAppointmentByID`
  - *GetAppointmentByID godoc @Summary      Get appointment details @Description  Retrieves the details of a specific appointment. @Tags         Appointments @Accept       json @Produce      json @Param        id   path      int  true  "Appointment ID" @Success      200  {object}  Appointment @Failure      400  {object}  map[string]string @Failure      404  {object}  map[string]string @Failure      500  {object}  map[string]string @Router       /appointments/id/{id} [get]*
- `GetAppointmentCategories`
  - *GetAppointmentCategoryList godoc @Summary      Get appointment categories @Description  Retrieves available appointment categories. @Tags         Appointments @Produce      json @Success      200  {object} []AppointmentCategory @Failure      500  {object} map[string]string @Router       /appointments/lookups/categories [get] GetAppointmentCategoryList retrieves all appointment concern categories.*
- `GetAppointmentDailyStats`
  - *GetDailyStatusCountList godoc @Summary      Get daily status count @Description  Retrieves appointment status counts by date. @Tags         Appointments @Produce      json @Param        start_date query string true "Start date (YYYY-MM-DD)" @Success      200  {object} []DailyStatusCount @Failure      400  {object} map[string]string @Failure      500  {object} map[string]string @Router       /appointments/calendar/stats [get]*
- `GetAppointmentMe`
  - *GetAppointmentListByIIR godoc @Summary      Get student's appointments @Description  Retrieves appointments for the authenticated student. @Tags         Appointments @Produce      json @Success      200  {object} map[string]interface{} @Failure      403  {object} map[string]string @Failure      500  {object} map[string]string @Router       /appointments/me [get]*
- `GetAppointmentSlots`
  - *GetAvailableTimeSlotList godoc @Summary      Get available time slots @Description  Retrieves available time slots for a date. @Tags         Appointments @Produce      json @Param        date query string true "Date (YYYY-MM-DD)" @Success      200  {object} []AvailableTimeSlotView @Failure      400  {object} map[string]string @Failure      500  {object} map[string]string @Router       /appointments/lookups/slots [get]*
- `GetAppointmentStats`
  - *GetAppointmentStatsList godoc @Summary      Get appointment statistics @Description  Retrieves appointment status counts. @Tags         Appointments @Produce      json @Success      200  {object} []StatusCount @Failure      500  {object} map[string]string @Router       /appointments/stats [get]*
- `GetAppointmentStatuses`
  - *GetAppointmentStatusList godoc @Summary      Get appointment statuses @Description  Retrieves available appointment statuses. @Tags         Appointments @Produce      json @Success      200  {object} []AppointmentStatus @Failure      500  {object} map[string]string @Router       /appointments/lookups/statuses [get]*
- `GetAppointments`
  - *GetAppointmentList godoc @Summary      Get all appointments @Description  Retrieves a list of all appointments. @Tags         Appointments @Accept       json @Produce      json @Param        status     query string false "Filter by status" @Param        start_date query string false "Filter by start date" @Param        end_date   query string false "Filter by end date" @Success      200     {object}  map[string]interface{} @Failure      400     {object}  map[string]string @Failure      500     {object}  map[string]string @Router       /appointments [get]*
- `PatchAppointment`
  - *PatchAppointment godoc*
- `PostAppointment`
  - *PostAppointment godoc @Summary      Book a new appointment @Description  Schedules a new appointment for a student. @Tags         Appointments @Accept       json @Produce      json @Param        request body      AppointmentDTO true "Appointment Details" @Success      201     {object}  map[string]interface{} @Failure      400     {object}  map[string]string @Failure      403     {object}  map[string]string @Failure      500     {object}  map[string]string @Router       /appointments [post]*
- `PostAppointmentCancel`

**Constructors/Factory Functions:**
- `NewHandler`

### `ListAppointmentsDTO`
**Kind**: Struct

### `ListAppointmentsRequest`
**Kind**: Struct

### `MockRepositoryInterface`
MockRepositoryInterface is a mock of RepositoryInterface interface.

**Kind**: Struct

**Methods:**
- `CreateAppointment`
  - *CreateAppointment mocks base method.*
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetAppointment`
  - *GetAppointment mocks base method.*
- `GetAppointmentCategoryByID`
  - *GetAppointmentCategoryByID mocks base method.*
- `GetAppointmentStats`
  - *GetAppointmentStats mocks base method.*
- `GetAvailableTimeSlots`
  - *GetAvailableTimeSlots mocks base method.*
- `GetCategories`
  - *GetCategories mocks base method.*
- `GetDB`
  - *GetDB mocks base method.*
- `GetDailyStatusCount`
  - *GetDailyStatusCount mocks base method.*
- `GetStatusByID`
  - *GetStatusByID mocks base method.*
- `GetStatuses`
  - *GetStatuses mocks base method.*
- `GetTimeSlotByID`
  - *GetTimeSlotByID mocks base method.*
- `GetTimeSlots`
  - *GetTimeSlots mocks base method.*
- `GetTotalAppointmentsCount`
  - *GetTotalAppointmentsCount mocks base method.*
- `GetUserIDByAppointmentID`
  - *GetUserIDByAppointmentID mocks base method.*
- `List`
  - *List mocks base method.*
- `ListByIIRID`
  - *ListByIIRID mocks base method.*
- `ListByUserID`
  - *ListByUserID mocks base method.*
- `UpdateAppointment`
  - *UpdateAppointment mocks base method.*
- `WithTransaction`
  - *WithTransaction mocks base method.*

**Constructors/Factory Functions:**
- `NewMockRepositoryInterface`

### `MockRepositoryInterfaceMockRecorder`
MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.

**Kind**: Struct

**Methods:**
- `CreateAppointment`
  - *CreateAppointment indicates an expected call of CreateAppointment.*
- `GetAppointment`
  - *GetAppointment indicates an expected call of GetAppointment.*
- `GetAppointmentCategoryByID`
  - *GetAppointmentCategoryByID indicates an expected call of GetAppointmentCategoryByID.*
- `GetAppointmentStats`
  - *GetAppointmentStats indicates an expected call of GetAppointmentStats.*
- `GetAvailableTimeSlots`
  - *GetAvailableTimeSlots indicates an expected call of GetAvailableTimeSlots.*
- `GetCategories`
  - *GetCategories indicates an expected call of GetCategories.*
- `GetDB`
  - *GetDB indicates an expected call of GetDB.*
- `GetDailyStatusCount`
  - *GetDailyStatusCount indicates an expected call of GetDailyStatusCount.*
- `GetStatusByID`
  - *GetStatusByID indicates an expected call of GetStatusByID.*
- `GetStatuses`
  - *GetStatuses indicates an expected call of GetStatuses.*
- `GetTimeSlotByID`
  - *GetTimeSlotByID indicates an expected call of GetTimeSlotByID.*
- `GetTimeSlots`
  - *GetTimeSlots indicates an expected call of GetTimeSlots.*
- `GetTotalAppointmentsCount`
  - *GetTotalAppointmentsCount indicates an expected call of GetTotalAppointmentsCount.*
- `GetUserIDByAppointmentID`
  - *GetUserIDByAppointmentID indicates an expected call of GetUserIDByAppointmentID.*
- `List`
  - *List indicates an expected call of List.*
- `ListByIIRID`
  - *ListByIIRID indicates an expected call of ListByIIRID.*
- `ListByUserID`
  - *ListByUserID indicates an expected call of ListByUserID.*
- `UpdateAppointment`
  - *UpdateAppointment indicates an expected call of UpdateAppointment.*
- `WithTransaction`
  - *WithTransaction indicates an expected call of WithTransaction.*

### `MockServiceInterface`
MockServiceInterface is a mock of ServiceInterface interface.

**Kind**: Struct

**Methods:**
- `CreateAppointment`
  - *CreateAppointment mocks base method.*
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetAppointmentByID`
  - *GetAppointmentByID mocks base method.*
- `GetAppointmentStats`
  - *GetAppointmentStats mocks base method.*
- `GetAppointmentStatuses`
  - *GetAppointmentStatuses mocks base method.*
- `GetAppointmentsByIIRID`
  - *GetAppointmentsByIIRID mocks base method.*
- `GetAppointmentsByUserID`
  - *GetAppointmentsByUserID mocks base method.*
- `GetAvailableTimeSlots`
  - *GetAvailableTimeSlots mocks base method.*
- `GetConcernCategories`
  - *GetConcernCategories mocks base method.*
- `GetDailyStatusCount`
  - *GetDailyStatusCount mocks base method.*
- `GetUserIDByAppointmentID`
  - *GetUserIDByAppointmentID mocks base method.*
- `ListAppointments`
  - *ListAppointments mocks base method.*
- `UpdateAppointment`
  - *UpdateAppointment mocks base method.*

**Constructors/Factory Functions:**
- `NewMockServiceInterface`

### `MockServiceInterfaceMockRecorder`
MockServiceInterfaceMockRecorder is the mock recorder for MockServiceInterface.

**Kind**: Struct

**Methods:**
- `CreateAppointment`
  - *CreateAppointment indicates an expected call of CreateAppointment.*
- `GetAppointmentByID`
  - *GetAppointmentByID indicates an expected call of GetAppointmentByID.*
- `GetAppointmentStats`
  - *GetAppointmentStats indicates an expected call of GetAppointmentStats.*
- `GetAppointmentStatuses`
  - *GetAppointmentStatuses indicates an expected call of GetAppointmentStatuses.*
- `GetAppointmentsByIIRID`
  - *GetAppointmentsByIIRID indicates an expected call of GetAppointmentsByIIRID.*
- `GetAppointmentsByUserID`
  - *GetAppointmentsByUserID indicates an expected call of GetAppointmentsByUserID.*
- `GetAvailableTimeSlots`
  - *GetAvailableTimeSlots indicates an expected call of GetAvailableTimeSlots.*
- `GetConcernCategories`
  - *GetConcernCategories indicates an expected call of GetConcernCategories.*
- `GetDailyStatusCount`
  - *GetDailyStatusCount indicates an expected call of GetDailyStatusCount.*
- `GetUserIDByAppointmentID`
  - *GetUserIDByAppointmentID indicates an expected call of GetUserIDByAppointmentID.*
- `ListAppointments`
  - *ListAppointments indicates an expected call of ListAppointments.*
- `UpdateAppointment`
  - *UpdateAppointment indicates an expected call of UpdateAppointment.*

### `Repository`
**Kind**: Struct

**Methods:**
- `CreateAppointment`
- `GetAppointment`
- `GetAppointmentCategoryByID`
- `GetAppointmentStats`
- `GetAvailableTimeSlots`
- `GetCategories`
- `GetDB`
- `GetDailyStatusCount`
- `GetStatusByID`
- `GetStatuses`
- `GetTimeSlotByID`
- `GetTimeSlots`
- `GetTotalAppointmentsCount`
- `GetUserIDByAppointmentID`
- `List`
- `ListByIIRID`
- `ListByUserID`
- `UpdateAppointment`
- `WithTransaction`

### `RepositoryInterface`
**Kind**: Interface

**Constructors/Factory Functions:**
- `NewRepository`

### `Service`
**Kind**: Struct

**Methods:**
- `CreateAppointment`
- `GetAppointmentByID`
- `GetAppointmentStats`
- `GetAppointmentStatuses`
- `GetAppointmentsByIIRID`
- `GetAppointmentsByUserID`
- `GetAvailableTimeSlots`
- `GetConcernCategories`
- `GetDailyStatusCount`
- `GetUserIDByAppointmentID`
- `ListAppointments`
- `UpdateAppointment`
  - *handles Status updates AND Rescheduling*

### `ServiceInterface`
**Kind**: Interface

**Constructors/Factory Functions:**
- `NewService`

### `StatusCount`
**Kind**: Struct

### `StatusCountDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Status Count mappers*

### `TimeSlot`
**Kind**: Struct

**Methods:**
- `ToPersistence`

### `TimeSlotDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *TimeSlot mappers*

## Package Level Functions

### `RegisterRoutes`
### `TestService_CreateAppointment`
### `TestService_GetAppointmentByID`
### `TestService_UpdateAppointment`
