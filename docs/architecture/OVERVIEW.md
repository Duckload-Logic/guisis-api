# Purist Clean Architecture Overview

This directory contains automatically generated technical breakdowns of every core and feature module within the `ogos-gin-api` project. These files are meant to be used as authoritative references for understanding module boundaries.

## The Guiding Principle

The defining characteristic of this project's architecture is the strict, undeniable separation between the **Domain Layer** (business rules) and the **Persistence Layer** (infrastructure & database). 

**Domain Models do not know about the database.**
**Persistence Models do not know about the business logic.**

## The Execution Flow

Every feature module follows a rigid execution pipeline:

1.  **Handler (`handler.go`)**: 
    - Receives HTTP requests.
    - Parses JSON bodies into Data Transfer Objects (DTOs).
    - Validates input.
    - Calls the Service.

2.  **Service (`service.go`)**: 
    - The brain of the feature.
    - Operates **EXTCLUSIVELY** on Domain Models.
    - Calls the Repository through an Interface.

3.  **Interface (`interface.go`)**:
    - Defines the contract between the Service and the Repository.
    - All method signatures must use pure Domain structures. You will never see `sql.NullString` in an interface signature.

4.  **Mapper (`mapper.go`)**: 
    - The bridge.
    - Converts Domain objects to Persistence objects before saving to the database.
    - Converts Persistence objects back into Domain objects before returning them to the Service.

5.  **Repository (`repository.go`)**: 
    - Executes SQL queries against the database.
    - Implements the Repository interface. Uses `mapper.go` internally to satisfy the Domain-only interface contract.

## Studying for your Presentation

To effectively command an understanding of this system on camera:
1. Do not memorize the code. 
2. Open the generated `.md` files in this directory.
3. Observe how the `Interface` types interact strictly with the `Struct` domain types.
4. If asked "How does a student get saved?", answer by walking down the Execution Flow listed above.

Do not wing it. Trust the boundaries.
