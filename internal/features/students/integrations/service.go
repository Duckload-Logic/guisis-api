package integrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
)

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListStudents(
	ctx context.Context,
	req OGOSListStudentsRequest,
) (OGOSListStudentsResponse, error) {
	req.SetDefaults("created_at")

	studentList, total, err := s.repo.ListStudents(ctx, req)
	if err != nil {
		return OGOSListStudentsResponse{}, err
	}

	var studentsDTO []OGOSStudentDTO
	for _, student := range studentList {
		studentsDTO = append(studentsDTO, OGOSStudentDTO{
			StudentNumber: student.StudentNumber,
			FirstName:     student.FirstName,
			MiddleName:    structs.FromSqlNull(student.MiddleName),
			LastName:      student.LastName,
			Email:         student.Email,
			MobileNumber:  student.MobileNumber,
			Course: students.Course{
				ID:   student.CourseID,
				Code: student.CourseCode,
				Name: student.CourseName,
			},
			YearLevel: student.YearLevel,
			Section:   student.Section,
		})
	}

	listResponse := OGOSListStudentsResponse{
		Students: studentsDTO,
		Meta:     structs.CalculateMetadata(total, req.Page, req.PageSize),
	}

	return listResponse, nil
}

func (s *Service) GetStudentByStudentNumber(
	ctx context.Context,
	studentNumber string,
) (*OGOSStudentDTO, error) {
	student, err := s.repo.GetStudentByStudentNumber(ctx, studentNumber)
	if err != nil {
		return nil, err
	}

	return &OGOSStudentDTO{
		StudentNumber: student.StudentNumber,
		FirstName:     student.FirstName,
		MiddleName:    structs.FromSqlNull(student.MiddleName),
		LastName:      student.LastName,
		Email:         student.Email,
		MobileNumber:  student.MobileNumber,
		Course: students.Course{
			ID:   student.CourseID,
			Code: student.CourseCode,
			Name: student.CourseName,
		},
		YearLevel: student.YearLevel,
		Section:   student.Section,
	}, nil
}

func (s *Service) GetStudentByEmail(
	ctx context.Context,
	email string,
) (*OGOSStudentDTO, error) {
	student, err := s.repo.GetStudentByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("student not found")
		}

		return nil, err
	}

	return &OGOSStudentDTO{
		StudentNumber: student.StudentNumber,
		FirstName:     student.FirstName,
		MiddleName:    structs.FromSqlNull(student.MiddleName),
		LastName:      student.LastName,
		Email:         student.Email,
		MobileNumber:  student.MobileNumber,
		Course: students.Course{
			ID:   student.CourseID,
			Code: student.CourseCode,
			Name: student.CourseName,
		},
		YearLevel: student.YearLevel,
		Section:   student.Section,
	}, nil
}

func (s *Service) GetPersonalInfoByStudentNumber(
	ctx context.Context,
	studentNumber string,
) (*OGOSStudentPersonalInfoDTO, error) {
	student, err := s.repo.GetPersonalInfoByStudentNumber(ctx, studentNumber)
	if err != nil {
		return nil, err
	}

	return &OGOSStudentPersonalInfoDTO{
		StudentNumber: student.StudentNumber,
		Gender: students.Gender{
			ID:   student.GenderID,
			Name: student.GenderName,
		},
		DateOfBirth:  student.DateOfBirth,
		PlaceOfBirth: student.PlaceOfBirth,
		HeightFt:     student.HeightFt,
		WeightKg:     student.WeightKg,
	}, nil
}

func (s *Service) GetAddressByStudentNumber(
	ctx context.Context,
	studentNumber string,
) ([]OGOSStudentAddressDTO, error) {
	studentAddresses, err := s.repo.GetAddressByStudentNumber(
		ctx,
		studentNumber,
	)
	if err != nil {
		return nil, err
	}

	if len(studentAddresses) == 0 {
		return nil, fmt.Errorf(
			"no addresses found for student number: %s",
			studentNumber,
		)
	}

	var addresesDTO []OGOSStudentAddressDTO
	for _, address := range studentAddresses {
		addresesDTO = append(addresesDTO, OGOSStudentAddressDTO{
			StudentNumber: address.StudentNumber,
			AddressType:   address.AddressType,
			StreetDetail:  address.StreetDetail,
			Barangay: locations.Barangay{
				Code: address.BarangayCode,
				Name: address.BarangayName,
			},
			City: locations.City{
				Code: address.CityCode,
				Name: address.CityName,
			},
			Province: &locations.ProvinceDTO{
				Code: structs.FromSqlNull(address.ProvinceCode),
				Name: structs.FromSqlNull(address.ProvinceName),
			},
			Region: locations.Region{
				Code: address.RegionCode,
				Name: address.RegionName,
			},
		})
	}

	return addresesDTO, nil
}
