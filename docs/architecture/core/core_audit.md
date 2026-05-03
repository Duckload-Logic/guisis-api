# Core Module: `audit`

## Overview
Package audit is a generated GoMock package.

## Types and Interfaces

### `DispatchParams`
DispatchParams holds the parameters for the Dispatch helper.

**Kind**: Struct

### `LogEntry`
LogEntry is the input struct used by other services to record a log.

**Kind**: Struct

### `LogMetadata`
LogMetadata defines a structured format for audit log metadata.

**Kind**: Struct

### `LogParams`
LogParams holds the parameters for a log entry.

**Kind**: Struct

### `Logger`
Logger defines the interface for recording system logs.

**Kind**: Interface

### `MockLogger`
MockLogger is a mock of Logger interface.

**Kind**: Struct

**Methods:**
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `Record`
  - *Record mocks base method.*

**Constructors/Factory Functions:**
- `NewMockLogger`

### `MockLoggerMockRecorder`
MockLoggerMockRecorder is the mock recorder for MockLogger.

**Kind**: Struct

**Methods:**
- `Record`
  - *Record indicates an expected call of Record.*

### `MockNotifier`
MockNotifier is a mock of Notifier interface.

**Kind**: Struct

**Methods:**
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `Send`
  - *Send mocks base method.*

**Constructors/Factory Functions:**
- `NewMockNotifier`

### `MockNotifierMockRecorder`
MockNotifierMockRecorder is the mock recorder for MockNotifier.

**Kind**: Struct

**Methods:**
- `Send`
  - *Send indicates an expected call of Send.*

### `MockUserGetter`
MockUserGetter is a mock of UserGetter interface.

**Kind**: Struct

**Methods:**
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetUserIDsByRole`
  - *GetUserIDsByRole mocks base method.*

**Constructors/Factory Functions:**
- `NewMockUserGetter`

### `MockUserGetterMockRecorder`
MockUserGetterMockRecorder is the mock recorder for MockUserGetter.

**Kind**: Struct

**Methods:**
- `GetUserIDsByRole`
  - *GetUserIDsByRole indicates an expected call of GetUserIDsByRole.*

### `NotificationEntry`
NotificationEntry is the input struct used by other services to send a
notification.

**Kind**: Struct

### `NotificationParams`
NotificationParams holds the parameters for a notification.

**Kind**: Struct

### `Notifier`
Notifier defines the interface for sending notifications.

**Kind**: Interface

### `UserGetter`
UserGetter defines the interface for fetching user IDs by role.

**Kind**: Interface

## Package Level Functions

### `Dispatch`
Dispatch is a centralized helper to record logs and send notifications.
It automatically extracts audit metadata from the context.

### `ExtractMeta`
ExtractMeta reads audit metadata from a context.

### `ExtractUserID`
### `WithContext`
WithContext enriches a context with audit metadata.

