# Feature Module: `slips`

## Overview
Package slips is a generated GoMock package.

## Types and Interfaces

### `AttachmentDTO`
**Kind**: Struct

### `CreateSlipRequest`
**Kind**: Struct

### `Handler`
**Kind**: Struct

**Methods:**
- `GetSlipAttachmentContent`
  - *GetAttachmentFile godoc @Summary      Download attachment @Description  Downloads a specific attachment file. @Tags         ExcuseSlips @Param        attachmentId path int true "Attachment ID" @Success      200  {file} binary @Failure      400  {object} map[string]string @Failure      404  {object} map[string]string @Failure      500  {object} map[string]string @Router       /slips/id/{id}/attachments/{attachmentId} [get]*
- `GetSlipAttachments`
  - *GetSlipAttachmentList godoc @Summary      Get slip attachments @Description  Retrieves attachments for a specific slip. @Tags         ExcuseSlips @Produce      json @Param        id   path      int  true  "Slip ID" @Success      200  {object} []AttachmentDTO @Failure      400  {object} map[string]string @Failure      500  {object} map[string]string @Router       /slips/id/{id}/attachments [get]*
- `GetSlipByID`
  - *GetSlipByID godoc @Summary      Get excuse slip by ID @Description  Retrieves details for a specific excuse slip. @Tags         ExcuseSlips @Produce      json @Param        id   path      string  true  "Slip ID" @Success      200  {object} SlipDTO @Failure      404  {object} map[string]string @Failure      500  {object} map[string]string @Router       /slips/id/{id} [get]*
- `GetSlipCategories`
  - *GetSlipCategoryList godoc @Summary      Get slip categories @Description  Retrieves available slip categories. @Tags         ExcuseSlips @Produce      json @Success      200  {object} []SlipCategory @Failure      500  {object} map[string]string @Router       /slips/lookups/categories [get]*
- `GetSlipMe`
  - *GetSlipListByIIR godoc @Summary      Get student's excuse slips @Description  Retrieves slips for the authenticated student. @Tags         ExcuseSlips @Produce      json @Success      200  {object} map[string]interface{} @Failure      403  {object} map[string]string @Failure      500  {object} map[string]string @Router       /slips/me [get]*
- `GetSlipStats`
  - *GetSlipStatsList godoc @Summary      Get slip statistics @Description  Retrieves slip status counts. @Tags         ExcuseSlips @Produce      json @Success      200  {object} []SlipStatusCount @Failure      500  {object} map[string]string @Router       /slips/stats [get]*
- `GetSlipStatuses`
  - *GetSlipStatusList godoc @Summary      Get slip statuses @Description  Retrieves available slip statuses. @Tags         ExcuseSlips @Produce      json @Success      200  {object} []SlipStatus @Failure      500  {object} map[string]string @Router       /slips/lookups/statuses [get]*
- `GetSlipUrgent`
  - *GetUrgentSlipList godoc @Summary      Get urgent excuse slips @Description  Retrieves urgent slips for counselor review. @Tags         ExcuseSlips @Produce      json @Success      200  {object} map[string]interface{} @Failure      500  {object} map[string]string @Router       /slips/urgent [get]*
- `GetSlips`
  - *GetSlipList godoc @Summary      Get all excuse slips @Description  Retrieves list of all submitted slips. @Tags         ExcuseSlips @Produce      json @Success      200  {object} map[string]interface{} @Failure      500  {object} map[string]string @Router       /slips [get]*
- `PatchSlip`
  - *PatchSlipStatus godoc @Summary      Update slip status @Description  Approve, reject, or request revision of slip. @Tags         ExcuseSlips @Accept       json @Produce      json @Param        id   path      int                  true  "Slip ID" @Param        body body      UpdateStatusRequest  true  "Status and notes" @Success      200  {object} map[string]string @Failure      400  {object} map[string]string @Failure      404  {object} map[string]string @Failure      500  {object} map[string]string @Router       /slips/id/{id}/status [patch]*
