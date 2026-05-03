# Feature Module: `files`

## Types and Interfaces

### `File`
File represents a pure business entity for a file asset.

**Kind**: Struct

**Constructors/Factory Functions:**
- `MapFileToDomain`

### `FileDB`
FileDB represents the database model for a file.

**Kind**: Struct

**Constructors/Factory Functions:**
- `MapFileToDB`

### `OCRResult`
OCRResult represents the pure result of an OCR process.

**Kind**: Struct

**Constructors/Factory Functions:**
- `MapOCRResultToDomain`

### `OCRResultDB`
OCRResultDB represents the database model for OCR results.

**Kind**: Struct

**Constructors/Factory Functions:**
- `MapOCRResultToDB`

### `Repository`
**Kind**: Struct

**Methods:**
- `Create`
- `CreateBulk`
- `Delete`
- `GetDB`
- `GetFileByID`
- `GetFilesByIDs`
- `GetOCRResultByFileID`
- `SaveOCRResult`
- `WithTransaction`

### `RepositoryInterface`
**Kind**: Interface

**Constructors/Factory Functions:**
- `NewRepository`

### `Service`
**Kind**: Struct

**Methods:**
- `DeleteFile`
- `GetFileByID`
- `GetOCRResult`
- `UploadFile`
- `UploadFiles`

### `ServiceInterface`
**Kind**: Interface

**Constructors/Factory Functions:**
- `NewService`

