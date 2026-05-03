# Feature Module: `m2mclients`

## Overview
Package m2mclients is a generated GoMock package.

## Types and Interfaces

### `CreateM2MClientRequest`
**Kind**: Struct

### `CreateM2MClientResponse`
**Kind**: Struct

### `Handler`
**Kind**: Struct

**Methods:**
- `DeleteM2MClient`
- `GetM2MClients`
- `GetService`
- `PatchM2MClientVerify`
- `PostM2MClient`
- `PostM2MClientSecret`
- `PostM2MToken`
  - *PostM2MToken godoc @Summary      M2M Token Exchange @Description  Exchanges client credentials (client_id and client_secret) for an M2M access token. @Tags         Auth @Accept       json @Produce      json @Param        request body      M2MTokenRequest true "M2M Credentials" @Success      200     {object}  M2MTokenSuccessResponse @Failure      400     {object}  response.CommonErrorResponse @Failure      401     {object}  response.CommonErrorResponse @Router       /auth/m2m/token [post]*
- `PostM2MTokenRefresh`

**Constructors/Factory Functions:**
- `NewHandler`

### `ListM2MClientsRequest`
**Kind**: Struct

### `M2MClient`
M2MClient represents a pure business entity for a Machine-to-Machine client.

**Kind**: Struct

**Constructors/Factory Functions:**
- `MapM2MClientToDomain`
- `MapM2MClientsToDomain`

### `M2MClientDB`
M2MClientDB represents the database model for the m2m_clients table.

**Kind**: Struct

**Constructors/Factory Functions:**
- `MapM2MClientToDB`

### `M2MClientDTO`
**Kind**: Struct

### `M2MClientListSuccessResponse`
M2MClientListSuccessResponse is a flat JSend response for GetM2MClients.

**Kind**: Struct

### `M2MCreateClientSuccessResponse`
M2MCreateClientSuccessResponse is a flat JSend response for PostM2MClient.

**Kind**: Struct

### `M2MMessageSuccessResponse`
M2MMessageSuccessResponse is a flat JSend response for generic messages.

**Kind**: Struct

### `M2MRefreshTokenRequest`
**Kind**: Struct

### `M2MSecretSuccessResponse`
M2MSecretSuccessResponse is a flat JSend response for PostM2MSecret.

**Kind**: Struct

### `M2MTokenRequest`
**Kind**: Struct

### `M2MTokenResponse`
**Kind**: Struct

### `M2MTokenSuccessResponse`
M2MTokenSuccessResponse is a flat JSend response for PostM2MToken and PostM2MRefresh.

**Kind**: Struct

### `MockRepositoryInterface`
MockRepositoryInterface is a mock of RepositoryInterface interface.

**Kind**: Struct

**Methods:**
- `Create`
  - *Create mocks base method.*
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetActiveByUserID`
  - *GetActiveByUserID mocks base method.*
- `GetByClientID`
  - *GetByClientID mocks base method.*
- `GetDB`
  - *GetDB mocks base method.*
- `List`
  - *List mocks base method.*
- `Revoke`
  - *Revoke mocks base method.*
- `TouchLastUsed`
  - *TouchLastUsed mocks base method.*
- `UpdateSecret`
  - *UpdateSecret mocks base method.*
- `UpdateVerificationStatus`
  - *UpdateVerificationStatus mocks base method.*
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
- `GetActiveByUserID`
  - *GetActiveByUserID indicates an expected call of GetActiveByUserID.*
- `GetByClientID`
  - *GetByClientID indicates an expected call of GetByClientID.*
- `GetDB`
  - *GetDB indicates an expected call of GetDB.*
- `List`
  - *List indicates an expected call of List.*
- `Revoke`
  - *Revoke indicates an expected call of Revoke.*
- `TouchLastUsed`
  - *TouchLastUsed indicates an expected call of TouchLastUsed.*
- `UpdateSecret`
  - *UpdateSecret indicates an expected call of UpdateSecret.*
- `UpdateVerificationStatus`
  - *UpdateVerificationStatus indicates an expected call of UpdateVerificationStatus.*
- `WithTransaction`
  - *WithTransaction indicates an expected call of WithTransaction.*

### `MockServiceInterface`
MockServiceInterface is a mock of ServiceInterface interface.

**Kind**: Struct

**Methods:**
- `Authenticate`
  - *Authenticate mocks base method.*
- `CreateClient`
  - *CreateClient mocks base method.*
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `IssueToken`
  - *IssueToken mocks base method.*
- `ListClients`
  - *ListClients mocks base method.*
- `RefreshToken`
  - *RefreshToken mocks base method.*
- `RegenerateSecret`
  - *RegenerateSecret mocks base method.*
- `RevokeClient`
  - *RevokeClient mocks base method.*
- `VerifyClient`
  - *VerifyClient mocks base method.*
- `WithTransaction`
  - *WithTransaction mocks base method.*

**Constructors/Factory Functions:**
- `NewMockServiceInterface`

### `MockServiceInterfaceMockRecorder`
MockServiceInterfaceMockRecorder is the mock recorder for MockServiceInterface.

**Kind**: Struct

**Methods:**
- `Authenticate`
  - *Authenticate indicates an expected call of Authenticate.*
- `CreateClient`
  - *CreateClient indicates an expected call of CreateClient.*
- `IssueToken`
  - *IssueToken indicates an expected call of IssueToken.*
- `ListClients`
  - *ListClients indicates an expected call of ListClients.*
- `RefreshToken`
  - *RefreshToken indicates an expected call of RefreshToken.*
- `RegenerateSecret`
  - *RegenerateSecret indicates an expected call of RegenerateSecret.*
- `RevokeClient`
  - *RevokeClient indicates an expected call of RevokeClient.*
- `VerifyClient`
  - *VerifyClient indicates an expected call of VerifyClient.*
- `WithTransaction`
  - *WithTransaction indicates an expected call of WithTransaction.*

### `Repository`
**Kind**: Struct

**Methods:**
- `Create`
- `GetActiveByUserID`
- `GetByClientID`
- `GetDB`
- `List`
- `Revoke`
- `TouchLastUsed`
- `UpdateSecret`
- `UpdateVerificationStatus`
- `WithTransaction`

**Constructors/Factory Functions:**
- `NewRepository`

### `RepositoryInterface`
**Kind**: Interface

### `Service`
**Kind**: Struct

**Methods:**
- `Authenticate`
- `CreateClient`
- `IssueToken`
- `ListClients`
- `RefreshToken`
- `RegenerateSecret`
- `RevokeClient`
- `VerifyClient`
- `WithTransaction`

**Constructors/Factory Functions:**
- `NewService`

### `ServiceInterface`
**Kind**: Interface

## Package Level Functions

### `RegisterRoutes`
RegisterRoutes sets up management and auth endpoints for M2M clients.

### `TestService_ListClients`
### `TestService_VerifyClient`
