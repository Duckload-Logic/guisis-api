# Feature Module: `students`

## Overview
Package students is a generated GoMock package.

## Types and Interfaces

### `ActivityOption`
ActivityOption represents the domain entity for student activities.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `ActivityOptionDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *ActivityOption mappers*

### `BulkUpdateStatusRequest`
**Kind**: Struct

### `CivilStatusType`
CivilStatusType represents the domain entity for civil status types.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `CivilStatusTypeDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *CivilStatusType mappers*

### `ComprehensiveProfileDTO`
**Kind**: Struct

### `Course`
Course represents the domain entity for academic courses.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `CourseDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Course mappers*

### `EducationalBGDTO`
**Kind**: Struct

### `EducationalBackground`
EducationalBackground represents the student's prior schooling status.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `EducationalBackgroundDB`
Education and Background

**Kind**: Struct

**Methods:**
- `ToDomain`
  - *EducationalBackground mapper*

### `EducationalBackgroundDTO`
**Kind**: Struct

### `EducationalLevel`
EducationalLevel represents the domain entity for educational levels.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `EducationalLevelDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *EducationalLevel mappers*

### `EmergencyContact`
EmergencyContact represents the domain entity for emergency notifications.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `EmergencyContactDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *EmergencyContact mappers*

### `EmergencyContactDTO`
**Kind**: Struct

### `EnrollmentReason`
EnrollmentReason represents the domain entity for enrollment reasons.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `EnrollmentReasonDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *EnrollmentReason mappers*

### `FamilyBackground`
FamilyBackground represents details about the student's family structure.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `FamilyBackgroundDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *FamilyBackground mapper*

### `FamilyBackgroundDTO`
**Kind**: Struct

### `Gender`
Gender represents the domain entity for student gender.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `GenderDB`
Lookup models

**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Gender mappers*

### `Handler`
**Kind**: Struct

**Methods:**
- `GetActivityOptions`
- `GetCivilStatusTypes`
- `GetCourses`
- `GetEducationalLevels`
- `GetEnrollmentReasons`
- `GetGenders`
- `GetIncomeRanges`
- `GetNatureOfResidenceTypes`
- `GetParentalStatusTypes`
- `GetReligions`
- `GetSiblingSupportTypes`
- `GetStudentActivities`
- `GetStudentAddresses`
- `GetStudentBasicInfo`
- `GetStudentCORByUserID`
- `GetStudentCORs`
- `GetStudentConsultations`
- `GetStudentEducationalBackground`
- `GetStudentEnrollmentReasons`
- `GetStudentFamilyBackground`
- `GetStudentFinancialInfo`
- `GetStudentHealthRecord`
- `GetStudentHobbies`
- `GetStudentIIRByIIRID`
- `GetStudentIIRByUserID`
- `GetStudentIIRDraft`
- `GetStudentIIRPDF`
  - *GenerateIIR godoc @Summary      Generate Student IIR PDF @Description  Generates and downloads the student's Initial Interview Record as a PDF. @Tags         Students @Produce      application/pdf @Param        iirID   path      string  true  "IIR ID" @Success      200     {file}    binary @Failure      400     {object}  map[string]string @Failure      500     {object}  map[string]string @Router       /students/{iirID}/iir/download [get]*
- `GetStudentList`
  - *GetStudentList godoc @Summary      List Students @Description  Retrieves a paginated list of students with optional filters. @Tags         Students @Accept       json @Produce      json @Param        course     query     string  false  "Filter by course" @Param        year_level query     int     false  "Filter by year level" @Param		gender_id query     int     false  "Filter by gender ID" @Param        page       query     int     false  "Page number" default(1) @Param        page_size  query     int     false  "Number of records per page" default(20) @Success      200        {object}  ListStudentsResponse @Failure      500        {object}  map[string]string "Failed to list students" @Router       /students [get]*
- `GetStudentPersonalInfo`
- `GetStudentProfile`
- `GetStudentRelatedPersons`
- `GetStudentRelationshipTypes`
- `GetStudentStatuses`
- `GetStudentSubjectPreferences`
- `GetStudentSupportTypes`
- `GetStudentTestResults`
- `PatchStudentStatusBulk`
  - *PatchStudentBulkStatus handles PATCH /students/inventory/records/bulk-status. It applies a lifecycle status transition to multiple student records in a single transaction. For the "Graduated" status the service layer enforces eligibility (Diploma/Year-3, Bachelor/Year-4) and quietly skips ineligible records — the response always reports the attempted count so the frontend can surface a warning if the actual count differs.*
- `PostStudentCOR`
- `PostStudentIIR`
- `PostStudentIIRDraft`

**Constructors/Factory Functions:**
- `NewHandler`

### `IIRDraft`
IIRDraft represents a student's in-progress Initial Interview Record.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `IIRDraftDB`
Core Student Records

**Kind**: Struct

**Methods:**
- `ToDomain`
  - *IIRDraft mappers*

### `IIRRecord`
IIRRecord represents a submitted Initial Interview Record entry.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `IIRRecordDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *IIRRecord mappers*

### `IncomeRange`
IncomeRange represents the domain entity for family income ranges.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `IncomeRangeDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *IncomeRange mappers*

### `ListStudentsRequest`
List Students

**Kind**: Struct

### `ListStudentsResponse`
**Kind**: Struct

### `MockRepositoryInterface`
MockRepositoryInterface is a mock of RepositoryInterface interface.

**Kind**: Struct

**Methods:**
- `BeginTx`
  - *BeginTx mocks base method.*
- `BulkUpdateStudentStatus`
  - *BulkUpdateStudentStatus mocks base method.*
- `CreateStudentActivity`
  - *CreateStudentActivity mocks base method.*
- `CreateStudentFinancialSupport`
  - *CreateStudentFinancialSupport mocks base method.*
- `CreateStudentHobby`
  - *CreateStudentHobby mocks base method.*
- `CreateStudentSelectedReason`
  - *CreateStudentSelectedReason mocks base method.*
