package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

func createStudent(
	index int,
	password string,
	userFromCSV *users.User,
	appointmentsDataset []map[string]string,
) {
	ctx := context.Background()
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	// generate core data needed later
	dob := gofakeit.DateRange(
		time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Now().AddDate(-18, 0, 0))
	birthYear := dob.Year()

	var user users.User
	if userFromCSV != nil {
		user = *userFromCSV
	} else {
		studentEmail := fmt.Sprintf("student%d@gmail.com",
			index+1) // guarantee unique

		// users
		user = users.User{
			ID:         uuid.New().String(),
			Roles:      []users.Role{{ID: 1, Name: "Student"}},
			FirstName:  gofakeit.FirstName(),
			MiddleName: randomMiddleName(),
			LastName:   gofakeit.LastName(),
			SuffixName: nullStringIf(rand.Float32() < 0.1, gofakeit.RandomString([]string{"Jr.", "Sr.", "III", "IV"})),
			Email:      studentEmail,
			PasswordHash: structs.NullableString{
				Valid:  true,
				String: password,
			},
			AuthType: "native",
			IsActive: true,
		}
	}

	err = usersRepo.CreateUser(ctx, tx, user)
	if err != nil {
		log.Fatalf("[Seeder] {Create Student User}: %v", err)
	}

	insertNotifications(ctx, tx, user.ID)

	generateFullStudentIIR(
		ctx,
		tx,
		user.ID,
		dob,
		birthYear,
		index,
		appointmentsDataset,
	)

	tx.Commit()
}

func generateFullStudentIIR(
	ctx context.Context,
	tx *sqlx.Tx,
	userID string,
	dob time.Time,
	birthYear int,
	index int,
	appointmentsDataset []map[string]string,
) {
	iir := students.IIRRecord{
		ID:          uuid.New().String(),
		UserID:      userID,
		IsSubmitted: true,
	}

	iirID, err := studentsRepo.UpsertIIRRecord(ctx, tx, &iir)
	if err != nil {
		log.Fatalf("[Seeder] {Upsert IIRRecord}: %v", err)
	}

	// addresses (residential & provincial)
	resAddr1 := insertAddress(ctx, tx, true)
	resAddr2 := insertAddress(ctx, tx, false)
	insertStudentAddress(ctx, tx, iirID, resAddr1, "Residential")
	insertStudentAddress(ctx, tx, iirID, resAddr2, "Provincial")

	// related persons (father, mother, optional guardian)
	father := insertRelatedPerson(ctx, tx)
	mother := insertRelatedPerson(ctx, tx)
	father.AddressID = &resAddr1
	mother.AddressID = &resAddr1

	guardianScenario := pickGuardianScenario()
	var guardian *relatedPersonSeed

	switch guardianScenario {
	case "father_guardian":
		father.AddressID = &resAddr2
	case "mother_guardian":
		mother.AddressID = &resAddr2
	case "separate_guardian":
		g := insertRelatedPerson(ctx, tx)
		g.AddressID = &resAddr2
		guardian = &g
	}

	// link related persons
	linkFamilyMembers(
		ctx,
		tx,
		iirID,
		father,
		mother,
		guardian,
		guardianScenario,
	)

	// student_personal_info
	emergency := deriveEmergencyContact(father, mother, guardian,
		guardianScenario, resAddr1, resAddr2)
	emergencyContactID := insertEmergencyContact(ctx, tx, iirID,
		emergency)
	insertPersonalInfo(ctx, tx, iirID, dob, index,
		emergencyContactID)

	// family background
	familyBgID := insertFamilyBackground(ctx, tx, iirID)

	// sibling supports (if employed siblings > 0)
	insertSiblingSupports(ctx, tx, familyBgID)

	// educational background
	ebID := insertEducationalBackground(ctx, tx, iirID)

	// school details for each educational level
	insertSchoolDetails(ctx, tx, ebID, birthYear, index)

	// health records
	insertStudentHealthRecords(ctx, tx, iirID)

	// consultations
	insertConsultations(ctx, tx, iirID)

	// test results
	insertTestResults(ctx, tx, iirID)

	// finances
	sfID := insertStudentFinances(ctx, tx, iirID)

	// financial supports
	insertFinancialSupports(ctx, tx, sfID)

	// activities
	insertActivities(ctx, tx, iirID)

	// subject preferences
	insertSubjectPreferences(ctx, tx, iirID)

	// hobbies
	insertHobbies(ctx, tx, iirID)

	// Collect appointment and admission slip IDs for notes
	appointmentIDs := []string{}
	admissionSlipIDs := []string{}

	// admission slip (30% chance)
	if rand.Float32() < 0.3 {
		slipID := insertAdmissionSlip(ctx, tx, iirID)
		if slipID != "" {
			admissionSlipIDs = append(admissionSlipIDs,
				slipID)
		}
	}

	// appointment (30% chance)
	if rand.Float32() < 0.3 {
		for i := 0; i < rand.Intn(5)+1; i++ {
			// up to 5 appointments per student
			apptID := insertAppointment(ctx, tx, iirID, appointmentsDataset)
			if apptID != "" {
				appointmentIDs = append(appointmentIDs,
					apptID)
			}
		}
	}

	// significant notes (after appointments/slips created)
	insertSignificantNotes(ctx, tx, iirID, appointmentIDs,
		admissionSlipIDs)

	fmt.Printf("Created student %d %s\n", index+1, iirID)
}

