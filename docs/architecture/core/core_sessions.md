# Core Module: `sessions`

## Types and Interfaces

### `JTIDTO`
JTIDTO encapsulates the JWT ID (jti) and provides helper methods
for generating Redis keys to ensure consistency and maintainability.

**Kind**: Struct

**Methods:**
- `ToIDPRefreshKey`
  - *ToIDPRefreshKey returns the Redis key for the linked IDP refresh token.*
- `ToSessionKey`
  - *ToSessionKey returns the Redis key for the primary session data.*

**Constructors/Factory Functions:**
- `NewJTI`

### `Service`
**Kind**: Struct

**Methods:**
- `DeleteToken`
  - *DeleteToken removes session data from Redis.*
- `DeleteUserToken`
  - *DeleteUserToken removes a session and its link to the user.*
- `GetToken`
  - *GetToken retrieves session data from Redis.*
- `ListUserSessions`
  - *ListUserSessions returns all active session data for a user.*
- `StoreToken`
  - *StoreToken saves session data in Redis with a JTI-based key.*
- `StoreUserToken`
  - *StoreUserToken saves session data and links it to a user for auditing.*

**Constructors/Factory Functions:**
- `NewService`

## Package Level Functions

### `TestService_DeleteToken`
### `TestService_GetToken`
### `TestService_StoreToken`
### `ToUserSessionsKey`
ToUserSessionsKey returns the Redis key for the set of sessions belonging
to a specific user.