- `CreateStudentSiblingSupport`
  - *CreateStudentSiblingSupport mocks base method.*
- `CreateStudentSubjectPreference`
  - *CreateStudentSubjectPreference mocks base method.*
- `DeleteSchoolDetailsByEBID`
  - *DeleteSchoolDetailsByEBID mocks base method.*
- `DeleteSignificantNotesByIIRID`
  - *DeleteSignificantNotesByIIRID mocks base method.*
- `DeleteStudentActivitiesByIIRID`
  - *DeleteStudentActivitiesByIIRID mocks base method.*
- `DeleteStudentFinancialSupportsByFinanceID`
  - *DeleteStudentFinancialSupportsByFinanceID mocks base method.*
- `DeleteStudentHobbiesByIIRID`
  - *DeleteStudentHobbiesByIIRID mocks base method.*
- `DeleteStudentRelatedPersons`
  - *DeleteStudentRelatedPersons mocks base method.*
- `DeleteStudentSelectedReasons`
  - *DeleteStudentSelectedReasons mocks base method.*
- `DeleteStudentSiblingSupportsByFamilyID`
  - *DeleteStudentSiblingSupportsByFamilyID mocks base method.*
- `DeleteStudentSubjectPreferencesByIIRID`
  - *DeleteStudentSubjectPreferencesByIIRID mocks base method.*