func linkFamilyMembers(
	ctx context.Context,
	tx *sqlx.Tx,
	iirID string,
	father, mother relatedPersonSeed,
	guardian *relatedPersonSeed,
	scenario string,
) {
	switch scenario {
	case "no_guardian":
		linkRelatedPerson(ctx, tx, iirID, father.ID, "Father",
			true, false, true)
		linkRelatedPerson(ctx, tx, iirID, mother.ID, "Mother",
			true, false, true)
	case "father_guardian":
		linkRelatedPerson(ctx, tx, iirID, father.ID, "Father",
			true, true, true)
		linkRelatedPerson(ctx, tx, iirID, mother.ID, "Mother",
			true, false, true)
	case "mother_guardian":
		linkRelatedPerson(ctx, tx, iirID, father.ID, "Father",
			true, false, true)
		linkRelatedPerson(ctx, tx, iirID, mother.ID, "Mother",
			true, true, true)
	case "separate_guardian":
		linkRelatedPerson(ctx, tx, iirID, father.ID, "Father",
			true, false, true)
		linkRelatedPerson(ctx, tx, iirID, mother.ID, "Mother",
			true, false, true)
		if guardian != nil {
			linkRelatedPerson(ctx, tx, iirID, guardian.ID,
				"Legal Guardian", false, true, true)
		}
	}
}

func pickGuardianScenario() string {
	r := rand.Float32()
	if r < 0.25 {
		return "no_guardian"
	}
	if r < 0.50 {
		return "father_guardian"
	}
	if r < 0.75 {
		return "mother_guardian"
	}
	return "separate_guardian"
}

func deriveEmergencyContact(
	father, mother relatedPersonSeed,
	guardian *relatedPersonSeed,
	guardianScenario string,
	residentialAddressID, provincialAddressID int,
) emergencyContactSeed {
	switch guardianScenario {
	case "father_guardian":
		return emergencyContactSeed{
			FirstName:      father.FirstName,
			MiddleName:     father.MiddleName,
			LastName:       father.LastName,
			Number:         validContact(father.ContactNumber),
			RelationshipID: relationshipID("Father"),
			AddressID: safeAddressID(
				father.AddressID,
				provincialAddressID,
			),
		}
	case "mother_guardian":
		return emergencyContactSeed{
			FirstName:      mother.FirstName,
			MiddleName:     mother.MiddleName,
			LastName:       mother.LastName,
			Number:         validContact(mother.ContactNumber),
			RelationshipID: relationshipID("Mother"),
			AddressID: safeAddressID(
				mother.AddressID,
				provincialAddressID,
			),
		}
	case "separate_guardian":
		return emergencyContactSeed{
			FirstName:      guardian.FirstName,
			MiddleName:     guardian.MiddleName,
			LastName:       guardian.LastName,
			Number:         validContact(guardian.ContactNumber),
			RelationshipID: relationshipID("Legal Guardian"),
			AddressID: safeAddressID(
				guardian.AddressID,
				provincialAddressID,
			),
		}
	default:
		// no guardian: choose one of the parents as emergency contact
		if rand.Intn(2) == 0 {
			return emergencyContactSeed{
				FirstName:      father.FirstName,
				MiddleName:     father.MiddleName,
				LastName:       father.LastName,
				Number:         validContact(father.ContactNumber),
				RelationshipID: relationshipID("Father"),
				AddressID: safeAddressID(
					father.AddressID,
					residentialAddressID,
				),
			}
		}
		return emergencyContactSeed{
			FirstName:      mother.FirstName,
			MiddleName:     mother.MiddleName,
			LastName:       mother.LastName,
			Number:         validContact(mother.ContactNumber),
			RelationshipID: relationshipID("Mother"),
			AddressID: safeAddressID(
				mother.AddressID,
				residentialAddressID,
			),
		}
	}
}

