# Feature Module: `notes`

## Overview
Package notes is a generated GoMock package.

## Types and Interfaces

### `Handler`
**Kind**: Struct

**Methods:**
- `GetSignificantNotes`
- `PostSignificantNote`

**Constructors/Factory Functions:**
- `NewHandler`

### `MockRepositoryInterface`
MockRepositoryInterface is a mock of RepositoryInterface interface.

**Kind**: Struct

**Methods:**
- `CreateSignificantNote`
  - *CreateSignificantNote mocks base method.*
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetStudentSignificantNotes`
  - *GetStudentSignificantNotes mocks base method.*
- `HasNoteForAppointment`
  - *HasNoteForAppointment mocks base method.*
- `WithTransaction`
  - *WithTransaction mocks base method.*

**Constructors/Factory Functions:**
- `NewMockRepositoryInterface`

### `MockRepositoryInterfaceMockRecorder`
MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.

**Kind**: Struct

**Methods:**
- `CreateSignificantNote`
  - *CreateSignificantNote indicates an expected call of CreateSignificantNote.*
- `GetStudentSignificantNotes`
  - *GetStudentSignificantNotes indicates an expected call of GetStudentSignificantNotes.*
- `HasNoteForAppointment`
  - *HasNoteForAppointment indicates an expected call of HasNoteForAppointment.*
- `WithTransaction`
  - *WithTransaction indicates an expected call of WithTransaction.*

### `MockServiceInterface`
MockServiceInterface is a mock of ServiceInterface interface.

**Kind**: Struct

**Methods:**
- `CreateSignificantNote`
  - *CreateSignificantNote mocks base method.*
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetStudentSignificantNotes`
  - *GetStudentSignificantNotes mocks base method.*
- `HasNoteForAppointment`
  - *HasNoteForAppointment mocks base method.*
- `WithTransaction`
  - *WithTransaction mocks base method.*

**Constructors/Factory Functions:**
- `NewMockServiceInterface`

### `MockServiceInterfaceMockRecorder`
MockServiceInterfaceMockRecorder is the mock recorder for MockServiceInterface.

**Kind**: Struct

**Methods:**
- `CreateSignificantNote`
  - *CreateSignificantNote indicates an expected call of CreateSignificantNote.*
- `GetStudentSignificantNotes`
  - *GetStudentSignificantNotes indicates an expected call of GetStudentSignificantNotes.*
- `HasNoteForAppointment`
  - *HasNoteForAppointment indicates an expected call of HasNoteForAppointment.*
- `WithTransaction`
  - *WithTransaction indicates an expected call of WithTransaction.*

### `Repository`
**Kind**: Struct

**Methods:**
- `CreateSignificantNote`
- `GetStudentSignificantNotes`
- `HasNoteForAppointment`
- `WithTransaction`

### `RepositoryInterface`
RepositoryInterface defines the data access layer for managing student notes.

**Kind**: Interface

**Constructors/Factory Functions:**
- `NewRepository`

### `Service`
**Kind**: Struct

**Methods:**
- `CreateSignificantNote`
- `GetStudentSignificantNotes`
- `HasNoteForAppointment`
- `WithTransaction`

**Constructors/Factory Functions:**
- `NewService`

### `ServiceInterface`
ServiceInterface defines the business logic for managing student notes.

**Kind**: Interface

### `SignificantNote`
SignificantNote represents a pure business entity for professional notes.

**Kind**: Struct

**Constructors/Factory Functions:**
- `MapSignificantNoteToDomain`
- `MapSignificantNotesToDomain`

### `SignificantNoteDB`
SignificantNoteDB represents the database model for professional notes.

**Kind**: Struct

**Constructors/Factory Functions:**
- `MapSignificantNoteToDB`

### `SignificantNoteDTO`
**Kind**: Struct

## Package Level Functions

### `RegisterRoutes`
### `TestService_GetStudentSignificantNotes`
### `TestService_HasNoteForAppointment`