- `PatchSlipStatus`
  - *PatchSlipStatus godoc*
- `PostSlip`
  - *PostSlip godoc @Summary      Submit an excuse slip @Description  Allows a student to submit an excuse slip with @Description  supporting document (file upload). @Tags         ExcuseSlips @Accept       multipart/form-data @Produce      json @Param        reason      formData string true "Reason for absence" @Param        absenceDate formData string true "Date (YYYY-MM-DD)" @Param        files       formData file   true "Supporting Document" @Success      201         {object} map[string]interface{} @Failure      400         {object} map[string]string @Failure      403         {object} map[string]string @Failure      500         {object} map[string]string @Router       /slips [post] PostSlip handles the submission of student excuse slips.*

**Constructors/Factory Functions:**
- `NewHandler`

### `ListSlipRequest`
**Kind**: Struct

### `ListSlipsDTO`
**Kind**: Struct

### `MockRepositoryInterface`
MockRepositoryInterface is a mock of RepositoryInterface interface.

**Kind**: Struct

**Methods:**
- `BeginTx`
  - *BeginTx mocks base method.*
- `CheckStudentExistence`
  - *CheckStudentExistence mocks base method.*
- `CreateSlip`
  - *CreateSlip mocks base method.*
- `Delete`
  - *Delete mocks base method.*
- `DeleteSlipAttachments`
  - *DeleteSlipAttachments mocks base method.*
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetAll`
  - *GetAll mocks base method.*
- `GetAttachmentByID`
  - *GetAttachmentByID mocks base method.*
- `GetByIIRID`
  - *GetByIIRID mocks base method.*
- `GetByUserID`
  - *GetByUserID mocks base method.*
- `GetDB`
  - *GetDB mocks base method.*
- `GetSlipAttachments`
  - *GetSlipAttachments mocks base method.*
- `GetSlipByID`
  - *GetSlipByID mocks base method.*
- `GetSlipByIDWithDetails`
  - *GetSlipByIDWithDetails mocks base method.*
- `GetSlipCategories`
  - *GetSlipCategories mocks base method.*
- `GetSlipStats`
  - *GetSlipStats mocks base method.*
- `GetSlipStatuses`
  - *GetSlipStatuses mocks base method.*
- `GetTotalSlipsCount`
  - *GetTotalSlipsCount mocks base method.*
- `GetTotalUrgentSlipsCount`
  - *GetTotalUrgentSlipsCount mocks base method.*
- `GetUrgentSlips`
  - *GetUrgentSlips mocks base method.*
- `GetUserIDBySlipID`
  - *GetUserIDBySlipID mocks base method.*
- `SaveSlipAttachment`
  - *SaveSlipAttachment mocks base method.*
- `UpdateSlip`
  - *UpdateSlip mocks base method.*
- `UpdateStatus`
  - *UpdateStatus mocks base method.*
- `WithTransaction`
  - *WithTransaction mocks base method.*

**Constructors/Factory Functions:**
- `NewMockRepositoryInterface`

### `MockRepositoryInterfaceMockRecorder`
MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.

**Kind**: Struct

**Methods:**
- `BeginTx`
  - *BeginTx indicates an expected call of BeginTx.*
- `CheckStudentExistence`
  - *CheckStudentExistence indicates an expected call of CheckStudentExistence.*
- `CreateSlip`
  - *CreateSlip indicates an expected call of CreateSlip.*
- `Delete`
  - *Delete indicates an expected call of Delete.*
- `DeleteSlipAttachments`
  - *DeleteSlipAttachments indicates an expected call of DeleteSlipAttachments.*
- `GetAll`
  - *GetAll indicates an expected call of GetAll.*
- `GetAttachmentByID`
  - *GetAttachmentByID indicates an expected call of GetAttachmentByID.*
- `GetByIIRID`
  - *GetByIIRID indicates an expected call of GetByIIRID.*
- `GetByUserID`
  - *GetByUserID indicates an expected call of GetByUserID.*
- `GetDB`
  - *GetDB indicates an expected call of GetDB.*
- `GetSlipAttachments`
  - *GetSlipAttachments indicates an expected call of GetSlipAttachments.*
- `GetSlipByID`
  - *GetSlipByID indicates an expected call of GetSlipByID.*
- `GetSlipByIDWithDetails`
  - *GetSlipByIDWithDetails indicates an expected call of GetSlipByIDWithDetails.*
- `GetSlipCategories`
  - *GetSlipCategories indicates an expected call of GetSlipCategories.*
- `GetSlipStats`
  - *GetSlipStats indicates an expected call of GetSlipStats.*
- `GetSlipStatuses`
  - *GetSlipStatuses indicates an expected call of GetSlipStatuses.*
- `GetTotalSlipsCount`
  - *GetTotalSlipsCount indicates an expected call of GetTotalSlipsCount.*
- `GetTotalUrgentSlipsCount`
  - *GetTotalUrgentSlipsCount indicates an expected call of GetTotalUrgentSlipsCount.*
- `GetUrgentSlips`
  - *GetUrgentSlips indicates an expected call of GetUrgentSlips.*
- `GetUserIDBySlipID`
  - *GetUserIDBySlipID indicates an expected call of GetUserIDBySlipID.*
- `SaveSlipAttachment`
  - *SaveSlipAttachment indicates an expected call of SaveSlipAttachment.*
- `UpdateSlip`
  - *UpdateSlip indicates an expected call of UpdateSlip.*
- `UpdateStatus`
  - *UpdateStatus indicates an expected call of UpdateStatus.*
- `WithTransaction`
  - *WithTransaction indicates an expected call of WithTransaction.*

### `MockServiceInterface`
MockServiceInterface is a mock of ServiceInterface interface.

**Kind**: Struct

**Methods:**
- `DownloadAttachment`
  - *DownloadAttachment mocks base method.*
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetAllExcuseSlips`
  - *GetAllExcuseSlips mocks base method.*