func safeAddressID(addrID *int, fallback int) int {
	if addrID != nil {
		return *addrID
	}
	return fallback
}

func validContact(contact structs.NullableString) string {
	if contact.Valid && strings.TrimSpace(contact.String) != "" {
		return contact.String
	}
	return gofakeit.Phone()
}

func insertEmergencyContact(
	ctx context.Context,
	tx *sqlx.Tx,
	iirID string,
	emergency emergencyContactSeed,
) int {
	ec := &students.EmergencyContact{
		IIRID:          iirID,
		FirstName:      emergency.FirstName,
		MiddleName:     emergency.MiddleName,
		LastName:       emergency.LastName,
		ContactNumber:  emergency.Number,
		RelationshipID: emergency.RelationshipID,
		AddressID:      emergency.AddressID,
	}

	id, err := studentsRepo.UpsertEmergencyContact(ctx, tx, ec)
	if err != nil {
		log.Fatalf("[Seeder] {Upsert EmergencyContact}: %v", err)
	}
	return id
}

func insertPersonalInfo(
	ctx context.Context,
	tx *sqlx.Tx,
	iirID string,
	dob time.Time,
	studentIndex int,
	emergencyContactID int,
) {
	statusID := studentStatusByName["active"]
	isEmployed := studentIndex%2 == 0
	yearLevel := rand.Intn(4) + 1
	enrollmentYear := time.Now().Year() - yearLevel

	studentNumber := fmt.Sprintf(
		"%d-%05d-TG-%d",
		enrollmentYear,
		rand.Intn(100000),
		statusID,
	)
	var employerName, employerAddress structs.NullableString
	if isEmployed {
		empName := gofakeit.Company()
		empAddr := gofakeit.Address().Address
		employerName = structs.NullableString{String: empName, Valid: true}
		employerAddress = structs.NullableString{String: empAddr, Valid: true}
	}

	mobileNumber := gofakeit.Phone()
	telephoneNumber := structs.NullableString{Valid: false}
	if studentIndex%3 != 0 {
		telephoneNumber = structs.NullableString{
			String: gofakeit.Phone(),
			Valid:  true,
		}
	}

	civilStatusID := chooseCivilStatusID()
	genderID := randomChoice(genderIDs).(int)

	info := &students.StudentPersonalInfo{
		IIRID:           iirID,
		StudentNumber:   studentNumber,
		GenderID:        genderID,
		CivilStatusID:   civilStatusID,
		ReligionID:      randomChoice(religionIDs).(int),
		HeightM:         gofakeit.Float64Range(1.37, 1.98),
		WeightKg:        gofakeit.Float64Range(40, 100),
		Complexion:      gofakeit.Color(),
		HighSchoolGWA:   gofakeit.Float64Range(75, 98),
		CourseID:        randomChoice(courseIDs).(int),
		YearLevel:       yearLevel,
		Section:         rand.Intn(5) + 1,
		PlaceOfBirth:    gofakeit.City(),
		DateOfBirth:     dob.Format("2006-01-02"),
		IsEmployed:      isEmployed,
		EmployerName:    employerName,
		EmployerAddress: employerAddress,
		MobileNumber:    mobileNumber,
		TelephoneNumber: telephoneNumber,
		StatusID:        statusID,
	}

	err := studentsRepo.UpsertStudentPersonalInfo(ctx, tx, info)
	if err != nil {
		log.Fatalf("[Seeder] {Upsert PersonalInfo}: %v", err)
	}
}

