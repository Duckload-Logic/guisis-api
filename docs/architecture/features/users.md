# Feature Module: `users`

## Overview
Package users is a generated GoMock package.

## Types and Interfaces

### `CreateUserRequest`
**Kind**: Struct

### `GetUserResponse`
**Kind**: Struct

### `Handler`
**Kind**: Struct

**Methods:**
- `DeleteUserSession`
- `GetMe`
  - *GetMe godoc @Summary      Get current user @Description  Retrieves information about the currently authenticated user. @Tags         Users @Accept       json @Produce      json @Success      200      {object}  GetUserResponse        "Returns current user details" @Failure      500      {object}  map[string]string     "Failed to get current user" @Router       /users/me [get] GetMe retrieves the currently authenticated user's information.*
- `GetRoleDistribution`
  - *GetRoleDistribution godoc @Summary      Get user role distribution @Description  Returns the count of users for each role in the system. @Tags         Users @Accept       json @Produce      json @Success      200      {array}   RoleDistributionDTO @Router       /users/distribution [get]*
- `GetUserActivity`
- `GetUserByEmail`
  - *GetUserByEmail godoc @Summary      Get user by email @Description  Retrieves user information based on the provided email. @Tags         Users @Accept       json @Produce      json @Param        email   query     string true "User Email" @Success      200      {object}  GetUserResponse        "Returns user details" @Failure      400      {object}  map[string]string     "Email query parameter is required" @Failure      500      {object}  map[string]string     "Failed to get user by email" @Router       /users [get] GetUserByEmail retrieves user information by their email address.*
- `GetUserSessions`
- `GetUsers`
  - *GetUsers godoc @Summary      List all users @Description  Retrieves a paginated list of all users with filtering options. @Tags         Users @Accept       json @Produce      json @Param        page       query     int     false  "Page number" @Param        page_size  query     int     false  "Items per page" @Param        role_id    query     int     false  "Filter by role" @Param        search     query     string  false  "Search by name/email" @Param        active     query     bool    false  "Filter by status" @Success      200        {object}  ListUsersResponse @Router       /users/all [get]*
- `PostUserBlock`
- `PostUserUnblock`

**Constructors/Factory Functions:**
- `NewHandler`

### `ListUsersParams`
**Kind**: Struct

### `ListUsersResponse`
**Kind**: Struct

### `MockRepositoryInterface`
MockRepositoryInterface is a mock of RepositoryInterface interface.

**Kind**: Struct

**Methods:**
- `BlockUser`
  - *BlockUser mocks base method.*
- `CheckUserWhitelist`
  - *CheckUserWhitelist mocks base method.*
- `CreateUser`
  - *CreateUser mocks base method.*
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetDB`
  - *GetDB mocks base method.*
- `GetRoleByID`
  - *GetRoleByID mocks base method.*
- `GetRoleDistribution`
  - *GetRoleDistribution mocks base method.*
- `GetUserByEmail`
  - *GetUserByEmail mocks base method.*
- `GetUserByID`
  - *GetUserByID mocks base method.*
- `GetUserIDsByRole`
  - *GetUserIDsByRole mocks base method.*
- `ListUsers`
  - *ListUsers mocks base method.*
- `PostProfilePicture`
  - *PostProfilePicture mocks base method.*
- `UnblockUser`
  - *UnblockUser mocks base method.*
- `WithTransaction`
  - *WithTransaction mocks base method.*

**Constructors/Factory Functions:**
- `NewMockRepositoryInterface`

### `MockRepositoryInterfaceMockRecorder`
MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.

**Kind**: Struct

**Methods:**
- `BlockUser`
  - *BlockUser indicates an expected call of BlockUser.*
- `CheckUserWhitelist`
  - *CheckUserWhitelist indicates an expected call of CheckUserWhitelist.*
- `CreateUser`
  - *CreateUser indicates an expected call of CreateUser.*
- `GetDB`
  - *GetDB indicates an expected call of GetDB.*
- `GetRoleByID`
  - *GetRoleByID indicates an expected call of GetRoleByID.*
- `GetRoleDistribution`
  - *GetRoleDistribution indicates an expected call of GetRoleDistribution.*
- `GetUserByEmail`
  - *GetUserByEmail indicates an expected call of GetUserByEmail.*
- `GetUserByID`
  - *GetUserByID indicates an expected call of GetUserByID.*
- `GetUserIDsByRole`
  - *GetUserIDsByRole indicates an expected call of GetUserIDsByRole.*
- `ListUsers`
  - *ListUsers indicates an expected call of ListUsers.*
- `PostProfilePicture`
  - *PostProfilePicture indicates an expected call of PostProfilePicture.*
- `UnblockUser`
  - *UnblockUser indicates an expected call of UnblockUser.*
- `WithTransaction`
  - *WithTransaction indicates an expected call of WithTransaction.*

### `MockServiceInterface`
MockServiceInterface is a mock of ServiceInterface interface.

**Kind**: Struct

**Methods:**
- `BlockUser`
  - *BlockUser mocks base method.*
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetRoleDistribution`
  - *GetRoleDistribution mocks base method.*
