# Feature Module: `auth`

## Overview
Package auth is a generated GoMock package.

## Types and Interfaces

### `Handler`
**Kind**: Struct

**Methods:**
- `GetAuthorizeURL`
  - *GetAuthorizeURL godoc @Summary      Get IDP authorization URL @Description  Redirects to OAuth 2.0 authorization page on the IDP. @Tags         Auth @Produce      json @Success      302 {string} string "Redirect to IDP login page" @Failure      500 {object} map[string]string @Router       /auth/idp/authorize [get]*
- `GetLogout`
  - *GetLogout godoc @Summary      User logout @Description  Invalidates the user's tokens by clearing cookies. @Tags         Auth @Success      200 {object} map[string]string @Router       /auth/logout [get]*
- `GetMe`
  - *GetMe godoc @Summary      Get current user info @Description  Retrieves information about the currently authenticated user (native or IDP). @Tags         Auth @Produce      json @Success      200 {object} MeResponse @Failure      401 {object} map[string]string @Router       /auth/me [get]*
- `PostIDPToken`
  - *PostIDPToken godoc @Summary      Exchange IDP authorization code for tokens @Description  Completes OAuth 2.0 flow and provisions user @Tags         Auth @Accept       json @Produce      json @Param        request body idp.IDPTokenExchangeRequest true "Code & State" @Success      200 {object} map[string]interface{} @Failure      400 {object} map[string]string @Failure      401 {object} map[string]string @Router       /auth/idp/token [post]*
- `PostLogin`
  - *PostLogin godoc @Summary      User login @Description  Authenticates a user and sets JWT cookies. @Tags         Auth @Accept       json @Produce      json @Param        request body      LoginDTO true "Login Credentials" @Success      200     {object}  map[string]interface{} "Returns user info (optional)" @Failure      400     {object}  map[string]string @Failure      401     {object}  map[string]string @Router       /auth/login [post]*
- `PostRefreshToken`
  - *PostRefreshToken godoc @Summary      Refresh JWT token @Description  Refreshes the JWT token using the refresh token cookie. @Tags         Auth @Accept       json @Produce      json @Success      200 {object} map[string]string "New access token (optional)" @Failure      401 {object} map[string]string @Router       /auth/refresh [post]*
- `PostRegister`
  - *PostRegister godoc @Summary      User registration @Description  Creates a new developer account. @Tags         Auth @Accept       json @Produce      json @Param        request body      RegisterDTO true "Registration Data" @Success      201     {object}  map[string]string "Success message" @Failure      400     {object}  map[string]string @Failure      409     {object}  map[string]string @Router       /auth/register [post]*
- `PostResendVerification`
- `PostVerify`

**Constructors/Factory Functions:**
- `NewHandler`

### `IDPRefreshRequest`
**Kind**: Struct

### `LoginDTO`
**Kind**: Struct

### `MeResponse`
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
- `GetRoleByID`
  - *GetRoleByID mocks base method.*
- `GetUserByEmail`
  - *GetUserByEmail mocks base method.*
- `GetUserByID`
  - *GetUserByID mocks base method.*
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
- `GetRoleByID`
  - *GetRoleByID indicates an expected call of GetRoleByID.*
- `GetUserByEmail`
  - *GetUserByEmail indicates an expected call of GetUserByEmail.*
- `GetUserByID`
  - *GetUserByID indicates an expected call of GetUserByID.*
- `UnblockUser`
  - *UnblockUser indicates an expected call of UnblockUser.*
- `WithTransaction`
  - *WithTransaction indicates an expected call of WithTransaction.*

### `MockServiceInterface`
MockServiceInterface is a mock of ServiceInterface interface.

**Kind**: Struct

**Methods:**
- `AuthenticateUser`
  - *AuthenticateUser mocks base method.*
