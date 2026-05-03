# Core Module: `middleware`

## Overview

owner.go

## Types and Interfaces

### `IPRateLimiter`

**Kind**: Struct

**Methods:**

- `GetLimiter`

**Constructors/Factory Functions:**

- `NewIPRateLimiter`

### `SecurityLogger`

SecurityLogger is an interface for recording security log events.
This interface is implemented by logs.Service and is used to break
the import cycle between middleware and logs packages.

**Kind**: Interface

## Package Level Functions

### `AuditContextMiddleware`

AuditContextMiddleware enriches the request context with audit metadata
(user email, IP address, User-Agent). Must be placed after AuthMiddleware
so that userEmail is available.

### `AuthMiddleware`

### `HydrateStudentIIRContext`

HydrateStudentIIRContext extracts student IIR ID from database and
sets it in the Gin context. Only applies to Student role users.
Day One students (no IIR record) can still proceed.

### `OwnershipMiddleware`

OwnershipMiddleware - Direct database access version

### `RateLimitMiddleware`

### `RoleMiddleware`

RoleMiddleware checks if the user's role is in the allowed list.
Optionally accepts a log service to record ACCESS_DENIED events.

### `TraceMiddleware`