func insertAddress(ctx context.Context, tx *sqlx.Tx, isResidential bool) int {
	// Pick a random city that has at least one barangay, along with its
	// region/province
	var cityRow struct {
		CityCode     string                 `db:"city_code"`
		RegionCode   string                 `db:"region_code"`
		ProvinceCode structs.NullableString `db:"province_code"`
	}

	query := `
		SELECT c.code AS city_code, c.region_code, c.province_code
		FROM cities c
		INNER JOIN barangays b ON c.code = b.city_code
		WHERE c.region_code IS NOT NULL`

	if isResidential {
		query += ` AND c.region_code = '1300000000'`
	}

	query += `
		GROUP BY c.code, c.region_code, c.province_code
		ORDER BY RAND() LIMIT 1
	`
	err := tx.Get(&cityRow, query)
	if err != nil {
		log.Fatal(
			"No cities with barangays found. Please run the address "+
				"seeder first (make locations): ",
			err,
		)
	}

	// Fetch a random barangay in that city
	var barangayCode string
	err = tx.Get(
		&barangayCode,
		"SELECT code FROM barangays WHERE city_code = ? ORDER BY RAND() LIMIT 1",
		cityRow.CityCode,
	)
	if err != nil {
		log.Fatal("No barangays found for city "+cityRow.CityCode+": ", err)
	}

	street := gofakeit.Street()

	// province_code is nullable (e.g. NCR has no provinces)
	var provinceCode *string
	if cityRow.ProvinceCode.Valid {
		s := cityRow.ProvinceCode.String
		provinceCode = &s
	}

	addr := &locations.Address{
		RegionCode:   cityRow.RegionCode,
		ProvinceCode: structs.PointerToNullableString(provinceCode),
		CityCode:     cityRow.CityCode,
		BarangayCode: barangayCode,
		StreetDetail: structs.PointerToNullableString(&street),
	}

	id, err := locationsRepo.UpsertAddress(ctx, tx, addr)
	if err != nil {
		log.Fatalf("[Seeder] {Upsert Address}: %v", err)
	}
	return id
}

func insertStudentAddress(
	ctx context.Context,
	tx *sqlx.Tx,
	iirID string,
	addressID int,
	addrType string,
) {
	sa := &students.StudentAddress{
		IIRID:       iirID,
		AddressID:   addressID,
		AddressType: addrType,
	}
	_, err := studentsRepo.UpsertStudentAddress(ctx, tx, sa)
	if err != nil {
		log.Fatalf("[Seeder] {Upsert StudentAddress}: %v", err)
	}
}

func insertRelatedPerson(ctx context.Context, tx *sqlx.Tx) relatedPersonSeed {
	fname := gofakeit.FirstName()
	lname := gofakeit.LastName()
	mname := randomMiddleName()
	dob := gofakeit.DateRange(
		time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Now().AddDate(-30, 0, 0),
	)
	educID := randomEducationalAttainmentID()
	occupation := structs.NullableString{
		String: gofakeit.JobTitle(),
		Valid:  gofakeit.Bool(),
	}
	employer := structs.NullableString{
		String: gofakeit.Company(),
		Valid:  occupation.Valid,
	}
	employerAddr := structs.NullableString{
		String: gofakeit.Address().Address,
		Valid:  occupation.Valid,
	}
	contact := structs.NullableString{
		String: gofakeit.Phone(),
		Valid:  gofakeit.Bool(),
	}

	rp := &students.RelatedPerson{
		EducationalAttainmentID: structs.Int64ToNullableInt64(educID),
		DateOfBirth:             dob.Format("2006-01-02"),
		LastName:                lname,
		FirstName:               fname,
		MiddleName:              mname,
		Occupation:              occupation,
		EmployerName:            employer,
		EmployerAddress:         employerAddr,
	}

	id, err := studentsRepo.UpsertRelatedPerson(ctx, tx, rp)
	if err != nil {
		log.Fatalf("[Seeder] {Upsert RelatedPerson}: %v", err)
	}

	return relatedPersonSeed{
		ID:            id,
		FirstName:     fname,
		MiddleName:    mname,
		LastName:      lname,
		ContactNumber: contact,
		AddressID:     nil,
	}
}