- `DeleteTestResultsByIIRID`
  - *DeleteTestResultsByIIRID mocks base method.*
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GetActivityOptionByID`
  - *GetActivityOptionByID mocks base method.*
- `GetActivityOptions`
  - *GetActivityOptions mocks base method.*
- `GetCivilStatusByID`
  - *GetCivilStatusByID mocks base method.*
- `GetCivilStatusTypes`
  - *GetCivilStatusTypes mocks base method.*
- `GetCourseByID`
  - *GetCourseByID mocks base method.*
- `GetCourses`
  - *GetCourses mocks base method.*
- `GetDB`
  - *GetDB mocks base method.*
- `GetEducationalLevelByID`
  - *GetEducationalLevelByID mocks base method.*
- `GetEducationalLevels`
  - *GetEducationalLevels mocks base method.*
- `GetEmergencyContactByIIRID`
  - *GetEmergencyContactByIIRID mocks base method.*
- `GetEnrollmentReasonByID`
  - *GetEnrollmentReasonByID mocks base method.*
- `GetEnrollmentReasons`
  - *GetEnrollmentReasons mocks base method.*
- `GetFinancialSupportTypeByFinanceID`
  - *GetFinancialSupportTypeByFinanceID mocks base method.*
- `GetGenderByID`
  - *GetGenderByID mocks base method.*
- `GetGenders`
  - *GetGenders mocks base method.*
- `GetIIRDraftByUserID`
  - *GetIIRDraftByUserID mocks base method.*
- `GetIncomeRangeByID`
  - *GetIncomeRangeByID mocks base method.*
- `GetIncomeRanges`
  - *GetIncomeRanges mocks base method.*
- `GetNatureOfResidenceByID`
  - *GetNatureOfResidenceByID mocks base method.*
- `GetNatureOfResidenceTypes`
  - *GetNatureOfResidenceTypes mocks base method.*
- `GetParentalStatusByID`
  - *GetParentalStatusByID mocks base method.*
- `GetParentalStatusTypes`
  - *GetParentalStatusTypes mocks base method.*
- `GetRelatedPersonByID`
  - *GetRelatedPersonByID mocks base method.*
- `GetReligionByID`
  - *GetReligionByID mocks base method.*
- `GetReligions`
  - *GetReligions mocks base method.*
- `GetSchoolDetailsByEBID`
  - *GetSchoolDetailsByEBID mocks base method.*
- `GetSiblingSupportTypeByID`
  - *GetSiblingSupportTypeByID mocks base method.*
- `GetSiblingSupportTypes`
  - *GetSiblingSupportTypes mocks base method.*
- `GetStudentActivities`
  - *GetStudentActivities mocks base method.*
- `GetStudentAddresses`
  - *GetStudentAddresses mocks base method.*
- `GetStudentBasicInfo`
  - *GetStudentBasicInfo mocks base method.*
- `GetStudentCORByUserID`
  - *GetStudentCORByUserID mocks base method.*
- `GetStudentCORsByUserID`
  - *GetStudentCORsByUserID mocks base method.*
- `GetStudentConsultations`
  - *GetStudentConsultations mocks base method.*
- `GetStudentEducationalBackground`
  - *GetStudentEducationalBackground mocks base method.*
- `GetStudentEnrollmentReasons`
  - *GetStudentEnrollmentReasons mocks base method.*
- `GetStudentFamilyBackground`
  - *GetStudentFamilyBackground mocks base method.*
- `GetStudentFinancialInfo`
  - *GetStudentFinancialInfo mocks base method.*
- `GetStudentHealthRecord`
  - *GetStudentHealthRecord mocks base method.*
- `GetStudentHobbies`
  - *GetStudentHobbies mocks base method.*
- `GetStudentIIR`
  - *GetStudentIIR mocks base method.*
- `GetStudentIIRByUserID`
  - *GetStudentIIRByUserID mocks base method.*
- `GetStudentPersonalInfo`
  - *GetStudentPersonalInfo mocks base method.*
- `GetStudentRelatedPersons`
  - *GetStudentRelatedPersons mocks base method.*
- `GetStudentRelationshipByID`
  - *GetStudentRelationshipByID mocks base method.*
- `GetStudentRelationshipTypes`
  - *GetStudentRelationshipTypes mocks base method.*
- `GetStudentSiblingSupport`
  - *GetStudentSiblingSupport mocks base method.*
- `GetStudentStatuses`
  - *GetStudentStatuses mocks base method.*
- `GetStudentSubjectPreferences`
  - *GetStudentSubjectPreferences mocks base method.*
- `GetStudentSupportByID`
  - *GetStudentSupportByID mocks base method.*
- `GetStudentSupportTypes`
  - *GetStudentSupportTypes mocks base method.*
- `GetStudentTestResults`
  - *GetStudentTestResults mocks base method.*
- `GetTotalStudentsCount`
  - *GetTotalStudentsCount mocks base method.*
- `IsStudentLocked`
  - *IsStudentLocked mocks base method.*
- `ListStudents`
  - *ListStudents mocks base method.*
- `SaveStudentCOR`
  - *SaveStudentCOR mocks base method.*
- `UpsertEducationalBackground`
  - *UpsertEducationalBackground mocks base method.*
- `UpsertEmergencyContact`
  - *UpsertEmergencyContact mocks base method.*
- `UpsertFamilyBackground`
  - *UpsertFamilyBackground mocks base method.*
- `UpsertIIRDraft`
  - *UpsertIIRDraft mocks base method.*
- `UpsertIIRRecord`
  - *UpsertIIRRecord mocks base method.*
- `UpsertRelatedPerson`
  - *UpsertRelatedPerson mocks base method.*
- `UpsertSchoolDetails`
  - *UpsertSchoolDetails mocks base method.*
- `UpsertStudentAddress`
  - *UpsertStudentAddress mocks base method.*
- `UpsertStudentConsultation`
  - *UpsertStudentConsultation mocks base method.*
- `UpsertStudentFinance`
  - *UpsertStudentFinance mocks base method.*
- `UpsertStudentHealthRecord`
  - *UpsertStudentHealthRecord mocks base method.*
- `UpsertStudentPersonalInfo`
  - *UpsertStudentPersonalInfo mocks base method.*
- `UpsertStudentRelatedPerson`
  - *UpsertStudentRelatedPerson mocks base method.*
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
- `BulkUpdateStudentStatus`
  - *BulkUpdateStudentStatus indicates an expected call of BulkUpdateStudentStatus.*
- `CreateStudentActivity`
  - *CreateStudentActivity indicates an expected call of CreateStudentActivity.*
- `CreateStudentFinancialSupport`
  - *CreateStudentFinancialSupport indicates an expected call of CreateStudentFinancialSupport.*
- `CreateStudentHobby`
  - *CreateStudentHobby indicates an expected call of CreateStudentHobby.*
- `CreateStudentSelectedReason`
  - *CreateStudentSelectedReason indicates an expected call of CreateStudentSelectedReason.*
- `CreateStudentSiblingSupport`
  - *CreateStudentSiblingSupport indicates an expected call of CreateStudentSiblingSupport.*
- `CreateStudentSubjectPreference`
  - *CreateStudentSubjectPreference indicates an expected call of CreateStudentSubjectPreference.*
- `DeleteSchoolDetailsByEBID`
  - *DeleteSchoolDetailsByEBID indicates an expected call of DeleteSchoolDetailsByEBID.*
- `DeleteSignificantNotesByIIRID`
  - *DeleteSignificantNotesByIIRID indicates an expected call of DeleteSignificantNotesByIIRID.*
- `DeleteStudentActivitiesByIIRID`
  - *DeleteStudentActivitiesByIIRID indicates an expected call of DeleteStudentActivitiesByIIRID.*
- `DeleteStudentFinancialSupportsByFinanceID`
  - *DeleteStudentFinancialSupportsByFinanceID indicates an expected call of DeleteStudentFinancialSupportsByFinanceID.*
- `DeleteStudentHobbiesByIIRID`
  - *DeleteStudentHobbiesByIIRID indicates an expected call of DeleteStudentHobbiesByIIRID.*
- `DeleteStudentRelatedPersons`
  - *DeleteStudentRelatedPersons indicates an expected call of DeleteStudentRelatedPersons.*
- `DeleteStudentSelectedReasons`
  - *DeleteStudentSelectedReasons indicates an expected call of DeleteStudentSelectedReasons.*
- `DeleteStudentSiblingSupportsByFamilyID`
  - *DeleteStudentSiblingSupportsByFamilyID indicates an expected call of DeleteStudentSiblingSupportsByFamilyID.*
- `DeleteStudentSubjectPreferencesByIIRID`
  - *DeleteStudentSubjectPreferencesByIIRID indicates an expected call of DeleteStudentSubjectPreferencesByIIRID.*
- `DeleteTestResultsByIIRID`
  - *DeleteTestResultsByIIRID indicates an expected call of DeleteTestResultsByIIRID.*
- `GetActivityOptionByID`
  - *GetActivityOptionByID indicates an expected call of GetActivityOptionByID.*
- `GetActivityOptions`
  - *GetActivityOptions indicates an expected call of GetActivityOptions.*
- `GetCivilStatusByID`
  - *GetCivilStatusByID indicates an expected call of GetCivilStatusByID.*
- `GetCivilStatusTypes`
  - *GetCivilStatusTypes indicates an expected call of GetCivilStatusTypes.*
- `GetCourseByID`
  - *GetCourseByID indicates an expected call of GetCourseByID.*
- `GetCourses`
  - *GetCourses indicates an expected call of GetCourses.*
- `GetDB`
  - *GetDB indicates an expected call of GetDB.*
- `GetEducationalLevelByID`
  - *GetEducationalLevelByID indicates an expected call of GetEducationalLevelByID.*
- `GetEducationalLevels`
  - *GetEducationalLevels indicates an expected call of GetEducationalLevels.*
- `GetEmergencyContactByIIRID`
  - *GetEmergencyContactByIIRID indicates an expected call of GetEmergencyContactByIIRID.*
- `GetEnrollmentReasonByID`
  - *GetEnrollmentReasonByID indicates an expected call of GetEnrollmentReasonByID.*
- `GetEnrollmentReasons`
  - *GetEnrollmentReasons indicates an expected call of GetEnrollmentReasons.*
- `GetFinancialSupportTypeByFinanceID`
  - *GetFinancialSupportTypeByFinanceID indicates an expected call of GetFinancialSupportTypeByFinanceID.*
- `GetGenderByID`
  - *GetGenderByID indicates an expected call of GetGenderByID.*
- `GetGenders`
  - *GetGenders indicates an expected call of GetGenders.*
- `GetIIRDraftByUserID`
  - *GetIIRDraftByUserID indicates an expected call of GetIIRDraftByUserID.*
- `GetIncomeRangeByID`
  - *GetIncomeRangeByID indicates an expected call of GetIncomeRangeByID.*
- `GetIncomeRanges`
  - *GetIncomeRanges indicates an expected call of GetIncomeRanges.*
- `GetNatureOfResidenceByID`
  - *GetNatureOfResidenceByID indicates an expected call of GetNatureOfResidenceByID.*
- `GetNatureOfResidenceTypes`
  - *GetNatureOfResidenceTypes indicates an expected call of GetNatureOfResidenceTypes.*
- `GetParentalStatusByID`
  - *GetParentalStatusByID indicates an expected call of GetParentalStatusByID.*
- `GetParentalStatusTypes`
  - *GetParentalStatusTypes indicates an expected call of GetParentalStatusTypes.*
- `GetRelatedPersonByID`
  - *GetRelatedPersonByID indicates an expected call of GetRelatedPersonByID.*
- `GetReligionByID`
  - *GetReligionByID indicates an expected call of GetReligionByID.*
- `GetReligions`
  - *GetReligions indicates an expected call of GetReligions.*
- `GetSchoolDetailsByEBID`
  - *GetSchoolDetailsByEBID indicates an expected call of GetSchoolDetailsByEBID.*
- `GetSiblingSupportTypeByID`
  - *GetSiblingSupportTypeByID indicates an expected call of GetSiblingSupportTypeByID.*
- `GetSiblingSupportTypes`
  - *GetSiblingSupportTypes indicates an expected call of GetSiblingSupportTypes.*
- `GetStudentActivities`
  - *GetStudentActivities indicates an expected call of GetStudentActivities.*
- `GetStudentAddresses`
  - *GetStudentAddresses indicates an expected call of GetStudentAddresses.*
- `GetStudentBasicInfo`
  - *GetStudentBasicInfo indicates an expected call of GetStudentBasicInfo.*
- `GetStudentCORByUserID`
  - *GetStudentCORByUserID indicates an expected call of GetStudentCORByUserID.*
- `GetStudentCORsByUserID`
  - *GetStudentCORsByUserID indicates an expected call of GetStudentCORsByUserID.*
- `GetStudentConsultations`
  - *GetStudentConsultations indicates an expected call of GetStudentConsultations.*
- `GetStudentEducationalBackground`
  - *GetStudentEducationalBackground indicates an expected call of GetStudentEducationalBackground.*
- `GetStudentEnrollmentReasons`
  - *GetStudentEnrollmentReasons indicates an expected call of GetStudentEnrollmentReasons.*
- `GetStudentFamilyBackground`
  - *GetStudentFamilyBackground indicates an expected call of GetStudentFamilyBackground.*
- `GetStudentFinancialInfo`
  - *GetStudentFinancialInfo indicates an expected call of GetStudentFinancialInfo.*
- `GetStudentHealthRecord`
  - *GetStudentHealthRecord indicates an expected call of GetStudentHealthRecord.*
- `GetStudentHobbies`
  - *GetStudentHobbies indicates an expected call of GetStudentHobbies.*
- `GetStudentIIR`
  - *GetStudentIIR indicates an expected call of GetStudentIIR.*
- `GetStudentIIRByUserID`
  - *GetStudentIIRByUserID indicates an expected call of GetStudentIIRByUserID.*
- `GetStudentPersonalInfo`
  - *GetStudentPersonalInfo indicates an expected call of GetStudentPersonalInfo.*
- `GetStudentRelatedPersons`
  - *GetStudentRelatedPersons indicates an expected call of GetStudentRelatedPersons.*
- `GetStudentRelationshipByID`
  - *GetStudentRelationshipByID indicates an expected call of GetStudentRelationshipByID.*
- `GetStudentRelationshipTypes`
  - *GetStudentRelationshipTypes indicates an expected call of GetStudentRelationshipTypes.*
- `GetStudentSiblingSupport`
  - *GetStudentSiblingSupport indicates an expected call of GetStudentSiblingSupport.*
- `GetStudentStatuses`
  - *GetStudentStatuses indicates an expected call of GetStudentStatuses.*
- `GetStudentSubjectPreferences`
  - *GetStudentSubjectPreferences indicates an expected call of GetStudentSubjectPreferences.*
- `GetStudentSupportByID`
  - *GetStudentSupportByID indicates an expected call of GetStudentSupportByID.*
- `GetStudentSupportTypes`
  - *GetStudentSupportTypes indicates an expected call of GetStudentSupportTypes.*
- `GetStudentTestResults`
  - *GetStudentTestResults indicates an expected call of GetStudentTestResults.*
- `GetTotalStudentsCount`
  - *GetTotalStudentsCount indicates an expected call of GetTotalStudentsCount.*
- `IsStudentLocked`
  - *IsStudentLocked indicates an expected call of IsStudentLocked.*
- `ListStudents`
  - *ListStudents indicates an expected call of ListStudents.*
- `SaveStudentCOR`
  - *SaveStudentCOR indicates an expected call of SaveStudentCOR.*
- `UpsertEducationalBackground`
  - *UpsertEducationalBackground indicates an expected call of UpsertEducationalBackground.*
- `UpsertEmergencyContact`
  - *UpsertEmergencyContact indicates an expected call of UpsertEmergencyContact.*
- `UpsertFamilyBackground`
  - *UpsertFamilyBackground indicates an expected call of UpsertFamilyBackground.*
- `UpsertIIRDraft`
  - *UpsertIIRDraft indicates an expected call of UpsertIIRDraft.*
- `UpsertIIRRecord`
  - *UpsertIIRRecord indicates an expected call of UpsertIIRRecord.*
- `UpsertRelatedPerson`
  - *UpsertRelatedPerson indicates an expected call of UpsertRelatedPerson.*
- `UpsertSchoolDetails`
  - *UpsertSchoolDetails indicates an expected call of UpsertSchoolDetails.*
- `UpsertStudentAddress`
  - *UpsertStudentAddress indicates an expected call of UpsertStudentAddress.*
- `UpsertStudentConsultation`
  - *UpsertStudentConsultation indicates an expected call of UpsertStudentConsultation.*
- `UpsertStudentFinance`
  - *UpsertStudentFinance indicates an expected call of UpsertStudentFinance.*
- `UpsertStudentHealthRecord`
  - *UpsertStudentHealthRecord indicates an expected call of UpsertStudentHealthRecord.*
- `UpsertStudentPersonalInfo`
  - *UpsertStudentPersonalInfo indicates an expected call of UpsertStudentPersonalInfo.*
- `UpsertStudentRelatedPerson`
  - *UpsertStudentRelatedPerson indicates an expected call of UpsertStudentRelatedPerson.*
- `WithTransaction`
  - *WithTransaction indicates an expected call of WithTransaction.*

### `MockServiceInterface`
MockServiceInterface is a mock of ServiceInterface interface.

**Kind**: Struct

**Methods:**
- `BulkUpdateStudentStatus`
  - *BulkUpdateStudentStatus mocks base method.*
- `EXPECT`
  - *EXPECT returns an object that allows the caller to indicate expected use.*
- `GenerateIIR`
  - *GenerateIIR mocks base method.*
- `GetActivityOptions`
  - *GetActivityOptions mocks base method.*
- `GetCivilStatusTypes`
  - *GetCivilStatusTypes mocks base method.*
- `GetCourses`
  - *GetCourses mocks base method.*
- `GetEducationalBackground`
  - *GetEducationalBackground mocks base method.*
- `GetEducationalLevels`
  - *GetEducationalLevels mocks base method.*
- `GetEnrollmentReasons`
  - *GetEnrollmentReasons mocks base method.*
- `GetGenders`
  - *GetGenders mocks base method.*
- `GetIIRDraft`
  - *GetIIRDraft mocks base method.*
- `GetIncomeRanges`
  - *GetIncomeRanges mocks base method.*
- `GetNatureOfResidenceTypes`
  - *GetNatureOfResidenceTypes mocks base method.*
- `GetParentalStatusTypes`
  - *GetParentalStatusTypes mocks base method.*
- `GetReligions`
  - *GetReligions mocks base method.*
- `GetSiblingSupportTypes`
  - *GetSiblingSupportTypes mocks base method.*
- `GetStudentActivities`
  - *GetStudentActivities mocks base method.*
- `GetStudentAddresses`
  - *GetStudentAddresses mocks base method.*
- `GetStudentBasicInfo`
  - *GetStudentBasicInfo mocks base method.*
- `GetStudentCOR`
  - *GetStudentCOR mocks base method.*
- `GetStudentCORs`
  - *GetStudentCORs mocks base method.*
- `GetStudentConsultations`
  - *GetStudentConsultations mocks base method.*
- `GetStudentEnrollmentReasons`
  - *GetStudentEnrollmentReasons mocks base method.*
- `GetStudentFamilyBackground`
  - *GetStudentFamilyBackground mocks base method.*
- `GetStudentFinancialInfo`
  - *GetStudentFinancialInfo mocks base method.*
- `GetStudentHealthRecord`
  - *GetStudentHealthRecord mocks base method.*
- `GetStudentHobbies`
  - *GetStudentHobbies mocks base method.*
- `GetStudentIIR`
  - *GetStudentIIR mocks base method.*
- `GetStudentIIRByUserID`
  - *GetStudentIIRByUserID mocks base method.*
- `GetStudentPersonalInfo`
  - *GetStudentPersonalInfo mocks base method.*
- `GetStudentProfile`
  - *GetStudentProfile mocks base method.*
- `GetStudentRelatedPersons`
  - *GetStudentRelatedPersons mocks base method.*
- `GetStudentRelationshipTypes`
  - *GetStudentRelationshipTypes mocks base method.*
- `GetStudentStatuses`
  - *GetStudentStatuses mocks base method.*
- `GetStudentSubjectPreferences`
  - *GetStudentSubjectPreferences mocks base method.*
- `GetStudentSupportTypes`
  - *GetStudentSupportTypes mocks base method.*
- `GetStudentTestResults`
  - *GetStudentTestResults mocks base method.*
- `IsStudentLocked`
  - *IsStudentLocked mocks base method.*
- `ListStudents`
  - *ListStudents mocks base method.*
- `SaveIIRDraft`
  - *SaveIIRDraft mocks base method.*
- `SubmitCOR`
  - *SubmitCOR mocks base method.*
- `SubmitStudentIIR`
  - *SubmitStudentIIR mocks base method.*

**Constructors/Factory Functions:**
- `NewMockServiceInterface`

### `MockServiceInterfaceMockRecorder`
MockServiceInterfaceMockRecorder is the mock recorder for MockServiceInterface.

**Kind**: Struct

**Methods:**
- `BulkUpdateStudentStatus`
  - *BulkUpdateStudentStatus indicates an expected call of BulkUpdateStudentStatus.*
- `GenerateIIR`
  - *GenerateIIR indicates an expected call of GenerateIIR.*
- `GetActivityOptions`
  - *GetActivityOptions indicates an expected call of GetActivityOptions.*
- `GetCivilStatusTypes`
  - *GetCivilStatusTypes indicates an expected call of GetCivilStatusTypes.*
- `GetCourses`
  - *GetCourses indicates an expected call of GetCourses.*
- `GetEducationalBackground`
  - *GetEducationalBackground indicates an expected call of GetEducationalBackground.*
- `GetEducationalLevels`
  - *GetEducationalLevels indicates an expected call of GetEducationalLevels.*
- `GetEnrollmentReasons`
  - *GetEnrollmentReasons indicates an expected call of GetEnrollmentReasons.*
- `GetGenders`
  - *GetGenders indicates an expected call of GetGenders.*
- `GetIIRDraft`
  - *GetIIRDraft indicates an expected call of GetIIRDraft.*
- `GetIncomeRanges`
  - *GetIncomeRanges indicates an expected call of GetIncomeRanges.*
- `GetNatureOfResidenceTypes`
  - *GetNatureOfResidenceTypes indicates an expected call of GetNatureOfResidenceTypes.*
- `GetParentalStatusTypes`
  - *GetParentalStatusTypes indicates an expected call of GetParentalStatusTypes.*
- `GetReligions`
  - *GetReligions indicates an expected call of GetReligions.*
- `GetSiblingSupportTypes`
  - *GetSiblingSupportTypes indicates an expected call of GetSiblingSupportTypes.*
- `GetStudentActivities`
  - *GetStudentActivities indicates an expected call of GetStudentActivities.*
- `GetStudentAddresses`
  - *GetStudentAddresses indicates an expected call of GetStudentAddresses.*
- `GetStudentBasicInfo`
  - *GetStudentBasicInfo indicates an expected call of GetStudentBasicInfo.*
- `GetStudentCOR`
  - *GetStudentCOR indicates an expected call of GetStudentCOR.*
- `GetStudentCORs`
  - *GetStudentCORs indicates an expected call of GetStudentCORs.*
- `GetStudentConsultations`
  - *GetStudentConsultations indicates an expected call of GetStudentConsultations.*
- `GetStudentEnrollmentReasons`
  - *GetStudentEnrollmentReasons indicates an expected call of GetStudentEnrollmentReasons.*
- `GetStudentFamilyBackground`
  - *GetStudentFamilyBackground indicates an expected call of GetStudentFamilyBackground.*
- `GetStudentFinancialInfo`
  - *GetStudentFinancialInfo indicates an expected call of GetStudentFinancialInfo.*
- `GetStudentHealthRecord`
  - *GetStudentHealthRecord indicates an expected call of GetStudentHealthRecord.*
- `GetStudentHobbies`
  - *GetStudentHobbies indicates an expected call of GetStudentHobbies.*
- `GetStudentIIR`
  - *GetStudentIIR indicates an expected call of GetStudentIIR.*
- `GetStudentIIRByUserID`
  - *GetStudentIIRByUserID indicates an expected call of GetStudentIIRByUserID.*
- `GetStudentPersonalInfo`
  - *GetStudentPersonalInfo indicates an expected call of GetStudentPersonalInfo.*
- `GetStudentProfile`
  - *GetStudentProfile indicates an expected call of GetStudentProfile.*
- `GetStudentRelatedPersons`
  - *GetStudentRelatedPersons indicates an expected call of GetStudentRelatedPersons.*
- `GetStudentRelationshipTypes`
  - *GetStudentRelationshipTypes indicates an expected call of GetStudentRelationshipTypes.*
- `GetStudentStatuses`
  - *GetStudentStatuses indicates an expected call of GetStudentStatuses.*
- `GetStudentSubjectPreferences`
  - *GetStudentSubjectPreferences indicates an expected call of GetStudentSubjectPreferences.*
- `GetStudentSupportTypes`
  - *GetStudentSupportTypes indicates an expected call of GetStudentSupportTypes.*
- `GetStudentTestResults`
  - *GetStudentTestResults indicates an expected call of GetStudentTestResults.*
- `IsStudentLocked`
  - *IsStudentLocked indicates an expected call of IsStudentLocked.*
- `ListStudents`
  - *ListStudents indicates an expected call of ListStudents.*
- `SaveIIRDraft`
  - *SaveIIRDraft indicates an expected call of SaveIIRDraft.*
- `SubmitCOR`
  - *SubmitCOR indicates an expected call of SubmitCOR.*
- `SubmitStudentIIR`
  - *SubmitStudentIIR indicates an expected call of SubmitStudentIIR.*

### `NatureOfResidenceType`
NatureOfResidenceType represents the domain entity for residence types.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `NatureOfResidenceTypeDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *NatureOfResidenceType mappers*

