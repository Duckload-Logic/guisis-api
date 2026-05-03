package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/hash"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/slips"
)

func insertAdmissionSlip(
	ctx context.Context,
	tx *sqlx.Tx,
	iirID string,
) string {
	now := time.Now().Truncate(24 * time.Hour)

	// Randomly decide if this slip is for a past or future absence
	// Most are past (80%), some are upcoming (20%)
	var daysOffset int
	if rand.Float32() < 0.2 {
		daysOffset = rand.Intn(14) + 1 // 1-14 days in future
	} else {
		daysOffset = -(rand.Intn(90) + 1) // 1-90 days in past
	}

	absDate := now.AddDate(0, 0, daysOffset)
	dateOfAbsence := absDate.Format("2006-01-02")

	// Status depends on time
	statusName := chooseAdmissionSlipStatus(absDate, now)
	statusID, ok := admissionSlipStatusesByName[statusName]
	if !ok {
		statusID = randomChoice(admissionSlipStatusIDs).(int)
	}

	categoryID := randomChoice(admissionSlipCategoryIDs).(int)
	reason := generateRealisticAdmissionReason()

	// Date needed: submission should be near the absence date
	// for future, it's usually current; for past, it was around then
	var dateNeeded string
	if absDate.After(now) {
		dateNeeded = now.Format("2006-01-02")
	} else {
		dateNeeded = absDate.AddDate(0, 0, rand.Intn(3)+1).Format("2006-01-02")
	}

	// Use PDF as primary extension...
	extensions := []string{".pdf", ".pdf", ".pdf", ".jpg", ".jpeg", ".png"}
	ext := extensions[rand.Intn(len(extensions))]
	basePath := "./uploads"
	subFolder := gofakeit.UUID()

	baseFileName := generateAdmissionFileName()
	readableFileName := baseFileName + ext
	fileName := hash.GetSHA256Hash(readableFileName, 16) + ext

	fullStoragePath := filepath.Join(basePath, "slips", subFolder, fileName)

	dbURL := fmt.Sprintf("/slips/%s/%s", subFolder, fileName)
	dir := filepath.Dir(fullStoragePath)
	os.MkdirAll(dir, os.ModePerm)
	f, _ := os.Create(fullStoragePath)
	content := fmt.Sprintf(
		"ADMISSION SLIP / EXCUSE SLIP\nStudent ID: %s\nDate of Absence: %s\nReason: %s\n",
		iirID,
		dateOfAbsence,
		reason,
	)
	f.WriteString(content)
	f.Close()

	adminNotes := generateAdmissionNotes(statusName)
	admissionSlipID := uuid.New().String()

	slip := &slips.Slip{
		ID:            admissionSlipID,
		IIRID:         iirID,
		Reason:        reason,
		CategoryID:    categoryID,
		DateOfAbsence: dateOfAbsence,
		DateNeeded:    dateNeeded,
		StatusID:      statusID,
		AdminNotes:    adminNotes,
	}

	_, err := slipsRepo.CreateSlip(ctx, tx, slip)
	if err != nil {
		log.Printf("[Seeder] {Insert AdmissionSlip}: %v", err)
		return ""
	}

	fileID := uuid.New().String()
	// Insert dummy record into files table so FK check passes
	_, err = tx.ExecContext(ctx, `
		INSERT INTO files (id, file_name, file_url, file_type, file_size, mime_type)
		VALUES (?, ?, ?, ?, ?, ?)
	`, fileID, readableFileName, dbURL, "application/pdf", 1024, "application/pdf")
	if err != nil {
		log.Printf("[Seeder] {Insert File Metadata}: %v", err)
		return ""
	}

	attachment := &slips.SlipAttachment{
		FileID:         fileID,
		SlipID:         structs.PointerToNullableString(&admissionSlipID),
		AttachmentType: "OTHER",
	}

	err = slipsRepo.SaveSlipAttachment(ctx, tx, attachment)
	if err != nil {
		log.Printf("[Seeder] {Insert SlipAttachment}: %v", err)
		return ""
	}

	return admissionSlipID
}

func chooseAdmissionSlipStatus(absDate, now time.Time) string {
	// Future absences can only be pending
	if absDate.After(now) {
		return "pending"
	}

	// Past absences are processed
	statuses := []struct {
		name   string
		weight int
	}{
		{"approved", 70},
		{"rejected", 20},
		{"for revision", 10},
	}

	totalWeight := 0
	for _, s := range statuses {
		totalWeight += s.weight
	}

	roll := rand.Intn(totalWeight)
	for _, s := range statuses {
		roll -= s.weight
		if roll < 0 {
			return s.name
		}
	}
	return "approved"
}

