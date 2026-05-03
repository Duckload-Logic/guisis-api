package main

import (
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/slips"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

var db *sqlx.DB

// global repositories
var (
	usersRepo        *users.Repository
	studentsRepo     *students.Repository
	appointmentsRepo *appointments.Repository
	slipsRepo        *slips.Repository
	locationsRepo    *locations.Repository
)

// notificationTypes is now defined in student.go or can be moved to a shared
// location if needed.
// Actually, I'll keep it here and make sure student.go doesn't redefine it.
var notificationTypes = []string{
	"System",
	"Appointment",
	"Slip",
	"General",
	"Guidance",
}

// lookup slices (IDs)
var (
	genderIDs                   []int
	civilStatusIDs              []int
	civilStatusByName           map[string]int
	religionIDs                 []int
	courseIDs                   []int
	enrollmentReasonIDs         []int
	supportTypeIDs              []int
	incomeRangeIDs              []int
	parentalStatusIDs           []int
	educationalLevelIDs         []int
	educationalLevelByName      map[string]int
	relationshipTypeIDs         map[string]int
	natureOfResidenceIDs        []int
	siblingSupportTypeIDs       []int
	activityOptionIDs           []int
	timeSlotIDs                 []int
	appointmentStatusIDs        []int
	appointmentStatusByName     map[string]int
	admissionSlipStatusIDs      []int
	admissionSlipStatusesByName map[string]int
	appointmentCategories       []map[string]string
	admissionSlipCategoryIDs    []int
	studentStatusIDs            []int
	studentStatusByName         map[string]int
	appointmentSlotMu           sync.Mutex
	reservedAppointmentSlots    = make(map[string]struct{})
)

type relatedPersonSeed struct {
	ID            int
	FirstName     string
	MiddleName    structs.NullableString
	LastName      string
	ContactNumber structs.NullableString
	AddressID     *int
}

type emergencyContactSeed struct {
	ID             int
	FirstName      string
	MiddleName     structs.NullableString
	LastName       string
	Number         string
	RelationshipID int
	AddressID      int
}
