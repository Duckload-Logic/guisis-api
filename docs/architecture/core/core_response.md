# Core Module: `response`

## Types and Interfaces

### `CommonErrorResponse`
CommonErrorResponse is a flat, unified struct for all 4xx and 5xx responses.
This is used for Swagger documentation to provide a clean, non-composite schema.

**Kind**: Struct

### `ErrorResponse`
ErrorResponse is a Swagger-friendly wrapper for a server-side JSend error.

**Kind**: Struct

### `FailResponse`
FailResponse is a Swagger-friendly wrapper for a client-side JSend failure.

**Kind**: Struct

### `JSendResponse`
JSendResponse is the generic wrapper for all API responses following the
JSend specification.

**Kind**: Struct

### `JSendStatus`
**Kind**: Type

### `SuccessResponse`
SuccessResponse is a Swagger-friendly wrapper for a successful JSend response.

**Kind**: Struct

## Package Level Functions

### `SendError`
SendError sends a JSend response for server-side errors.

### `SendFail`
SendFail sends a JSend response for client-side errors (e.g., validation).

### `SendSuccess`
SendSuccess sends a successful JSend response.

### `TestSendError`
### `TestSendFail`
### `TestSendSuccess`
