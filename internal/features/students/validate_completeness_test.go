package students

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateIIRCompleteness(t *testing.T) {
	svc := &Service{}

	tests := []struct {
		name        string
		profile     *ComprehensiveProfileDTO
		wantErr     bool
		expectedErr error
		errContains string
	}{
		{
			name: "Both Elementary, Father, and Mother present - Valid",
			profile: &ComprehensiveProfileDTO{
				Education: EducationalBackgroundDTO{
					School: []SchoolDetailsDTO{
						{
							EducationalLevel: EducationalLevel{
								Name: "Elementary",
							},
							SchoolName: "PUP Lab High/Elem",
						},
					},
				},
				Family: struct {
					FamilyBackgroundDTO `json:"background"`
					RelatedPersons      []RelatedPersonDTO `json:"relatedPersons"`
					Finance             StudentFinanceDTO  `json:"finance"`
				}{
					RelatedPersons: []RelatedPersonDTO{
						{
							FirstName: "Juan",
							LastName:  "Santos",
							Relationship: StudentRelationshipType{
								ID:   1,
								Name: "Father",
							},
						},
						{
							FirstName: "Maria",
							LastName:  "Santos",
							Relationship: StudentRelationshipType{
								ID:   2,
								Name: "Mother",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Missing Elementary - Invalid",
			profile: &ComprehensiveProfileDTO{
				Education: EducationalBackgroundDTO{
					School: []SchoolDetailsDTO{
						{
							EducationalLevel: EducationalLevel{
								Name: "High School",
							},
							SchoolName: "Test HS",
						},
					},
				},
				Family: struct {
					FamilyBackgroundDTO `json:"background"`
					RelatedPersons      []RelatedPersonDTO `json:"relatedPersons"`
					Finance             StudentFinanceDTO  `json:"finance"`
				}{
					RelatedPersons: []RelatedPersonDTO{
						{
							FirstName: "Juan",
							LastName:  "Santos",
							Relationship: StudentRelationshipType{
								ID:   1,
								Name: "Father",
							},
						},
						{
							FirstName: "Maria",
							LastName:  "Santos",
							Relationship: StudentRelationshipType{
								ID:   2,
								Name: "Mother",
							},
						},
					},
				},
			},
			wantErr:     true,
			expectedErr: ErrIncompleteIIR,
			errContains: "Elementary School",
		},
		{
			name: "Missing Father only - Invalid",
			profile: &ComprehensiveProfileDTO{
				Education: EducationalBackgroundDTO{
					School: []SchoolDetailsDTO{
						{
							EducationalLevel: EducationalLevel{
								Name: "Elementary",
							},
							SchoolName: "Test Elem",
						},
					},
				},
				Family: struct {
					FamilyBackgroundDTO `json:"background"`
					RelatedPersons      []RelatedPersonDTO `json:"relatedPersons"`
					Finance             StudentFinanceDTO  `json:"finance"`
				}{
					RelatedPersons: []RelatedPersonDTO{
						{
							FirstName: "Maria",
							LastName:  "Santos",
							Relationship: StudentRelationshipType{
								ID:   2,
								Name: "Mother",
							},
						},
					},
				},
			},
			wantErr:     true,
			expectedErr: ErrIncompleteIIR,
			errContains: "Father's information",
		},
		{
			name: "Missing Mother only - Invalid",
			profile: &ComprehensiveProfileDTO{
				Education: EducationalBackgroundDTO{
					School: []SchoolDetailsDTO{
						{
							EducationalLevel: EducationalLevel{
								Name: "Elementary",
							},
							SchoolName: "Test Elem",
						},
					},
				},
				Family: struct {
					FamilyBackgroundDTO `json:"background"`
					RelatedPersons      []RelatedPersonDTO `json:"relatedPersons"`
					Finance             StudentFinanceDTO  `json:"finance"`
				}{
					RelatedPersons: []RelatedPersonDTO{
						{
							FirstName: "Juan",
							LastName:  "Santos",
							Relationship: StudentRelationshipType{
								ID:   1,
								Name: "Father",
							},
						},
					},
				},
			},
			wantErr:     true,
			expectedErr: ErrIncompleteIIR,
			errContains: "Mother's information",
		},
		{
			name: "Missing Both Father and Mother - Invalid",
			profile: &ComprehensiveProfileDTO{
				Education: EducationalBackgroundDTO{
					School: []SchoolDetailsDTO{
						{
							EducationalLevel: EducationalLevel{
								Name: "Elementary",
							},
							SchoolName: "Test Elem",
						},
					},
				},
				Family: struct {
					FamilyBackgroundDTO `json:"background"`
					RelatedPersons      []RelatedPersonDTO `json:"relatedPersons"`
					Finance             StudentFinanceDTO  `json:"finance"`
				}{
					RelatedPersons: []RelatedPersonDTO{},
				},
			},
			wantErr:     true,
			expectedErr: ErrIncompleteIIR,
			errContains: "Father and Mother information",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := svc.validateIIRCompleteness(tc.profile)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !errors.Is(err, tc.expectedErr) {
					t.Errorf(
						"expected error to wrap %v, got %v",
						tc.expectedErr,
						err,
					)
				}
				if tc.errContains != "" &&
					!strings.Contains(err.Error(), tc.errContains) {
					t.Errorf(
						"expected error to contain %q, got %q",
						tc.errContains,
						err.Error(),
					)
				}
			} else if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
