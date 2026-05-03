# Contributor's Guide: OGOS Architecture

Welcome. If you are reading this, you are likely a third-year student taking over this project. Don't panic. This project is built using **Purist Clean Architecture** and **Vertical Slicing**. It is designed to be maintainable, testable, and resilient to change.

## 1. The Core Philosophy
We do not write "spaghetti code." We separate **Business Logic** from **Infrastructure**.

- **Domain (The Soul)**: What a student *is* in the real world.
- **Persistence (The Shell)**: How a student is stored in a database.
- **Service (The Brain)**: The rules and logic (e.g., "A student needs a COR to enroll").
- **Repository (The Hand)**: The part that actually talks to the Database.

## 2. Project Structure (Vertical Slices)
Instead of putting all "Handlers" in one folder, we group everything by **Feature**. Look inside `internal/features/`. Each folder (e.g., `students`) is a self-contained unit.

### Inside a Feature Folder:
| File | Purpose |
| :--- | :--- |
| `domain.go` | **Pure** structs. No database tags. No dependencies. |
| `persistence.go` | Database-specific structs. Contains `db` and `json` tags. |
| `mapper.go` | The "Bridge." Functions to convert between Domain and Persistence. |
| `interface.go` | The "Contract." Defines what the Service and Repo can do. |
| `service.go` | Business logic. Only talks to the Repository and Domain. |
| `repository.go` | SQL queries. Fetches Persistence models, maps them to Domain. |
| `handler.go` | HTTP logic. Receives requests, calls Service, sends JSON. |
| `routes.go` | Maps URLs to specific Handlers. |

## 3. Strict Coding Standards
If you break these rules, the senior devs (and the linter) will be unhappy.

### 80-Character Line Limit
No line of code should exceed 80 characters. If a function call or struct tag is too long, break it into multiple lines. This ensures the code is readable on any screen without horizontal scrolling.

### Structured Error Logging
Use the following format for all errors in Handlers:
`fmt.Printf("[HandlerName] {Specific Step}: error message\n")`

### Standardized Naming
Handlers MUST follow: `[HTTP Method][Resource]`
- *Correct:* `PostStudentProfile`, `GetStudentIIR`
- *Incorrect:* `HandleSave`, `StudentGet`

### Pure Domain
Never, ever put a `sql.NullString` or a `db` tag inside `domain.go`. If you need a nullable string in the domain, use `structs.NullableString`.

## 4. How to Add a Feature
1.  **Define the Domain**: Write the core structs in `domain.go`.
2.  **Define the Interface**: Decide what methods you need in `interface.go`.
3.  **Build the Persistence**: Create the DB models in `persistence.go`.
4.  **Write the Mapper**: Implement `ToDomain` and `ToPersistence` in `mapper.go`.
5.  **Implement the Repo**: Write the SQL in `repository.go`.
6.  **Implement the Service**: Write the logic in `service.go`.
7.  **Hook it up**: Create the `handler.go` and register it in `routes.go`.

**Reference Implementation:** If you are confused, look at `internal/features/students/`. It is the gold standard for this architecture.