func generateRealisticAdmissionReason() string {
	reasons := []string{
		"Medical consultation - Doctor's appointment",
		"Family emergency requiring immediate attention",
		"Legal matter - Court appearance required",
		"University business - Office of the Registrar",
		"Health center visit - Medical examination",
		"Dental appointment - Preventive care",
		"Job interview - Employment opportunity",
		"Scholarship office - Financial aid matters",
		"Hospital admission - Medical treatment",
		"Burial ceremony - Family death",
		"Important family matter - Guardian meeting",
		"University activity - Student organization duty",
		"Mental health counseling session",
		"Physical therapy appointment",
		"Government office work - Official business",
	}
	return reasons[rand.Intn(len(reasons))]
}

func generateAdmissionFileName() string {
	fileTypes := []string{
		"medical_certificate",
		"excuse_letter",
		"hospital_clearance",
		"doctor_note",
		"appointment_slip",
		"official_letter",
		"emergency_report",
		"proof_of_attendance",
	}
	return fileTypes[rand.Intn(len(fileTypes))]
}

func generateAdmissionNotes(status string) structs.NullableString {
	switch status {
	case "approved":
		notes := []string{
			"Verified with hospital records",
			"Legitimate documentation provided",
			"Approved - excused absence",
			"Admission approved as submitted",
			"",
		}
		selected := notes[rand.Intn(len(notes))]
		if selected == "" {
			return structs.NullableString{Valid: false}
		}
		return structs.NullableString{String: selected, Valid: true}
	case "rejected":
		notes := []string{
			"Insufficient documentation",
			"Outside acceptable timeframe",
			"Documentation not properly verified",
			"Required supporting documents missing",
		}
		return structs.NullableString{
			String: notes[rand.Intn(len(notes))],
			Valid:  true,
		}
	case "for revision":
		notes := []string{
			"Requires additional supporting documents",
			"Need clarification on dates",
			"Missing required signatures",
			"Please provide recent medical certificate",
		}
		return structs.NullableString{
			String: notes[rand.Intn(len(notes))],
			Valid:  true,
		}
	default:
		return structs.NullableString{Valid: false}
	}
}

func insertAppointment(
	ctx context.Context,
	tx *sqlx.Tx,
	iirID string,
	appointmentsDataset []map[string]string,
) string {
	if len(timeSlotIDs) == 0 || len(appointmentCategories) == 0 ||
		len(appointmentStatusIDs) == 0 {
		return ""
	}

	hasDataset := false
	if len(appointmentsDataset) > 0 {
		hasDataset = true
	}

	selectedIdx := 0
	var selectedAppt map[string]string
	if hasDataset {
		selectedIdx = rand.Intn(len(appointmentsDataset))
		selectedAppt = (appointmentsDataset)[selectedIdx]
	} else {
		selectedAppt = map[string]string{
			"text":          gofakeit.Sentence(rand.Intn(5) + 5),
			"urgency_level": []string{"LOW", "MEDIUM", "HIGH", "CRITICAL"}[rand.Intn(4)],
			"category":      []string{"ACADEMIC", "CAREER GUIDANCE", "PERSONAL", "MENTAL HEALTH", "FINANCIAL", "OTHER"}[rand.Intn(6)],
		}
	}

	now := time.Now().Truncate(24 * time.Hour)
	whenDateStr, timeSlotID, whenDate := reserveAppointmentSlot(now)
	statusID := chooseAppointmentStatusID(whenDate, now)
	concernCategoryID := func() string {
		for _, category := range appointmentCategories {
			if strings.EqualFold(category["name"], selectedAppt["category"]) {
				return category["id"]
			}
		}

		return "0"
	}()

	concernCategoryIDInt, _ := strconv.Atoi(concernCategoryID)

	// Determine admin_notes value
	var adminNotes structs.NullableString
	statusName := ""
	for name, id := range appointmentStatusByName {
		if id == statusID {
			statusName = strings.ToLower(name)
			break
		}
	}
	if statusName == "cancelled" || statusName == "rejected" ||
		strings.Contains(statusName, "show") || statusName == "completed" {
		adminNotes = structs.NullableString{
			String: gofakeit.Sentence(rand.Intn(5) + 5),
			Valid:  true,
		}
	} else {
		adminNotes = structs.NullableString{Valid: false}
	}

	appointmentID := uuid.New().String()

	appt := &appointments.Appointment{
		ID:                    appointmentID,
		IIRID:                 iirID,
		Reason:                structs.StringToNullableString(selectedAppt["text"]),
		AdminNotes:            adminNotes,
		WhenDate:              whenDateStr,
		TimeSlotID:            timeSlotID,
		UrgencyLevel:          selectedAppt["urgency_level"],
		UrgencyScore:          0.95,
		AppointmentCategoryID: concernCategoryIDInt,
		StatusID:              statusID,
	}

	err := appointmentsRepo.CreateAppointment(ctx, tx, appt)
	if err != nil {
		log.Printf("[Seeder] {Insert Appointment}: %v", err)
		return ""
	}

	return appointmentID
}