### `ParentalStatusType`
ParentalStatusType represents the domain entity for parental status types.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `ParentalStatusTypeDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *ParentalStatusType mappers*

### `RelatedPerson`
RelatedPerson represents a person related to the student
(parent, guardian, etc.).

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `RelatedPersonDB`
Family and Related Persons

**Kind**: Struct

**Methods:**
- `ToDomain`
  - *RelatedPerson mapper*

### `RelatedPersonDTO`
**Kind**: Struct

### `Religion`
Religion represents the domain entity for religious affiliations.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `ReligionDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *Religion mappers*

### `Repository`
**Kind**: Struct

**Methods:**
- `BeginTx`
- `BulkUpdateStudentStatus`
  - *BulkUpdateStudentStatus updates status_id (and optionally graduation_year) for a set of students. When SelectAllMatching is true it builds a WHERE clause from the request filters instead of an IN list, so every matching record across all pages is affected in one query.*
- `CreateStudentActivity`
- `CreateStudentFinancialSupport`
- `CreateStudentHobby`
- `CreateStudentSelectedReason`
- `CreateStudentSiblingSupport`
- `CreateStudentSubjectPreference`
- `CreateTestResult`
- `DeleteSchoolDetailsByEBID`
- `DeleteSignificantNotesByIIRID`
- `DeleteStudentActivitiesByIIRID`
- `DeleteStudentFinancialSupportsByFinanceID`
- `DeleteStudentHobbiesByIIRID`
- `DeleteStudentRelatedPersons`
- `DeleteStudentSelectedReasons`
- `DeleteStudentSiblingSupportsByFamilyID`
- `DeleteStudentSubjectPreferencesByIIRID`
- `DeleteTestResultsByIIRID`
- `GetActivityOptionByID`
- `GetActivityOptions`
- `GetCivilStatusByID`
- `GetCivilStatusTypes`
- `GetCourseByID`
- `GetCourses`
- `GetDB`
- `GetEducationalLevelByID`
- `GetEducationalLevels`
- `GetEmergencyContactByIIRID`
- `GetEnrollmentReasonByID`
- `GetEnrollmentReasons`
- `GetFinancialSupportTypeByFinanceID`
- `GetGenderByID`
- `GetGenders`
  - *Lookup*
