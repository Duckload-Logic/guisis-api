# Feature Module: `notifications`

## Overview
Package notifications is a generated GoMock package.

## Types and Interfaces

### `Handler`
**Kind**: Struct

**Methods:**
- `GetNotifications`
- `PatchNotificationRead`

**Constructors/Factory Functions:**
- `NewHandler`

### `ListNotificationsResponse`
**Kind**: Struct

### `MockRepositoryInterface`
MockRepositoryInterface is a mock of RepositoryInterface interface.

**Kind**: Struct

**Methods:**
- `Create`
  - *Create mocks base method.*
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetByUserID`
  - *GetByUserID mocks base method.*
- `MarkAsRead`
  - *MarkAsRead mocks base method.*
- `WithTransaction`
  - *WithTransaction mocks base method.*

**Constructors/Factory Functions:**
- `NewMockRepositoryInterface`

### `MockRepositoryInterfaceMockRecorder`
MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.

**Kind**: Struct

**Methods:**
- `Create`
  - *Create indicates an expected call of Create.*
- `GetByUserID`
  - *GetByUserID indicates an expected call of GetByUserID.*
- `MarkAsRead`
  - *MarkAsRead indicates an expected call of MarkAsRead.*
- `WithTransaction`
  - *WithTransaction indicates an expected call of WithTransaction.*

### `MockServiceInterface`
MockServiceInterface is a mock of ServiceInterface interface.

**Kind**: Struct

**Methods:**
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetUserNotifications`
  - *GetUserNotifications mocks base method.*
- `MarkAsRead`
  - *MarkAsRead mocks base method.*
- `Send`
  - *Send mocks base method.*

**Constructors/Factory Functions:**
- `NewMockServiceInterface`

### `MockServiceInterfaceMockRecorder`
MockServiceInterfaceMockRecorder is the mock recorder for MockServiceInterface.

**Kind**: Struct

**Methods:**
- `GetUserNotifications`
  - *GetUserNotifications indicates an expected call of GetUserNotifications.*
- `MarkAsRead`
  - *MarkAsRead indicates an expected call of MarkAsRead.*
- `Send`
  - *Send indicates an expected call of Send.*

### `Notification`
Notification represents a pure business entity for system notifications.

**Kind**: Struct

**Constructors/Factory Functions:**
- `MapNotificationToDomain`
- `MapNotificationsToDomain`

### `NotificationDB`
NotificationDB represents the database model for the notifications table.

**Kind**: Struct

**Constructors/Factory Functions:**
- `MapNotificationToDB`

### `NotificationResponse`
**Kind**: Struct

### `Repository`
**Kind**: Struct

**Methods:**
- `Create`
- `GetByUserID`
- `MarkAsRead`
- `WithTransaction`

**Constructors/Factory Functions:**
- `NewRepository`

### `RepositoryInterface`
RepositoryInterface defines the data access layer for managing notifications.

**Kind**: Interface

### `Service`
**Kind**: Struct

**Methods:**
- `GetUserNotifications`
- `MarkAsRead`
- `Send`

**Constructors/Factory Functions:**
- `NewService`

### `ServiceInterface`
ServiceInterface defines the business logic for managing notifications.

**Kind**: Interface

## Package Level Functions

### `RegisterRoutes`
### `TestService_GetUserNotifications`
### `TestService_Send`