func reserveAppointmentSlot(now time.Time) (string, int, time.Time) {
	if len(timeSlotIDs) == 0 {
		log.Fatal("no time slots found in time_slots")
	}

	for attempts := 0; attempts < 1000; attempts++ {
		// Randomly decide if past or future
		var when time.Time
		if rand.Float32() < 0.3 {
			// Future (up to 2 weeks)
			when = now.AddDate(0, 0, rand.Intn(14))
		} else {
			// Past (up to 90 days)
			when = now.AddDate(0, 0, -(rand.Intn(90) + 1))
		}

		weekday := when.Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			continue
		}
		whenDate := when.Format("2006-01-02")
		timeSlotID := randomChoice(timeSlotIDs).(int)
		key := fmt.Sprintf("%s|%d", whenDate, timeSlotID)

		appointmentSlotMu.Lock()
		_, exists := reservedAppointmentSlots[key]
		if !exists {
			reservedAppointmentSlots[key] = struct{}{}
			appointmentSlotMu.Unlock()
			return whenDate, timeSlotID, when
		}
		appointmentSlotMu.Unlock()
	}

	log.Fatal("unable to reserve unique appointment slot")
	return "", 0, time.Time{}
}

func chooseAppointmentStatusID(apptDate, now time.Time) int {
	if len(appointmentStatusByName) == 0 {
		return randomChoice(appointmentStatusIDs).(int)
	}

	type weightedStatus struct {
		id     int
		weight int
	}

	weighted := make([]weightedStatus, 0)
	used := make(map[int]bool)

	add := func(name string, weight int) {
		if id, ok := appointmentStatusByName[name]; ok {
			weighted = append(weighted, weightedStatus{id: id, weight: weight})
			used[id] = true
		}
	}

	isFuture := apptDate.After(now) || apptDate.Equal(now)

	if isFuture {
		// Future: pending, approved (scheduled), rescheduled
		add("pending", 60)
		add("approved", 30) // means Scheduled
		add("rescheduled", 10)
	} else {
		// Past: completed, no-show, rejected, cancelled
		add("completed", 60)
		add("no-show", 20)
		add("rejected", 10)
		add("cancelled", 10)
	}

	if len(weighted) == 0 {
		return randomChoice(appointmentStatusIDs).(int)
	}

	totalWeight := 0
	for _, item := range weighted {
		totalWeight += item.weight
	}

	roll := rand.Intn(totalWeight)
	for _, item := range weighted {
		roll -= item.weight
		if roll < 0 {
			return item.id
		}
	}

	return weighted[len(weighted)-1].id
}

func chooseCivilStatusID() int {
	if len(civilStatusByName) == 0 {
		return randomChoice(civilStatusIDs).(int)
	}

	type weightedStatus struct {
		id     int
		weight int
	}

	weighted := make([]weightedStatus, 0)
	used := make(map[int]bool)

	add := func(name string, weight int) {
		if id, ok := civilStatusByName[name]; ok {
			weighted = append(weighted, weightedStatus{id: id, weight: weight})
			used[id] = true
		}
	}

	add("single", 75)
	add("married", 8)
	add("widowed", 2)
	add("divorced", 2)
	add("separated", 6)

	for _, id := range civilStatusIDs {
		if !used[id] {
			weighted = append(weighted, weightedStatus{id: id, weight: 4})
		}
	}

	if len(weighted) == 0 {
		return randomChoice(civilStatusIDs).(int)
	}

	totalWeight := 0
	for _, item := range weighted {
		totalWeight += item.weight
	}

	roll := rand.Intn(totalWeight)
	for _, item := range weighted {
		roll -= item.weight
		if roll < 0 {
			return item.id
		}
	}

	return weighted[len(weighted)-1].id
}