- `GetIIRDraftByUserID`
- `GetIncomeRangeByID`
- `GetIncomeRanges`
- `GetNatureOfResidenceByID`
- `GetNatureOfResidenceTypes`
- `GetParentalStatusByID`
- `GetParentalStatusTypes`
- `GetRelatedPersonByID`
- `GetReligionByID`
- `GetReligions`
- `GetSchoolDetailsByEBID`
- `GetSiblingSupportTypeByID`
- `GetSiblingSupportTypes`
- `GetStudentActivities`
- `GetStudentAddresses`
- `GetStudentBasicInfo`
- `GetStudentCORByUserID`
- `GetStudentCORsByUserID`
- `GetStudentConsultations`
- `GetStudentEducationalBackground`
- `GetStudentEnrollmentReasons`
- `GetStudentFamilyBackground`
- `GetStudentFinancialInfo`
- `GetStudentHealthRecord`
- `GetStudentHobbies`
- `GetStudentIIR`
- `GetStudentIIRByUserID`
- `GetStudentPersonalInfo`
- `GetStudentRelatedPersons`
- `GetStudentRelationshipByID`
- `GetStudentRelationshipTypes`
- `GetStudentSiblingSupport`
- `GetStudentStatuses`
  - *GetStudentStatuses retrieves all available student statuses.*