- `BlockUser`
  - *BlockUser mocks base method.*
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetAuthorizeURL`
  - *GetAuthorizeURL mocks base method.*
- `GetIDPUserInfo`
  - *GetIDPUserInfo mocks base method.*
- `GetMe`
  - *GetMe mocks base method.*
- `Logout`
  - *Logout mocks base method.*
- `PostIDPTokenExchange`
  - *PostIDPTokenExchange mocks base method.*
- `RefreshIDPToken`
  - *RefreshIDPToken mocks base method.*
- `RefreshToken`
  - *RefreshToken mocks base method.*
- `RegisterUser`
  - *RegisterUser mocks base method.*
- `ResendVerification`
  - *ResendVerification mocks base method.*
- `UnblockUser`
  - *UnblockUser mocks base method.*
- `ValidateIDPSession`
  - *ValidateIDPSession mocks base method.*
- `VerifyUser`
  - *VerifyUser mocks base method.*
- `WithTransaction`
  - *WithTransaction mocks base method.*

**Constructors/Factory Functions:**
- `NewMockServiceInterface`

### `MockServiceInterfaceMockRecorder`
MockServiceInterfaceMockRecorder is the mock recorder for MockServiceInterface.

**Kind**: Struct

**Methods:**
- `AuthenticateUser`
  - *AuthenticateUser indicates an expected call of AuthenticateUser.*
- `BlockUser`
  - *BlockUser indicates an expected call of BlockUser.*
- `GetAuthorizeURL`
  - *GetAuthorizeURL indicates an expected call of GetAuthorizeURL.*
- `GetIDPUserInfo`
  - *GetIDPUserInfo indicates an expected call of GetIDPUserInfo.*
- `GetMe`
  - *GetMe indicates an expected call of GetMe.*
- `Logout`
  - *Logout indicates an expected call of Logout.*
- `PostIDPTokenExchange`
  - *PostIDPTokenExchange indicates an expected call of PostIDPTokenExchange.*
- `RefreshIDPToken`
  - *RefreshIDPToken indicates an expected call of RefreshIDPToken.*
- `RefreshToken`
  - *RefreshToken indicates an expected call of RefreshToken.*
- `RegisterUser`
  - *RegisterUser indicates an expected call of RegisterUser.*
- `ResendVerification`
  - *ResendVerification indicates an expected call of ResendVerification.*
- `UnblockUser`
  - *UnblockUser indicates an expected call of UnblockUser.*
- `ValidateIDPSession`
  - *ValidateIDPSession indicates an expected call of ValidateIDPSession.*
- `VerifyUser`
  - *VerifyUser indicates an expected call of VerifyUser.*
- `WithTransaction`
  - *WithTransaction indicates an expected call of WithTransaction.*

### `RegisterDTO`
**Kind**: Struct

### `RepositoryInterface`
**Kind**: Interface

### `Service`
**Kind**: Struct

**Methods:**
- `AuthenticateUser`
  - *AuthenticateUser handles native email/password authentication.*
- `BlockUser`
- `GetAuthorizeURL`
  - *GetAuthorizeURL generates the complete OAuth 2.0 authorization URL with PKCE parameters. This method creates a state parameter for CSRF protection, generates PKCE verifier and challenge, stores the state with metadata, and builds the authorization URL.  Parameters:   - cfg: Application configuration containing IDP endpoints  Returns the authorization URL and state parameter, or an error if generation fails.*
- `GetIDPUserInfo`
  - *GetIDPUserInfo fetches user information from the IDP userinfo endpoint using the provided access token. This is typically called after a successful token exchange to retrieve user details for provisioning.  Parameters:   - ctx: Context for the HTTP request   - accessToken: Access token obtained from IDP token exchange   - cfg: Application configuration containing IDP endpoints  Returns the IDP user information or an error if retrieval fails.*
- `GetMe`
  - *GetMe retrieves the currently authenticated user's profile information.*
- `Logout`
  - *Logout invalidates the user's session in Redis and optionally the IDP.*
- `PostIDPTokenExchange`
  - *PostIDPTokenExchange orchestrates the complete IDP login flow: validates state, exchanges code for token, fetches user info, provisions user, and generates application JWT tokens.  Parameters:   - ctx: Context for database and HTTP operations   - code: Authorization code from IDP callback   - state: State parameter from IDP callback   - cfg: Application configuration  Returns user ID and JWT tokens, or an error if any step fails.*
- `RefreshIDPToken`
  - *RefreshIDPToken handles token refresh for IDP-authenticated sessions.*
- `RefreshToken`
  - *RefreshToken generates a new access token using a valid session handle.*
- `RegisterUser`
  - *RegisterUser handles native user registration.*
- `ResendVerification`
- `UnblockUser`
- `ValidateIDPSession`
  - *ValidateIDPSession checks if the provided session ID is valid on the IDP.*
- `VerifyUser`
- `WithTransaction`

**Constructors/Factory Functions:**
- `NewService`

### `ServiceInterface`
**Kind**: Interface

### `TTL`
**Kind**: Type

### `VerifyDTO`
**Kind**: Struct

## Package Level Functions

### `RegisterDebugRoutes`
### `RegisterRoutes`
### `TestService_BlockUser`
### `TestService_GetMe`
### `TestService_VerifyUser`