- `GetAttachmentFile`
  - *GetAttachmentFile mocks base method.*
- `GetExcuseSlipsByIIRID`
  - *GetExcuseSlipsByIIRID mocks base method.*
- `GetSlipAttachments`
  - *GetSlipAttachments mocks base method.*
- `GetSlipByID`
  - *GetSlipByID mocks base method.*
- `GetSlipCategories`
  - *GetSlipCategories mocks base method.*
- `GetSlipStats`
  - *GetSlipStats mocks base method.*
- `GetSlipStatuses`
  - *GetSlipStatuses mocks base method.*
- `GetUrgentSlips`
  - *GetUrgentSlips mocks base method.*
- `SubmitExcuseSlip`
  - *SubmitExcuseSlip mocks base method.*
- `UpdateExcuseSlip`
  - *UpdateExcuseSlip mocks base method.*
- `UpdateExcuseSlipStatus`
  - *UpdateExcuseSlipStatus mocks base method.*
- `WithTransaction`
  - *WithTransaction mocks base method.*

**Constructors/Factory Functions:**
- `NewMockServiceInterface`

### `MockServiceInterfaceMockRecorder`
MockServiceInterfaceMockRecorder is the mock recorder for MockServiceInterface.

**Kind**: Struct

**Methods:**
- `DownloadAttachment`
  - *DownloadAttachment indicates an expected call of DownloadAttachment.*
- `GetAllExcuseSlips`
  - *GetAllExcuseSlips indicates an expected call of GetAllExcuseSlips.*
- `GetAttachmentFile`
  - *GetAttachmentFile indicates an expected call of GetAttachmentFile.*
- `GetExcuseSlipsByIIRID`
  - *GetExcuseSlipsByIIRID indicates an expected call of GetExcuseSlipsByIIRID.*
- `GetSlipAttachments`
  - *GetSlipAttachments indicates an expected call of GetSlipAttachments.*