- `GetStudentSubjectPreferences`
- `GetStudentSupportByID`
- `GetStudentSupportTypes`
- `GetStudentTestResults`
- `GetTotalStudentsCount`
  - *Retrieve - Count*
- `IsStudentLocked`
  - *IsStudentLocked checks if a student's record is locked (Graduated, Archived, or Withdrawn).*
- `ListStudents`
  - *Retrieve - List*
- `SaveStudentCOR`
- `UpsertEducationalBackground`
- `UpsertEmergencyContact`
- `UpsertFamilyBackground`
- `UpsertIIRDraft`
  - *Save and Upsert*
- `UpsertIIRRecord`
- `UpsertRelatedPerson`
- `UpsertSchoolDetails`
- `UpsertStudentAddress`
- `UpsertStudentConsultation`
- `UpsertStudentFinance`
- `UpsertStudentHealthRecord`
- `UpsertStudentPersonalInfo`
- `UpsertStudentRelatedPerson`
- `WithTransaction`

### `RepositoryInterface`
**Kind**: Interface

**Constructors/Factory Functions:**
- `NewRepository`

### `SchoolDetails`
SchoolDetails represents specific schooling history entries.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `SchoolDetailsDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *SchoolDetails mapper*

### `SchoolDetailsDTO`
**Kind**: Struct

### `Service`
Service provides student-related business logic and data access.

**Kind**: Struct

**Methods:**
- `BulkUpdateStudentStatus`
  - *BulkUpdateStudentStatus applies a lifecycle status change to a batch of students. For the "Graduated" status a graduation eligibility guard runs first: only Diploma year-3 and Bachelor year-4 students are eligible. Any ineligible IDs are stripped from the request before delegation to the repository. The number of skipped records is returned as an error so the caller can surface a warning without aborting the whole operation.*
- `GenerateIIR`
  - *GenerateIIR generates a student's IIR as a PDF using an HTML template and Gotenberg.*
- `GetActivityOptions`
  - *GetActivityOptions retrieves all available student activity options.*
- `GetCivilStatusTypes`
  - *GetCivilStatusTypes retrieves all available civil status types.*
- `GetCourses`
  - *GetCourses retrieves all available academic courses.*
- `GetEducationalBackground`
  - *GetEducationalBackground retrieves the educational history for a student.*
- `GetEducationalLevels`
  - *GetEducationalLevels retrieves all available educational levels.*
- `GetEnrollmentReasons`
  - *GetEnrollmentReasons retrieves all available enrollment reason types.*
- `GetFormattedDate`
- `GetGenders`
  - *GetGenders retrieves all available gender types.*
- `GetIIRDraft`
  - *GetIIRDraft retrieves the latest IIR draft for a user.*
- `GetIncomeRanges`
  - *GetIncomeRanges retrieves all available family income range types.*
- `GetNatureOfResidenceTypes`
  - *GetNatureOfResidenceTypes retrieves all available nature of residence types.*
- `GetParentalStatusTypes`
  - *GetParentalStatusTypes retrieves all available parental status types.*
- `GetReligions`
  - *GetReligions retrieves all available religions.*
- `GetSiblingSupportTypes`
  - *GetSiblingSupportTypes retrieves all available sibling support types.*
- `GetStudentActivities`
  - *GetStudentActivities retrieves all extracurricular activities for a student.*
- `GetStudentAddresses`
  - *GetStudentAddresses retrieves all addresses associated with a student.*
- `GetStudentBasicInfo`
  - *GetStudentBasicInfo retrieves basic identification info for a student.*
- `GetStudentCOR`
- `GetStudentCORs`
  - *GetStudentCORs retrieves all COR records for a specific student.*
- `GetStudentConsultations`
  - *GetStudentConsultations retrieves past consultations for a student.*
- `GetStudentEnrollmentReasons`
  - *GetStudentEnrollmentReasons retrieves the reasons why a student enrolled.*
- `GetStudentFamilyBackground`
  - *GetStudentFamilyBackground retrieves the family background information for a student.*
- `GetStudentFinancialInfo`
  - *GetStudentFinancialInfo retrieves religious and financial info for a student.*
- `GetStudentHealthRecord`
  - *GetStudentHealthRecord retrieves the health record for a student.*
- `GetStudentHobbies`
  - *GetStudentHobbies retrieves hobbies and interests for a student.*
- `GetStudentIIR`
  - *GetStudentIIR retrieves a specific IIR record by its ID.*
- `GetStudentIIRByUserID`
  - *GetStudentIIRByUserID retrieves the IIR record for a specific user ID.*
- `GetStudentPersonalInfo`
  - *GetStudentPersonalInfo retrieves detailed personal information for a student.*
- `GetStudentProfile`
  - *GetStudentProfile retrieves the full comprehensive profile of a student.*
- `GetStudentRelatedPersons`
  - *GetStudentRelatedPersons retrieves all related persons (parents, guardians, etc.) for a student.*
- `GetStudentRelationshipTypes`
  - *GetStudentRelationshipTypes retrieves all available student relationship types.*
- `GetStudentStatuses`
  - *GetStudentStatuses retrieves all available student statuses.*
- `GetStudentSubjectPreferences`
  - *GetStudentSubjectPreferences retrieves academic subject preferences for a student.*
- `GetStudentSupportTypes`
  - *GetStudentSupportTypes retrieves all available student support types.*
- `GetStudentTestResults`
  - *GetStudentTestResults retrieves standardized test results for a student.*
- `IsStudentLocked`
  - *IsStudentLocked checks if a student's record is locked from further appointments or slips.*
- `ListStudents`
  - *Retrieve - List ListStudents retrieves a paginated list of student profiles.*
- `SaveIIRDraft`
  - *SaveIIRDraft persists a student's IIR draft data to the datastore.*
- `SubmitCOR`
  - *SubmitCOR uploads a Certificate of Registration and links it to the student.*
- `SubmitStudentIIR`
  - *SubmitStudentIIR processes and persists a full student IIR submission.*

### `ServiceInterface`
**Kind**: Interface

**Constructors/Factory Functions:**
- `NewService`

### `SibilingSupportType`
SibilingSupportType represents the domain entity for sibling support.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `SibilingSupportTypeDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *SibilingSupportType mappers*

