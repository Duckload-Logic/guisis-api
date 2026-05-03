# Core Module: `structs`

## Types and Interfaces

### `JSendResponse`
JSendResponse is the generic wrapper for all API responses following the
JSend specification.

**Kind**: Struct

**Constructors/Factory Functions:**
- `NewErrorResponse`
- `NewFailResponse`
- `NewSuccessResponse`

### `JSendStatus`
**Kind**: Type

### `NullableInt64`
**Kind**: Type

**Methods:**
- `MarshalJSON`
- `UnmarshalJSON`

**Constructors/Factory Functions:**
- `FromSqlNullInt64`

### `NullableString`
**Kind**: Type

**Methods:**
- `MarshalJSON`
- `UnmarshalJSON`

**Constructors/Factory Functions:**
- `FromSqlNull`
- `PointerToNullableString`
- `StringToNullableString`

### `NullableTime`
**Kind**: Type

**Methods:**
- `MarshalJSON`
- `UnmarshalJSON`

**Constructors/Factory Functions:**
- `FromSqlNullTime`
- `TimeToNullableTime`

### `PaginationMetadata`
PaginationMetadata contains calculated pagination information for responses.

**Kind**: Struct

**Constructors/Factory Functions:**
- `CalculateMetadata`

### `PaginationRequest`
PaginationRequest represents a standard paginated request.

**Kind**: Struct

**Methods:**
- `GetOffset`
  - *GetOffset calculates the SQL offset.*
- `SetDefaults`
  - *SetDefaults sets default pagination values if not provided or invalid.*

## Package Level Functions

### `TestCalculateMetadata`
### `TestNewErrorResponse`
### `TestNewFailResponse`
### `TestNewSuccessResponse`
### `TestNullableInt64_JSON`
### `TestNullableString_MarshalJSON`
### `TestNullableString_UnmarshalJSON`
### `TestPaginationRequest_GetOffset`
### `TestPaginationRequest_SetDefaults`
### `TestSqlNullConversion`
### `TestStringToNullableString`
### `ToSqlNull`
### `ToSqlNullInt64`
### `ToSqlNullTime`