func randomEducationalAttainmentID() int64 {
	// IDs match migration 000016_add_educational_attainments_table.up.sql
	// 1: Elementary, 2: Junior High School, 3: Senior High School, 4: College, 5: Post Grad, 6: Vocational, etc.
	return int64(rand.Intn(7) + 1)
}

func linkRelatedPerson(
	ctx context.Context,
	tx *sqlx.Tx,
	iirID string,
	personID int,
	relType string,
	isParent, isGuardian, isLiving bool,
) {
	relID := relationshipID(relType)
	srp := &students.StudentRelatedPerson{
		IIRID:           iirID,
		RelatedPersonID: personID,
		RelationshipID:  relID,
		IsParent:        isParent,
		IsGuardian:      isGuardian,
		IsLiving:        isLiving,
	}
	err := studentsRepo.UpsertStudentRelatedPerson(ctx, tx, srp)
	if err != nil {
		log.Fatalf("[Seeder] {Upsert StudentRelatedPerson}: %v", err)
	}
}

func insertFamilyBackground(
	ctx context.Context,
	tx *sqlx.Tx,
	iirID string,
) int {
	brothers := rand.Intn(4)
	sisters := rand.Intn(4)
	employedSibs := rand.Intn(brothers + sisters + 1)
	ordinal := rand.Intn(brothers+sisters+1) + 1
	parentalID := randomChoice(parentalStatusIDs).(int)
	var details structs.NullableString
	if parentalID == 5 { // "Other"
		details = structs.NullableString{
			String: gofakeit.Sentence(3),
			Valid:  true,
		}
	}
	quiet := gofakeit.Bool()
	sharing := gofakeit.Bool()
	var shareDetails structs.NullableString
	if sharing {
		shareDetails = structs.NullableString{
			String: "Shares with " + gofakeit.FirstName(),
			Valid:  true,
		}
	}
	natureID := randomChoice(natureOfResidenceIDs).(int)

	fb := &students.FamilyBackground{
		IIRID:                 iirID,
		ParentalStatusID:      parentalID,
		ParentalStatusDetails: details,
		Brothers:              brothers,
		Sisters:               sisters,
		EmployedSiblings:      employedSibs,
		OrdinalPosition:       ordinal,
		HaveQuietPlaceToStudy: quiet,
		IsSharingRoom:         sharing,
		RoomSharingDetails:    shareDetails,
		NatureOfResidenceId:   natureID,
	}

	id, err := studentsRepo.UpsertFamilyBackground(ctx, tx, fb)
	if err != nil {
		log.Fatalf("[Seeder] {Upsert FamilyBackground}: %v", err)
	}
	return id
}

func insertSiblingSupports(ctx context.Context, tx *sqlx.Tx, familyBgID int) {
	// only if employed siblings > 0, maybe add 1-2 supports
	if rand.Intn(3) == 0 { // 1/3 chance
		num := rand.Intn(2) + 1
		selected := make(map[int]bool)
		for len(selected) < num {
			stID := randomChoice(siblingSupportTypeIDs).(int)
			selected[stID] = true
		}
		for stID := range selected {
			sss := &students.StudentSiblingSupport{
				FamilyBackgroundID: familyBgID,
				SupportTypeID:      stID,
			}
			err := studentsRepo.CreateStudentSiblingSupport(ctx, tx, sss)
			if err != nil {
				log.Fatalf("[Seeder] {Create StudentSiblingSupport}: %v", err)
			}
		}
	}
}

func insertEducationalBackground(
	ctx context.Context,
	tx *sqlx.Tx,
	iirID string,
) int {
	nature := "Continuous"
	var details structs.NullableString
	if rand.Float32() < 0.1 { // interrupted
		nature = "Interrupted"
		details = structs.NullableString{
			String: gofakeit.Sentence(5),
			Valid:  true,
		}
	}

	eb := &students.EducationalBackground{
		IIRID:              iirID,
		NatureOfSchooling:  nature,
		InterruptedDetails: details,
	}

	id, err := studentsRepo.UpsertEducationalBackground(ctx, tx, eb)
	if err != nil {
		log.Fatalf("[Seeder] {Upsert EducationalBackground}: %v", err)
	}
	return id
}