- `GetUserByEmail`
  - *GetUserByEmail mocks base method.*
- `GetUserByID`
  - *GetUserByID mocks base method.*
- `GetUserIDsByRole`
  - *GetUserIDsByRole mocks base method.*
- `ListUsers`
  - *ListUsers mocks base method.*
- `PostProfilePicture`
  - *PostProfilePicture mocks base method.*
- `UnblockUser`
  - *UnblockUser mocks base method.*

**Constructors/Factory Functions:**
- `NewMockServiceInterface`

### `MockServiceInterfaceMockRecorder`
MockServiceInterfaceMockRecorder is the mock recorder for MockServiceInterface.

**Kind**: Struct

**Methods:**
- `BlockUser`
  - *BlockUser indicates an expected call of BlockUser.*
- `GetRoleDistribution`
  - *GetRoleDistribution indicates an expected call of GetRoleDistribution.*
- `GetUserByEmail`
  - *GetUserByEmail indicates an expected call of GetUserByEmail.*
- `GetUserByID`
  - *GetUserByID indicates an expected call of GetUserByID.*
- `GetUserIDsByRole`
  - *GetUserIDsByRole indicates an expected call of GetUserIDsByRole.*
- `ListUsers`
  - *ListUsers indicates an expected call of ListUsers.*
- `PostProfilePicture`
  - *PostProfilePicture indicates an expected call of PostProfilePicture.*
- `UnblockUser`
  - *UnblockUser indicates an expected call of UnblockUser.*

### `ProfilePicture`
**Kind**: Struct

**Methods:**
- `ToPersistence`

### `ProfilePictureDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *ProfilePicture mappers*

### `Repository`
**Kind**: Struct

**Methods:**
- `BlockUser`
- `CheckUserWhitelist`
- `CreateUser`
  - *CreateUser inserts or updates a user using the persistence model.*
- `GetDB`
- `GetRoleByID`
- `GetRoleDistribution`
- `GetUserByEmail`
- `GetUserByID`
  - *GetUserByID fetches a user by their ID and maps to Domain.*
- `GetUserIDsByRole`
- `ListUsers`
- `PostProfilePicture`
- `UnblockUser`
- `WithTransaction`

### `RepositoryInterface`
**Kind**: Interface

**Constructors/Factory Functions:**
- `NewRepository`

### `Role`
**Kind**: Struct

**Methods:**
- `ToPersistence`

### `RoleDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Role mappers*

### `RoleDistributionDTO`
**Kind**: Struct

### `Service`
**Kind**: Struct

**Methods:**
- `BlockUser`
- `GetRoleDistribution`
- `GetUserByEmail`
  - *GetUserByEmail retrieves a user by their email and auth type.*
- `GetUserByID`
  - *GetUserByID retrieves a user by their ID.*
- `GetUserIDsByRole`
- `ListUsers`
- `PostProfilePicture`
- `UnblockUser`

### `ServiceInterface`
**Kind**: Interface

**Constructors/Factory Functions:**
- `NewService`

### `User`
**Kind**: Struct

**Methods:**
- `ToPersistence`

### `UserDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *User mappers*

## Package Level Functions

### `RegisterRoutes`
### `TestService_BlockUser`
### `TestService_GetUserByID`