### `StudentActivity`
StudentActivity represents student interests and club involvements.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `StudentActivityDB`
Interests and Activities

**Kind**: Struct

**Methods:**
- `ToDomain`
  - *StudentActivity mapper*

### `StudentActivityDTO`
**Kind**: Struct

### `StudentAddress`
StudentAddress represents the domain mapping for student locations.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `StudentAddressDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *StudentAddress mappers*

### `StudentAddressDTO`
**Kind**: Struct

### `StudentBasicInfoView`
StudentBasicInfoView represents a simplified view of student
identification info.

**Kind**: Struct

### `StudentBasicInfoViewDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *StudentBasicInfoView mappers*

### `StudentBasicInfoViewDTO`
**Kind**: Struct

### `StudentCOR`
StudentCOR represents the Certificate of Registration domain entity.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `StudentCORDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *StudentCOR mappers*

### `StudentConsultation`
StudentConsultation represents record of professional consultations.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `StudentConsultationDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *StudentConsultation mapper*

### `StudentConsultationDTO`
**Kind**: Struct

### `StudentFinance`
StudentFinance represents student financial background data.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `StudentFinanceDB`
Financial Support

**Kind**: Struct

**Methods:**
- `ToDomain`
  - *StudentFinance mapper*

### `StudentFinanceDTO`
**Kind**: Struct

