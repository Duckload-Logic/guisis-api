# Feature Module: `logs`

## Overview
Package logs is a generated GoMock package.

## Types and Interfaces

### `Handler`
**Kind**: Struct

**Methods:**
- `GetLogs`
  - *HandleListSystemLogs godoc @Summary      List system logs @Description  Retrieves a paginated list of audit, system, and security logs. Super Admin only. @Tags         SystemLogs @Accept       json @Produce      json @Param        page        query     int    false "Page number" @Param        page_size   query     int    false "Number of entries per page" @Param        category    query     string false "Filter by category (AUDIT, SYSTEM, SECURITY)" @Param        action      query     string false "Filter by action" @Param        user_email  query     string false "Filter by user email" @Param        start_date  query     string false "Filter from date (YYYY-MM-DD)" @Param        end_date    query     string false "Filter to date (YYYY-MM-DD)" @Param        search      query     string false "Search in message, action, or user email" @Success      200         {object}  ListSystemLogsDTO @Failure      400         {object}  map[string]string "Bad request" @Failure      500         {object}  map[string]string "Internal server error" @Router       /system-logs [get]*
- `GetLogsActivity`
  - *GetActivityStats returns log counts grouped by hour for the last 24 hours*
- `GetLogsAudit`
- `GetLogsMe`
  - *GetMyLogs retrieves activity logs for the currently authenticated user.*
- `GetLogsSecurity`
  - *GetSecurityLogs returns only SECURITY category logs*
- `GetLogsStats`
  - *GetLogStats returns log counts by category*
- `GetLogsSystem`
- `GetService`

**Constructors/Factory Functions:**
- `NewHandler`

### `ListSystemLogsDTO`
ListSystemLogsDTO is the paginated response for system logs

**Kind**: Struct

### `ListSystemLogsRequest`
ListSystemLogsRequest holds query parameters for listing system logs

**Kind**: Struct

### `LogActivityDTO`
**Kind**: Struct

### `LogStatsDTO`
LogStatsDTO returns summary counts by category

**Kind**: Struct

### `MockRepositoryInterface`
MockRepositoryInterface is a mock of RepositoryInterface interface.

**Kind**: Struct

**Methods:**
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetActivityStats`
  - *GetActivityStats mocks base method.*
- `GetDB`
  - *GetDB mocks base method.*
- `GetStats`
  - *GetStats mocks base method.*
- `GetTotalCount`
  - *GetTotalCount mocks base method.*
- `List`
  - *List mocks base method.*
- `Record`
  - *Record mocks base method.*
- `WithTransaction`
  - *WithTransaction mocks base method.*

**Constructors/Factory Functions:**
- `NewMockRepositoryInterface`

### `MockRepositoryInterfaceMockRecorder`
MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.

**Kind**: Struct

**Methods:**
- `GetActivityStats`
  - *GetActivityStats indicates an expected call of GetActivityStats.*
- `GetDB`
  - *GetDB indicates an expected call of GetDB.*
- `GetStats`
  - *GetStats indicates an expected call of GetStats.*
- `GetTotalCount`
  - *GetTotalCount indicates an expected call of GetTotalCount.*
- `List`
  - *List indicates an expected call of List.*
- `Record`
  - *Record indicates an expected call of Record.*
- `WithTransaction`
  - *WithTransaction indicates an expected call of WithTransaction.*

### `MockServiceInterface`
MockServiceInterface is a mock of ServiceInterface interface.

**Kind**: Struct

**Methods:**
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetActivityStats`
  - *GetActivityStats mocks base method.*
- `GetDB`
  - *GetDB mocks base method.*
- `GetStats`
  - *GetStats mocks base method.*
- `ListLogs`
  - *ListLogs mocks base method.*
- `Record`
  - *Record mocks base method.*
- `RecordSecurity`
  - *RecordSecurity mocks base method.*
- `WithTransaction`
  - *WithTransaction mocks base method.*

**Constructors/Factory Functions:**
- `NewMockServiceInterface`

### `MockServiceInterfaceMockRecorder`
MockServiceInterfaceMockRecorder is the mock recorder for MockServiceInterface.

**Kind**: Struct

**Methods:**
- `GetActivityStats`
  - *GetActivityStats indicates an expected call of GetActivityStats.*
- `GetDB`
  - *GetDB indicates an expected call of GetDB.*
- `GetStats`
  - *GetStats indicates an expected call of GetStats.*
- `ListLogs`
  - *ListLogs indicates an expected call of ListLogs.*
- `Record`
  - *Record indicates an expected call of Record.*
- `RecordSecurity`
  - *RecordSecurity indicates an expected call of RecordSecurity.*
- `WithTransaction`
  - *WithTransaction indicates an expected call of WithTransaction.*

### `Repository`
**Kind**: Struct

**Methods:**
- `GetActivityStats`
  - *GetActivityStats returns log counts grouped by hour for the last 24 hours*
- `GetDB`
- `GetStats`
  - *GetStats returns log counts grouped by category*
- `GetTotalCount`
  - *GetTotalCount returns the total count of system log entries matching filters*
- `List`
- `Record`
- `WithTransaction`

**Constructors/Factory Functions:**
- `NewRepository`

### `RepositoryInterface`
**Kind**: Interface

### `Service`
**Kind**: Struct

**Methods:**
- `GetActivityStats`
- `GetDB`
- `GetStats`
- `ListLogs`
- `Record`
- `RecordSecurity`
- `WithTransaction`

**Constructors/Factory Functions:**
- `NewService`

### `ServiceInterface`
ServiceInterface defines the business logic for system logging.

**Kind**: Interface

### `SystemLog`
SystemLog represents a pure business entity for a system log entry.

**Kind**: Struct

**Constructors/Factory Functions:**
- `MapSystemLogToDomain`
- `MapSystemLogsToDomain`

### `SystemLogDB`
SystemLogDB represents the database model for the system_logs table.

**Kind**: Struct

**Constructors/Factory Functions:**
- `MapSystemLogToDB`

### `SystemLogDTO`
SystemLogDTO is the response DTO for a single system log entry

**Kind**: Struct

## Package Level Functions

### `RegisterRoutes`
### `TestService_ListLogs`
### `TestService_Record`
### `TestService_RecordSecurity`