- `GetSlipByID`
  - *GetSlipByID indicates an expected call of GetSlipByID.*
- `GetSlipCategories`
  - *GetSlipCategories indicates an expected call of GetSlipCategories.*
- `GetSlipStats`
  - *GetSlipStats indicates an expected call of GetSlipStats.*
- `GetSlipStatuses`
  - *GetSlipStatuses indicates an expected call of GetSlipStatuses.*
- `GetUrgentSlips`
  - *GetUrgentSlips indicates an expected call of GetUrgentSlips.*
- `SubmitExcuseSlip`
  - *SubmitExcuseSlip indicates an expected call of SubmitExcuseSlip.*
- `UpdateExcuseSlip`
  - *UpdateExcuseSlip indicates an expected call of UpdateExcuseSlip.*
- `UpdateExcuseSlipStatus`
  - *UpdateExcuseSlipStatus indicates an expected call of UpdateExcuseSlipStatus.*
- `WithTransaction`
  - *WithTransaction indicates an expected call of WithTransaction.*

### `Repository`
**Kind**: Struct

**Methods:**
- `BeginTx`
- `CheckStudentExistence`
- `CreateSlip`
- `Delete`
- `DeleteSlipAttachments`
- `GetAll`
- `GetAttachmentByID`
- `GetByIIRID`
- `GetByUserID`
- `GetDB`
- `GetSlipAttachments`
- `GetSlipByID`
- `GetSlipByIDWithDetails`
- `GetSlipCategories`
- `GetSlipStats`
- `GetSlipStatuses`
- `GetTotalSlipsCount`
- `GetTotalUrgentSlipsCount`
- `GetUrgentSlips`
- `GetUserIDBySlipID`
- `SaveSlipAttachment`
- `UpdateSlip`
- `UpdateStatus`
- `WithTransaction`

### `RepositoryInterface`
RepositoryInterface defines the data access layer for managing excuse slips.

**Kind**: Interface

**Constructors/Factory Functions:**
- `NewRepository`

### `Service`
**Kind**: Struct

**Methods:**
- `DownloadAttachment`
  - *DownloadAttachment streams the attachment from Azure Blob Storage.*
- `GetAllExcuseSlips`
- `GetAttachmentFile`
- `GetExcuseSlipsByIIRID`
- `GetSlipAttachments`
- `GetSlipByID`
- `GetSlipCategories`
- `GetSlipStats`
- `GetSlipStatuses`
- `GetUrgentSlips`
- `GetUserIDBySlipID`
- `SubmitExcuseSlip`
  - *SubmitExcuseSlip creates a new slip with attachments.*
- `UpdateExcuseSlip`
- `UpdateExcuseSlipStatus`

### `ServiceInterface`
ServiceInterface defines the business logic for managing excuse slips.

**Kind**: Interface

**Constructors/Factory Functions:**
- `NewService`

### `Slip`
**Kind**: Struct

**Methods:**
- `ToPersistence`

### `SlipAttachment`
**Kind**: Struct

**Methods:**
- `ToPersistence`

### `SlipAttachmentDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Attachment mappers*

### `SlipCategory`
**Kind**: Struct

**Methods:**
- `ToPersistence`

### `SlipCategoryDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Category mappers*

### `SlipDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Slip mappers*

### `SlipDTO`
**Kind**: Struct

### `SlipStatus`
**Kind**: Struct

**Methods:**
- `ToPersistence`

### `SlipStatusCount`
**Kind**: Struct

### `SlipStatusCountDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Status Count mappers*

### `SlipStatusDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Status mappers*

### `SlipWithDetailsView`
**Kind**: Struct

### `SlipWithDetailsViewDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Details View mappers*

### `UpdateStatusRequest`
**Kind**: Struct

## Package Level Functions

### `RegisterRoutes`
### `TestService_GetSlipByID`
### `TestService_UpdateExcuseSlipStatus`