### `StudentFinancialSupport`
StudentFinancialSupport maps financial support types to student finance.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `StudentFinancialSupportDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *StudentFinancialSupport mappers*

### `StudentHealthRecord`
StudentHealthRecord represents physical wellness parameters.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `StudentHealthRecordDB`
Health and Wellness

**Kind**: Struct

**Methods:**
- `ToDomain`
  - *StudentHealthRecord mapper*

### `StudentHealthRecordDTO`
**Kind**: Struct

### `StudentHobby`
StudentHobby represents extracurricular interests.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `StudentHobbyDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *StudentHobby mappers*

### `StudentHobbyDTO`
**Kind**: Struct

### `StudentPersonalInfo`
StudentPersonalInfo represents detailed student personal parameters.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `StudentPersonalInfoDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *StudentPersonalInfo mappers*

### `StudentPersonalInfoDTO`
**Kind**: Struct

### `StudentProfileDTO`
**Kind**: Struct

### `StudentProfileView`
StudentProfileView represents a summary view of a student's profile.

**Kind**: Struct

### `StudentProfileViewDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *StudentProfileView mappers*

### `StudentRelatedPerson`
StudentRelatedPerson represents the mapping between a student
and a related person.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `StudentRelatedPersonDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *StudentRelatedPerson mappers*

### `StudentRelationshipType`
StudentRelationshipType represents the domain entity for
relationships (e.g., Father, Mother).

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `StudentRelationshipTypeDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *StudentRelationshipType mappers*

### `StudentSelectedReason`
StudentSelectedReason represents the domain mapping for enrollment reasons.

**Kind**: Struct

### `StudentSelectedReasonDB`
Enrollment and Reasons

**Kind**: Struct

### `StudentSelectedReasonDTO`
**Kind**: Struct

### `StudentSiblingSupport`
StudentSiblingSupport maps family backgrounds to support types.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `StudentSiblingSupportDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *StudentSiblingSupport mappers*

### `StudentStatus`
StudentStatus represents the domain entity for academic status.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `StudentStatusDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *StudentStatus mappers*

### `StudentSubjectPreference`
StudentSubjectPreference represents academic subject interests.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `StudentSubjectPreferenceDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *StudentSubjectPreference mappers*

### `StudentSubjectPreferenceDTO`
**Kind**: Struct

### `StudentSupportType`
StudentSupportType represents the domain entity for student support types.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `StudentSupportTypeDB`
**Kind**: Struct

**Methods:**
- `ToDomain`
  - *StudentSupportType mappers*

### `TestResult`
TestResult represents academic assessment outcome records.

**Kind**: Struct

**Methods:**
- `ToPersistence`

### `TestResultDB`
Test Results and Assessments

**Kind**: Struct

**Methods:**
- `ToDomain`
  - *TestResult mappers*

### `TestResultDTO`
**Kind**: Struct

## Package Level Functions

### `RegisterRoutes`
### `TestService_GetStudentBasicInfo`
### `TestService_IsStudentLocked`
### `TestService_ListStudents`