func insertSchoolDetails(
	ctx context.Context,
	tx *sqlx.Tx,
	ebID, birthYear, studentIndex int,
) {
	nowYear := time.Now().Year()

	levelStartAge := map[string]int{
		"Pre-Elementary":     4,
		"Elementary":         6,
		"Junior High School": 12,
		"Senior High School": 16,
		"College":            18,
	}

	levelPath := []string{"pre-elementary", "elementary", "high school"}
	if studentIndex%3 == 0 {
		levelPath = append(levelPath, "vocational")
	} else {
		levelPath = append(levelPath, "college")
	}

	for _, levelKey := range levelPath {
		levelID, ok := educationalLevelByName[levelKey]
		if !ok {
			continue
		}

		levelName := cases.Title(language.English).String(levelKey)
		if levelKey == "pre-elementary" {
			levelName = "Pre-Elementary"
		}
		if levelKey == "high school" {
			levelName = "High School"
		}

		startAge := levelStartAge[levelName]
		startYear := birthYear + startAge
		duration := 0
		switch levelName {
		case "Pre-Elementary":
			duration = 2
		case "Elementary":
			duration = 6
		case "High School":
			duration = 4
		case "Vocational":
			duration = 2
		case "College":
			duration = 4
		}
		completedYear := startYear + duration
		if completedYear > nowYear {
			completedYear = nowYear + rand.Intn(3)
		}
		schoolType := randomChoice([]string{"Private", "Public"}).(string)
		awards := randomAwards()

		sd := &students.SchoolDetails{
			EBID:               ebID,
			EducationalLevelID: levelID,
			SchoolName:         gofakeit.Company(),
			SchoolAddress:      gofakeit.Address().Address,
			SchoolType:         schoolType,
			YearStarted:        startYear,
			YearCompleted:      completedYear,
			Awards:             awards,
		}

		_, err := studentsRepo.UpsertSchoolDetails(ctx, tx, sd)
		if err != nil {
			log.Fatalf("[Seeder] {Upsert SchoolDetails}: %v", err)
		}
	}
}

func insertStudentHealthRecords(
	ctx context.Context,
	tx *sqlx.Tx,
	iirID string,
) {
	visionProb := false
	hearingProb := false
	speechProb := false
	genProb := false

	roll := rand.Float32()
	if roll >= 0.35 {
		impairments := []string{"vision", "hearing", "speech", "general"}
		rand.Shuffle(len(impairments), func(i, j int) {
			impairments[i], impairments[j] = impairments[j], impairments[i]
		})

		count := 1
		if roll >= 0.75 {
			count = 2
		}

		for _, impairment := range impairments[:count] {
			switch impairment {
			case "vision":
				visionProb = true
			case "hearing":
				hearingProb = true
			case "speech":
				speechProb = true
			case "general":
				genProb = true
			}
		}
	}

	visionDet := nullStringIf(visionProb, gofakeit.Sentence(3))
	hearingDet := nullStringIf(hearingProb, gofakeit.Sentence(3))
	speechDet := nullStringIf(speechProb, gofakeit.Sentence(3))
	genDet := nullStringIf(genProb, gofakeit.Sentence(5))

	hr := &students.StudentHealthRecord{
		IIRID:                   iirID,
		VisionHasProblem:        visionProb,
		VisionDetails:           visionDet,
		HearingHasProblem:       hearingProb,
		HearingDetails:          hearingDet,
		SpeechHasProblem:        speechProb,
		SpeechDetails:           speechDet,
		GeneralHealthHasProblem: genProb,
		GeneralHealthDetails:    genDet,
	}

	_, err := studentsRepo.UpsertStudentHealthRecord(ctx, tx, hr)
	if err != nil {
		log.Fatalf("[Seeder] {Upsert HealthRecord}: %v", err)
	}
}

