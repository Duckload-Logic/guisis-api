# Core Module: `tokens`

## Types and Interfaces

### `Claims`
**Kind**: Struct

### `Service`
**Kind**: Struct

**Methods:**
- `GenerateToken`
- `ParseTokenUnverified`
  - *ParseTokenUnverified extracts claims from a token string without verifying its signature or expiration. Use this ONLY to identify a session for refresh logic, never for authorization.*
- `ValidateToken`

**Constructors/Factory Functions:**
- `NewService`

## Package Level Functions

### `TestNewService`
### `TestService_TokenWorkflow`