func insertConsultations(ctx context.Context, tx *sqlx.Tx, iirID string) {
	if rand.Float32() < 0.55 {
		num := rand.Intn(3) + 1
		for i := 0; i < num; i++ {
			profType := randomChoice(
				[]string{"Psychiatrist", "Psychologist", "Counselor"},
			).(string)
			has := true
			when := sql.NullString{
				String: gofakeit.Date().Format("2006-01-02"),
				Valid:  true,
			}
			what := sql.NullString{String: gofakeit.Sentence(8), Valid: true}

			sc := &students.StudentConsultation{
				IIRID:            iirID,
				ProfessionalType: profType,
				HasConsulted:     has,
				WhenDate:         structs.FromSqlNull(when),
				ForWhat:          structs.FromSqlNull(what),
			}

			_, err := studentsRepo.UpsertStudentConsultation(ctx, tx, sc)
			if err != nil {
				log.Fatalf("[Seeder] {Upsert Consultation}: %v", err)
			}
		}
	}
}

func insertTestResults(ctx context.Context, tx *sqlx.Tx, iirID string) {
	numOfTests := 2
	for i := 0; i < numOfTests; i++ {
		tr := &students.TestResult{
			IIRID:       iirID,
			TestDate:    gofakeit.Date().Format("2006-01-02"),
			TestName:    gofakeit.JobTitle() + " Test",
			RawScore:    fmt.Sprintf("%d", rand.Intn(50)+50),
			Percentile:  fmt.Sprintf("%d", rand.Intn(30)+70),
			Description: gofakeit.Sentence(5),
		}

		// Since we don't have a specific UpsertTestResult in the interface yet
		// (based on view),
		// we use raw exec for now, but following the pattern.
		// Actually, I should check if I missed it in the interface.
		// Interface shows DeleteTestResultsByIIRID but not Upsert.
		// I will stick to raw exec for those NOT in repository yet to avoid
		// breaking interface.
		// Wait, I see UpsertStudentHealthRecord etc but not TestResults.
		_, err := tx.Exec(`
			INSERT INTO test_results (
				iir_id, test_date, test_name, raw_score, percentile, description
			) VALUES (?, ?, ?, ?, ?, ?)
		`, tr.IIRID, tr.TestDate, tr.TestName, tr.RawScore, tr.Percentile, tr.Description)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func insertSignificantNotes(ctx context.Context, tx *sqlx.Tx, iirID string,
	appointmentIDs, admissionSlipIDs []string,
) {
	if rand.Float32() < 0.3 {
		num := rand.Intn(3) + 1
		for i := 0; i < num; i++ {
			var appointmentID sql.NullString
			var admissionSlipID sql.NullString

			if rand.Float32() < 0.5 && len(appointmentIDs) > 0 {
				appointmentID = sql.NullString{
					String: appointmentIDs[rand.Intn(len(appointmentIDs))],
					Valid:  true,
				}
			}

			if rand.Float32() < 0.5 && len(admissionSlipIDs) > 0 {
				admissionSlipID = sql.NullString{
					String: admissionSlipIDs[rand.Intn(len(admissionSlipIDs))],
					Valid:  true,
				}
			}

			_, err := tx.Exec(`
				INSERT INTO significant_notes (
					id, iir_id, appointment_id, admission_slip_id,
					note, remarks
				) VALUES (?, ?, ?, ?, ?, ?)
			`, uuid.New().String(), iirID, appointmentID, admissionSlipID,
				gofakeit.Sentence(8),
				gofakeit.Sentence(5))
			if err != nil {
				log.Printf("[Seeder] {Insert SignificantNote}: %v",
					err)
				return
			}
		}
	}
}

func insertStudentFinances(ctx context.Context, tx *sqlx.Tx, iirID string) int {
	incRangeID := randomChoice(incomeRangeIDs).(int)
	var otherInc sql.NullString
	if incRangeID == 10 {
		otherInc = sql.NullString{
			String: "Business income",
			Valid:  gofakeit.Bool(),
		}
	}
	allowance := gofakeit.Float64Range(100, 1500)

	sf := &students.StudentFinance{
		IIRID:           iirID,
		IncomeRangeID:   incRangeID,
		OtherIncome:     structs.FromSqlNull(otherInc),
		WeeklyAllowance: allowance,
	}

	id, err := studentsRepo.UpsertStudentFinance(ctx, tx, sf)
	if err != nil {
		log.Fatalf("[Seeder] {Upsert StudentFinance}: %v", err)
	}
	return id
}

func insertFinancialSupports(ctx context.Context, tx *sqlx.Tx, sfID int) {
	num := rand.Intn(3) + 1
	selected := make(map[int]bool)
	for len(selected) < num {
		stID := randomChoice(supportTypeIDs).(int)
		selected[stID] = true
	}
	for stID := range selected {
		sfs := &students.StudentFinancialSupport{
			StudentFinanceID: sfID,
			SupportTypeID:    stID,
		}
		err := studentsRepo.CreateStudentFinancialSupport(ctx, tx, sfs)
		if err != nil {
			log.Fatalf("[Seeder] {Create StudentFinancialSupport}: %v", err)
		}
	}
}

func insertActivities(ctx context.Context, tx *sqlx.Tx, iirID string) {
	if len(activityOptionIDs) == 0 {
		return
	}
	num := rand.Intn(4)
	for i := 0; i < num; i++ {
		optID := randomChoice(activityOptionIDs).(int)
		role := randomChoice([]string{"Officer", "Member", "Other"}).(string)
		var roleSpec structs.NullableString
		if role == "Other" {
			roleSpec = structs.NullableString{
				String: gofakeit.Word(),
				Valid:  true,
			}
		}

		sa := &students.StudentActivity{
			IIRID:             iirID,
			OptionID:          optID,
			Role:              role,
			RoleSpecification: roleSpec,
		}

		_, err := studentsRepo.CreateStudentActivity(ctx, tx, sa)
		if err != nil {
			log.Fatalf("[Seeder] {Create StudentActivity}: %v", err)
		}
	}
}

func insertSubjectPreferences(ctx context.Context, tx *sqlx.Tx, iirID string) {
	if rand.Float32() < 0.2 {
		return
	}

	subjects := []string{
		"Math",
		"Science",
		"English",
		"History",
		"PE",
		"Art",
		"Music",
		"Computer",
	}
	num := rand.Intn(5) + 2
	selected := pickUniqueStrings(subjects, num)
	fav := gofakeit.Bool()

	for _, sub := range selected {
		ssp := &students.StudentSubjectPreference{
			IIRID:       iirID,
			SubjectName: sub,
			IsFavorite:  fav,
		}
		_, err := studentsRepo.CreateStudentSubjectPreference(ctx, tx, ssp)
		if err != nil {
			log.Fatalf("[Seeder] {Create StudentSubjectPreference}: %v", err)
		}
	}
}

func insertHobbies(ctx context.Context, tx *sqlx.Tx, iirID string) {
	num := rand.Intn(4) + 1
	for i := 1; i <= num; i++ {
		rank := i
		sh := &students.StudentHobby{
			IIRID:        iirID,
			HobbyName:    gofakeit.Hobby(),
			PriorityRank: rank,
		}
		_, err := studentsRepo.CreateStudentHobby(ctx, tx, sh)
		if err != nil {
			log.Fatalf("[Seeder] {Create StudentHobby}: %v", err)
		}
	}
}

func randomAwards() structs.NullableString {
	awardPool := []string{
		"With Honors",
		"Best in Mathematics",
		"Best in Conduct",
		"Leadership Award",
		"Academic Excellence",
		"Science Fair Winner",
	}

	r := rand.Float32()
	if r < 0.35 {
		return structs.NullableString{Valid: false}
	}

	if r < 0.7 {
		return structs.NullableString{
			String: awardPool[rand.Intn(len(awardPool))],
			Valid:  true,
		}
	}

	count := rand.Intn(3) + 2
	selected := pickUniqueStrings(awardPool, count)
	return structs.NullableString{
		String: strings.Join(selected, ", "),
		Valid:  true,
	}
}

func insertNotifications(ctx context.Context, tx *sqlx.Tx, userID string) {
	count := rand.Intn(3) + 3

	for i := 0; i < count; i++ {
		title := gofakeit.Sentence(3)
		message := gofakeit.Sentence(10)
		randomType := notificationTypes[rand.Intn(len(notificationTypes))]
		isRead := rand.Intn(2)

		_, err := tx.Exec(`
            INSERT INTO notifications (id, receiver_id, title, message, type, is_read, created_at)
            VALUES (?, ?, ?, ?, ?, ?, NOW())
        `, uuid.New().String(), userID, title, message, randomType, isRead)
		if err != nil {
			log.Printf(
				"failed to insert fake notification for %s with type %s: %v",
				userID,
				randomType,
				err,
			)
		}
	}
}
